package utils

import (
	"bufio"
	"os"
)

// ReadFirstLineFromFile returns the firstline from a txt file path
func ReadFirstLineFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	return firstLine, nil
}
