package main

import (
	"encoding/json"
	// "fmt"
	"io"

	// "log"
	"os"
)

// func validateJSON(jsonFilePath string) {
// 	// Open the JSON file
// 	jsonFile, err := os.Open(jsonFilePath)
// 	if err != nil {
// 		log.Fatalf("Error opening JSON file: %v", err)
// 	}
// 	defer jsonFile.Close()

// 	// Read the JSON file
// 	jsonBytes, err := io.ReadAll(jsonFile)
// 	if err != nil {
// 		log.Fatalf("Error reading JSON file: %v", err)
// 	}

// 	// Validate JSON format
// 	var jsonData map[string]interface{}
// 	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
// 		log.Fatalf("Error unmarshalling JSON: %v", err)
// 	}

// 	fmt.Println("JSON is valid and correctly formatted.")
// }

func validateJSON(jsonFilePath string) error {
	// Open the JSON file
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Read the JSON file
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	// Validate JSON format
	var parsedData interface{}
	if err := json.Unmarshal(jsonData, &parsedData); err != nil {
		return err
	}

	// Additional validation logic can be added here

	return nil
}
