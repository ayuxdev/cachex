package scanner

import (
	"fmt"
	"maps"
	"math/rand"
	urlparser "net/url"
	"os"
	"strings"
	"time"
)

// MergeMaps merges two maps and returns a new map with the contents of both maps
// If the same key is present in both maps, the value from map2 is used
func MergeMaps(map1, map2 map[string]string) map[string]string {
	// Create a new map to store the merged map
	mergedMap := make(map[string]string)

	// Copy the contents of map1 and map2 to mergedMap
	maps.Copy(mergedMap, map1)
	maps.Copy(mergedMap, map2)

	return mergedMap
}

// SetCacheBusterURL appends a cache-busting query parameter to the URL with a randomly generated value of the specified length.
func (s *ScannerArgs) SetCacheBusterURL() {
	cacheBusterValue := generateRandomString(5) // Todo: make this length configurable
	s.cacheBusterURL = createCacheBusterURL(s.URL, cacheBusterValue)
}

// createCacheBusterURL generates a URL with a cache-busting query parameter.
func createCacheBusterURL(url, paramValue string) string {
	return fmt.Sprintf("%s?cache=%s", removeURLQueryParams(url), paramValue)
}

// RemoveURLQueryParams removes the query parameters from the given URL
func removeURLQueryParams(url string) string {
	parsedURL, _ := urlparser.Parse(url) // ignoring errors
	parsedURL.RawQuery = ""

	return parsedURL.String()
}

// charset is the set of characters used to generate the random string
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a random string of the specified length
func generateRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano())) // seed the random number generator with the current time
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))] // select a random character from the charset
	}
	return string(result)
}

// writeToFile appends the given content to the specified file path
func writeToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	if !strings.HasSuffix(content, "\n") {
		content += "\n" // Ensure the content ends with a newline
	}
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}
