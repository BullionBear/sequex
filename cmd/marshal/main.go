package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Define flags
	debugFlag := flag.Bool("d", false, "debug mode - add line numbers to output [line]")
	negativeFlag := flag.Bool("s", false, "negative mode - add negative line numbers to output [-line]")
	flag.Parse()

	// Validate flags - exactly one of -d or -s must be specified
	if *debugFlag && *negativeFlag {
		fmt.Fprintf(os.Stderr, "Error: cannot use both -d and -s flags together\n")
		flag.Usage()
		os.Exit(1)
	}
	if !*debugFlag && !*negativeFlag {
		fmt.Fprintf(os.Stderr, "Error: must specify either -d or -s flag\n")
		flag.Usage()
		os.Exit(1)
	}

	var reader io.Reader
	var filename string

	// Check if we have a file argument
	args := flag.Args()
	if len(args) > 0 {
		// Read from file
		filename = args[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		// Check if stdin has data available (not a terminal)
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking stdin: %v\n", err)
			flag.Usage()
			os.Exit(1)
		}

		// If stdin is a terminal (no piped data), show error
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintf(os.Stderr, "Error: no input provided. Please specify a file or pipe data to stdin\n")
			flag.Usage()
			os.Exit(1)
		}

		// Read from stdin (piped input)
		reader = os.Stdin
		filename = "stdin"
	}

	// Process the input
	if err := processInput(reader, *debugFlag, *negativeFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing input: %v\n", err)
		os.Exit(1)
	}
}

func processInput(reader io.Reader, debug bool, negative bool) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if debug {
			fmt.Printf("[%d] %s\n", lineNumber, line)
		} else if negative {
			fmt.Printf("[-%d] %s\n", lineNumber, line)
		} else {
			fmt.Println(line)
		}
		lineNumber++
	}

	return scanner.Err()
}
