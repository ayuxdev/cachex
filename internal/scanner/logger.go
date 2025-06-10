// Description: This file contains the implementation of the logger for ScannerOutput.
package scanner

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// LogMode defines the mode of logging
type LogMode int

const (
	PrettyLog LogMode = iota
	JsonLog
)

type LogTarget int

const (
	StdoutLog LogTarget = 1 << iota
	FileLog
	BothLog = StdoutLog | FileLog
)

// Log prints the scan result based on the log mode and optionally writes it to a file.
func (so *ScannerOutput) Log(outputFilePath string, mode LogMode, logTarget LogTarget, skipTentative bool) error {
	if !so.IsVulnerable && (!so.IsResponseManipulable || skipTentative) {
		return nil
	}

	var message string
	var plainMessage string

	switch mode {
	case JsonLog:
		jsonData, err := json.Marshal(so)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		message = string(jsonData)
		plainMessage = message

	case PrettyLog:
		headerInfo := formatHeaderForLog(so.PayloadHeaders)

		url := color.CyanString(so.URL)
		header := color.MagentaString(headerInfo)

		if so.IsVulnerable {
			tag := color.New(color.FgRed, color.Bold).Sprint("[vuln]")
			poc := color.BlueString(so.PersistenceCheckResult.PoCLink)

			message = fmt.Sprintf(
				"%s [%s] [%s] [header: %s] [poc: %s]",
				tag,
				url,
				mapResponseChangeToVulnType(so.ManipulationType),
				header,
				poc,
			)

			plainMessage = fmt.Sprintf(
				"[vuln] [%s] [%s] [header: %s] [poc: %s]",
				so.URL,
				mapResponseChangeToVulnType(so.ManipulationType),
				headerInfo,
				so.PersistenceCheckResult.PoCLink,
			)

		} else if so.IsResponseManipulable && !skipTentative {
			tag := color.New(color.FgYellow, color.Bold).Sprint("[tentative-vuln]")

			message = fmt.Sprintf(
				"%s [%s] [%s] [header: %s]",
				tag,
				url,
				mapResponseChangeTypeToManipulationType(so.ManipulationType),
				header,
			)

			plainMessage = fmt.Sprintf(
				"[tentative-vuln] [%s] [%s] [header: %s]",
				so.URL,
				mapResponseChangeTypeToManipulationType(so.ManipulationType),
				headerInfo,
			)
		}
	}

	if logTarget&StdoutLog != 0 {
		fmt.Println(message)
	}
	if logTarget&FileLog != 0 {
		if outputFilePath == "" {
			return fmt.Errorf("output file path is required for file logging")
		}
		if err := writeToFile(outputFilePath, plainMessage); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}

// formatHeaderForLog converts a map of headers to a readable string for logging.
func formatHeaderForLog(headers map[string]string) string {
	if len(headers) == 0 {
		return "-"
	}

	if len(headers) == 1 {
		for k, v := range headers {
			return fmt.Sprintf("%s: %s", k, v)
		}
	}

	parts := make([]string, 0, len(headers))
	for k, v := range headers {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "; ")
}

// mapResponseChangeToVulnType maps the change type to a user-friendly vuln name.
func mapResponseChangeToVulnType(r ResponseChangeType) string {
	switch r {
	case ChangedLocationHeader:
		return "Location Poisoning"
	case ChangedStatusCode:
		return "Status Code Poisoning"
	case ChangedBody:
		return "Response Body Poisoning"
	}
	return "Unknown"
}

// mapResponseChangeTypeToManipulationType maps the change type to a user-friendly manipulation type.
func mapResponseChangeTypeToManipulationType(r ResponseChangeType) string {
	switch r {
	case ChangedLocationHeader:
		return "Location Header Manipulation"
	case ChangedStatusCode:
		return "Status Code Manipulation"
	case ChangedBody:
		return "Response Body Manipulation"
	}
	return "Unknown"
}
