package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

// responseChangeTypeToString maps ResponseChangeType values to string representations.
var responseChangeTypeToString = map[ResponseChangeType]string{
	ChangedLocationHeader: "ChangedLocationHeader",
	ChangedStatusCode:     "ChangedStatusCode",
	ChangedBody:           "ChangedBody",
	NoChange:              "NoChange",
}

// MarshalJSON ensures that ResponseChangeType is serialized as a string.
func (r ResponseChangeType) MarshalJSON() ([]byte, error) {
	if str, exists := responseChangeTypeToString[r]; exists {
		return json.Marshal(str)
	}
	return json.Marshal("Unknown")
}

// MarshalScannerOutput converts the ScannerOutput struct to JSON
func MarshalScannerOutput(scanResult ScannerOutput, outputFile string) ([]byte, error) {
	// Convert struct to JSON
	jsonData, err := json.MarshalIndent(scanResult, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonData, nil
}

// ExportJSONToFile appends the JSON data to a file
func ExportJSONToFile(jsonData []byte, outputFile string) error {
	// Write JSON data to file
	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
