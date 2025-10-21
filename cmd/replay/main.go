package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/BullionBear/sequex/internal/model/protobuf"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"google.golang.org/protobuf/proto"
)

var (
	inputFile   = flag.String("input", "messages-20250915.raw", "Input file containing serialized protobuf messages")
	showLimit   = flag.Int("limit", 100, "Number of messages to display (0 for all)")
	showSummary = flag.Bool("summary", true, "Show summary statistics")
	verbose     = flag.Bool("verbose", false, "Show verbose output")
)

func main() {
	flag.Parse()

	fmt.Println("Sequex Trade Message Replay Tool")
	fmt.Println(strings.Repeat("=", 40))

	if *verbose {
		fmt.Printf("Input file: %s\n", *inputFile)
		fmt.Printf("Display limit: %d\n", *showLimit)
		fmt.Printf("Show summary: %v\n", *showSummary)
		fmt.Println()
	}

	successCount, totalProcessed, err := replayTradeMessages(*inputFile)
	if err != nil {
		log.Fatalf("Failed to replay messages: %v", err)
	}

	if *showSummary {
		printSummary(successCount, totalProcessed)
	}
}

func replayTradeMessages(filename string) (successCount, totalProcessed int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024) // 1MB buffer
	var accumulated []byte

	for {
		n, readErr := file.Read(buffer)

		if n > 0 {
			accumulated = append(accumulated, buffer[:n]...)
		}

		if readErr == io.EOF && len(accumulated) == 0 {
			break
		}

		// Process accumulated data
		for len(accumulated) >= 10 { // Minimum viable message size
			messageData, consumed, found := parseNextMessage(accumulated)
			if !found {
				// Skip one byte and try again
				accumulated = accumulated[1:]
				continue
			}

			trade := &protobuf.Trade{}
			if err := proto.Unmarshal(messageData, trade); err == nil {
				totalProcessed++

				// Validate trade message
				if isValidTradeMessage(trade) {
					successCount++

					// Display message if within limit
					if *showLimit == 0 || successCount <= *showLimit {
						displayTradeMessage(successCount, trade)
					} else if successCount == *showLimit+1 {
						fmt.Printf("... (limiting output to first %d messages)\n\n", *showLimit)
					}
				}
			}

			accumulated = accumulated[consumed:]
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return successCount, totalProcessed, fmt.Errorf("error reading file: %w", readErr)
		}
	}

	return successCount, totalProcessed, nil
}

// parseNextMessage parses the next complete protobuf message from the data
func parseNextMessage(data []byte) (messageData []byte, consumed int, found bool) {
	if len(data) < 10 {
		return nil, 0, false
	}

	// Parse protobuf wire format to find message boundaries
	offset := 0
	fieldsSeen := make(map[int]bool)

	for offset < len(data) && offset < 200 { // Reasonable upper bound
		if offset+1 >= len(data) {
			break
		}

		// Read field header (field number + wire type)
		fieldHeader := data[offset]
		fieldNum := int(fieldHeader >> 3)
		wireType := int(fieldHeader & 0x7)
		offset++

		// Skip invalid field numbers (protobuf fields start at 1)
		if fieldNum == 0 || fieldNum > 20 {
			break
		}

		// Skip the field data based on wire type
		fieldLength, ok := skipFieldData(data[offset:], wireType)
		if !ok {
			break
		}
		offset += fieldLength
		fieldsSeen[fieldNum] = true

		// Check if we have a complete Trade message
		// Trade has fields: 1=id, 2=exchange, 3=instrument, 4=symbol, 5=side, 7=price, 8=quantity, 9=timestamp
		if hasAllExpectedFields(fieldsSeen) {
			// We've seen all expected fields, try to parse
			candidate := data[:offset]
			trade := &protobuf.Trade{}
			if err := proto.Unmarshal(candidate, trade); err == nil && isValidTradeMessage(trade) {
				return candidate, offset, true
			}
		}
	}

	return nil, 0, false
}

// skipFieldData skips over field data based on wire type
func skipFieldData(data []byte, wireType int) (int, bool) {
	switch wireType {
	case 0: // Varint
		return skipVarint(data)
	case 1: // 64-bit fixed
		if len(data) < 8 {
			return 0, false
		}
		return 8, true
	case 2: // Length-delimited (strings, bytes, embedded messages)
		return skipLengthDelimited(data)
	case 5: // 32-bit fixed
		if len(data) < 4 {
			return 0, false
		}
		return 4, true
	default:
		return 0, false // Unknown wire type
	}
}

