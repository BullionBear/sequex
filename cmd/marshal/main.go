package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BullionBear/sequex/internal/model/protobuf"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Define flags
	deserializeFlag := flag.Bool("d", false, "deserialize mode - convert .raw protobuf file to JSON format")
	serializeFlag := flag.Bool("s", false, "serialize mode - convert JSON to protobuf .raw format")
	outputFile := flag.String("o", "", "output file (default: stdout for -d, required for -s)")
	flag.Parse()

	// Validate flags - exactly one of -d or -s must be specified
	if *deserializeFlag && *serializeFlag {
		fmt.Fprintf(os.Stderr, "Error: cannot use both -d and -s flags together\n")
		flag.Usage()
		os.Exit(1)
	}
	if !*deserializeFlag && !*serializeFlag {
		fmt.Fprintf(os.Stderr, "Error: must specify either -d (deserialize) or -s (serialize) flag\n")
		flag.Usage()
		os.Exit(1)
	}

	// Get input file (optional - if not provided, read from stdin)
	args := flag.Args()
	var inputFile string
	if len(args) > 0 {
		inputFile = args[0]
	}

	// Validate output file for serialize mode
	if *serializeFlag && *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: output file (-o) required for serialize mode (-s)\n")
		flag.Usage()
		os.Exit(1)
	}

	// Process based on mode
	if *deserializeFlag {
		if err := deserializeMode(inputFile, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error in deserialize mode: %v\n", err)
			os.Exit(1)
		}
	} else if *serializeFlag {
		if err := serializeMode(inputFile, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error in serialize mode: %v\n", err)
			os.Exit(1)
		}
	}
}

// deserializeMode reads a .raw protobuf file and outputs JSON
func deserializeMode(inputFile, outputFile string) error {
	var file *os.File
	var err error

	if inputFile == "" {
		// Read from stdin
		file = os.Stdin
	} else {
		file, err = os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("failed to open input file %s: %w", inputFile, err)
		}
		defer file.Close()
	}

	// Setup output writer
	var writer io.Writer = os.Stdout
	if outputFile != "" {
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
		}
		defer outFile.Close()
		writer = outFile
	}

	buffer := make([]byte, 1024*1024) // 1MB buffer
	var accumulated []byte
	messageCount := 0

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
				// Convert to SQX format and output as JSON
				sqxTrade := &sqx.Trade{}
				if err := sqxTrade.FromProtobuf(trade); err == nil {
					jsonData, err := json.Marshal(sqxTrade)
					if err == nil {
						fmt.Fprintf(writer, "%s\n", string(jsonData))
						messageCount++
					}
				}
			}

			accumulated = accumulated[consumed:]
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return fmt.Errorf("error reading file: %w", readErr)
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully deserialized %d messages\n", messageCount)
	return nil
}

// serializeMode reads JSON input and writes protobuf .raw file
func serializeMode(inputFile, outputFile string) error {
	var inputReader *os.File
	var err error

	if inputFile == "" {
		// Read from stdin
		inputReader = os.Stdin
	} else {
		inputReader, err = os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("failed to open input file %s: %w", inputFile, err)
		}
		defer inputReader.Close()
	}

	outputWriter, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer outputWriter.Close()

	scanner := bufio.NewScanner(inputReader)
	messageCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		// Parse JSON to SQX Trade
		var sqxTrade sqx.Trade
		if err := json.Unmarshal([]byte(line), &sqxTrade); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse JSON line %d: %v\n", messageCount+1, err)
			continue
		}

		// Convert to protobuf and marshal
		pbTrade := sqxTrade.ToProtobuf()
		data, err := proto.Marshal(pbTrade)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to marshal protobuf for line %d: %v\n", messageCount+1, err)
			continue
		}

		// Write raw protobuf data
		if _, err := outputWriter.Write(data); err != nil {
			return fmt.Errorf("failed to write protobuf data: %w", err)
		}

		messageCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Successfully serialized %d messages\n", messageCount)
	return nil
}

// parseNextMessage parses the next complete protobuf message from the data
// This is adapted from the replay tool
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
