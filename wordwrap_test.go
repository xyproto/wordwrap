package wordwrap

import (
	"reflect"
	"testing"
)

func TestWordWrap(t *testing.T) {
	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
	maxWidth := 40
	expected := []string{
		"Lorem ipsum dolor sit amet, consectetur",
		"adipiscing elit. Sed do eiusmod tempor",
		"incididunt ut labore et dolore magna",
		"aliqua.",
	}

	result, err := WordWrap(text, maxWidth)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
