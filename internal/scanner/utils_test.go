package scanner

import (
	"testing"
)

func TestCreateCacheBusterURL(t *testing.T) {
	s := &ScannerArgs{URL: "https://example.com/?removeme=123"}
	s.SetCacheBusterURL()
	t.Logf(s.cacheBusterURL)
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]string{"key1": "value1", "key2": "value2"}
	map2 := map[string]string{"key2": "value3", "key3": "value4"}
	mergedMap := MergeMaps(map1, map2)
	t.Logf("Final Map: %v", mergedMap)
}
