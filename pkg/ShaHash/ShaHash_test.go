package ShaHash

import (
	"reflect"
	"testing"
)

func TestSha256Hash(t *testing.T) {
	input := []byte("Hello, World!")
	expectedHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"

	result := Sha256Hash(input)

	if !reflect.DeepEqual(result, expectedHash) {
		t.Errorf("Expected hash %s, but got %s", expectedHash, result)
	}
}
