package internal

import (
	"testing"
	"time"
)

func TestCreateDumpFilenameFromPatternAndDate(t *testing.T) {
	date := time.Date(2023, time.August, 29, 12, 0, 0, 0, time.UTC)

	result := CreateDumpFilenameFromPatternAndDate("dump-{yyyy}{MM}{dd}-{hh}{mm}{ss}", date)
	expected := "dump-20230829-120000"

	if result != expected {
		t.Errorf("Expected: %s\nResult: %s", expected, result)
	}

	result = CreateDumpFilenameFromPatternAndDate("backup-{yy}{M}{d}-{hh}{m}{s}", date)
	expected = "backup-23829-1200"

	if result != expected {
		t.Errorf("Expected: %s\nResult: %s", expected, result)
	}
}

func TestFindAllCommandOptionsFromString(t *testing.T) {
	result := FindAllCommandOptionsFromString("-B --only-data    -u me")
	expected := []string{"-B", "--only-data", "-u me"}

	lenIsValid := len(result) == len(expected)
	result0IsValid := result[0] == expected[0]
	result1IsValid := result[1] == expected[1]
	result2IsValid := result[2] == expected[2]

	if !lenIsValid || !result0IsValid || !result1IsValid || !result2IsValid {
		t.Errorf("Expected: %s\nResult: %s", expected, result)
	}
}
