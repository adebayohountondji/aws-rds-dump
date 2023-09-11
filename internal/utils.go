package internal

import (
	"compress/gzip"
	"io"
	"os"
	"regexp"
)

func CreateGzipFile(sourceFilePath string, targetFilePath string) (err error) {
	inputFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer func(inputFile *os.File) {
		_ = inputFile.Close()
	}(inputFile)

	outputFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	gzipWriter := gzip.NewWriter(outputFile)
	defer func(gzipWriter *gzip.Writer) {
		_ = gzipWriter.Close()
	}(gzipWriter)

	_, err = io.Copy(gzipWriter, inputFile)

	return err
}

func RemoveSomeValuesFromSliceByRegExp(slice []string, regexp *regexp.Regexp) []string {
	var result []string

	for _, value := range slice {
		if !regexp.MatchString(value) {
			result = append(result, value)
		}
	}

	return result
}
