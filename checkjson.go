package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func checkjson() {
	// Path to the JSON file
	jsonFilePath := "edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json"

	// Read the JSON file
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Validate JSON format
	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	fmt.Println("JSON is valid and correctly formatted.")
}
