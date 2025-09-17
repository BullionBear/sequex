package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BullionBear/sequex/internal/model/protobuf"
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
		for len(accumulated) >= 80 { // Minimum message size
			found := false

			// Try common message sizes (based on protobuf analysis)
			for _, size := range getMessageSizes() {
				if len(accumulated) < size {
					continue
				}

				candidate := accumulated[:size]
				trade := &protobuf.Trade{}

				if err := proto.Unmarshal(candidate, trade); err == nil {
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

						accumulated = accumulated[size:]
						found = true
						break
					}
				}
			}

			if !found {
				// Skip one byte and try again
				accumulated = accumulated[1:]
			}
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

// getMessageSizes returns the common message sizes observed in the protobuf data
func getMessageSizes() []int {
	return []int{45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60}
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
	fmt.Printf("Trade %d:\n", messageNum)
	fmt.Printf("  ID: %d\n", trade.Id)
	fmt.Printf("  Exchange: %s (%d)\n", trade.Exchange.String(), int(trade.Exchange))
	fmt.Printf("  Instrument: %s (%d)\n", trade.Instrument.String(), int(trade.Instrument))

	if trade.Symbol != nil {
		fmt.Printf("  Symbol: %s/%s\n", trade.Symbol.Base, trade.Symbol.Quote)
	} else {
		fmt.Printf("  Symbol: <nil>\n")
	}

	fmt.Printf("  Side: %s (%d)\n", trade.Side.String(), int(trade.Side))
	fmt.Printf("  Price: %.8f\n", trade.Price)
	fmt.Printf("  Quantity: %.8f\n", trade.Quantity)

	if trade.Timestamp > 0 {
		t := time.Unix(trade.Timestamp/1000, (trade.Timestamp%1000)*1000000)
		fmt.Printf("  Timestamp: %d (%s)\n", trade.Timestamp, t.Format("2006-01-02 15:04:05.000"))
	} else {
		fmt.Printf("  Timestamp: %d\n", trade.Timestamp)
	}

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