// skipVarint skips over a varint-encoded value
func skipVarint(data []byte) (int, bool) {
	for i := 0; i < len(data) && i < 10; i++ { // Max 10 bytes for varint
		if data[i]&0x80 == 0 {
			return i + 1, true
		}
	}
	return 0, false
}

// skipLengthDelimited skips over a length-delimited field
func skipLengthDelimited(data []byte) (int, bool) {
	// Decode the length varint
	length := uint64(0)
	lengthBytes := 0

	for i := 0; i < len(data) && i < 10; i++ {
		length |= uint64(data[i]&0x7F) << (7 * i)
		lengthBytes++
		if data[i]&0x80 == 0 {
			break
		}
	}

	if lengthBytes == 0 {
		return 0, false
	}

	// Skip the length prefix + the data
	totalLength := lengthBytes + int(length)
	if len(data) < totalLength {
		return 0, false
	}
	return totalLength, true
}

// hasAllExpectedFields checks if we've seen all the expected fields for a complete Trade message
func hasAllExpectedFields(fieldsSeen map[int]bool) bool {
	// All expected Trade fields: 1=id, 2=exchange, 3=instrument, 4=symbol, 5=side, 7=price, 8=quantity, 9=timestamp
	expectedFields := []int{1, 2, 3, 4, 5, 7, 8, 9}

	for _, field := range expectedFields {
		if !fieldsSeen[field] {
			return false
		}
	}

	return true
}

// isValidTradeMessage validates that a Trade message contains reasonable data
func isValidTradeMessage(trade *protobuf.Trade) bool {
	validFields := 0

	// ID should be positive
	if trade.Id > 0 {
		validFields++
	}

	// Exchange should be valid (1-3 for known exchanges)
	if trade.Exchange >= 1 && trade.Exchange <= 3 {
		validFields++
	}

	// Instrument should be valid
	if trade.Instrument >= 1 && trade.Instrument <= 6 {
		validFields++
	}

	// Symbol should exist and have reasonable values
	if trade.Symbol != nil && len(trade.Symbol.Base) >= 2 && len(trade.Symbol.Quote) >= 3 {
		validFields++
	}

	// Side should be buy or sell
	if trade.Side >= 1 && trade.Side <= 2 {
		validFields++
	}

	// Price should be reasonable (between $0.01 and $1M)
	if trade.Price >= 0.01 && trade.Price <= 1000000 {
		validFields++
	}

	// Quantity should be positive
	if trade.Quantity > 0 {
		validFields++
	}

	// Timestamp should be reasonable (2020-2030)
	if trade.Timestamp >= 1577836800000 && trade.Timestamp <= 1893456000000 {
		validFields++
	}

	// Require at least 6 out of 8 fields to be valid
	return validFields >= 6
}

// displayTradeMessage prints a formatted trade message
func displayTradeMessage(messageNum int, trade *protobuf.Trade) {
	sqxTrade := &sqx.Trade{}
	if err := sqxTrade.FromProtobuf(trade); err != nil {
		fmt.Printf("Failed to deserialize trade message: %v\n", err)
		return
	}
	fmt.Printf("Trade %d:\n", messageNum)
	data, err := json.MarshalIndent(sqxTrade, "", "  ")
	if err != nil {
		fmt.Printf("Failed to serialize trade message to JSON: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(data))

	fmt.Printf("\n")
}

// printSummary displays summary statistics
func printSummary(successCount, totalProcessed int) {
	fmt.Printf(strings.Repeat("=", 50) + "\n")
	fmt.Printf("Summary:\n")
	fmt.Printf("Successfully deserialized: %d complete messages\n", successCount)
	fmt.Printf("Total messages processed: %d\n", totalProcessed)
	if totalProcessed > 0 {
		fmt.Printf("Success rate: %.2f%%\n", float64(successCount)/float64(totalProcessed)*100)
	}
	fmt.Printf("Input file: %s\n", *inputFile)

	// Additional statistics
	if successCount > 0 {
		fmt.Printf("\nReplay completed successfully!\n")
	} else {
		fmt.Printf("\nNo valid trade messages found. Check input file format.\n")
	}
}
