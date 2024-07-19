package io

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"moba-converter-go/internal/config"
	"os"
)

// Deal with loading and saving mxtsession files and json files.

func loadFile(path *string) []byte {
	// Load a file from a path or from stdin if path is "".
	var data []byte
	var err error
	if *path == "" {
		// Read from file
		data, err = os.ReadFile(*path)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
	} else {
		// Read from stdin
		fmt.Fprintln(os.Stderr, "<<Reading from stdin.>>")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data = append(data, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}

	}
	return data
}

func LoadJsonInput(path *string) config.JSONInput {
	//Load a Json Input file from a path or from stdin if path is "".
	data := loadFile(path)

	var input config.JSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		log.Fatalf("Error parsing input file: %v", err)
	}

	return input
}

func LoadMxtsessionInput(path *string) []byte {
	//Load a mxtsession file from a path or from stdin if path is "".
	return loadFile(path)
}
