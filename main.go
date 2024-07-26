package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

func main() {
	resp := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	reader := bufio.NewReader(bytes.NewBufferString(resp))

	commands, err := parseRESPCommand(reader)
	if err != nil {
		log.Fatalf("Error parsing RESP: %v", err)
	}

	fmt.Println("Parsed commands:", commands)
}
