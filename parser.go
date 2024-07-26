package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n -> ["ECHO", "hey"]
// '+', '-', ':', '$', '~', '%' // these characters (and some others) are to be ignored since they don't represent the start of an RESP Array.
// A '*' does
func parseRESPCommand(reader *bufio.Reader) ([]string, error) {
	// Read string - after reading the string up to the delimiter, the reader's internal buffer is advanced.
	// This means the next read operation will start from where the previous read ended,
	// skipping the data that has already been consumed.
	start, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	start = strings.TrimSpace(start)
	if len(start) == 0 {
		return nil, fmt.Errorf("error: empty payload")
	}

	if start[0] != '*' {
		return nil, fmt.Errorf("error: invalid command")
	}

	numOfElements, _ := strconv.Atoi(start[1:]) // the count of total commands and their arguments

	return parseRESPArrayFromReader(reader, numOfElements)
}

// $4\r\nECHO\r\n$3\r\nhey\r\n -> ["ECHO", "hey"]
func parseRESPArrayFromReader(reader *bufio.Reader, numOfElements int) ([]string, error) {
	respArgs := make([]string, numOfElements)

	for i := 0; i < numOfElements; i++ {
		start, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		start = strings.TrimSpace(start)
		if len(start) == 0 || start[0] != '$' {
			return nil, fmt.Errorf("error: expected a bulk string")
		}
		argumentLength, _ := strconv.Atoi(start[1:]) // get the length of the current argument

		// Allocate space for the argument and an extra two bytes to account for the newline characters (`\r\n`).
		// We're reading from the reader.
		argument := make([]byte, argumentLength+2)
		reader.Read(argument)
		respArgs[i] = strings.TrimSpace(string(argument))
	}
	return respArgs, nil
}
