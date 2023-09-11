package internal

import (
	"regexp"
	"testing"
)

func TestRemoveSomeValuesFromSliceByRegExp(t *testing.T) {
	result := RemoveSomeValuesFromSliceByRegExp(
		[]string{
			"-U DatabaseUser",
			"-h=127.0.0.1",
			"-p 5432",
			"-d=DatabaseName",
			"-f dump.zip",
			"--data-only",
			"-B",
		},
		regexp.MustCompile(`^(-U|-h|-p|-d|-f)(=|\s)?[\S*]?`),
	)

	expected := []string{"--data-only", "-B"}

	if len(expected) != len(result) {
		t.Errorf("Expected: %s\nResult: %s", expected, result)
	}
}
