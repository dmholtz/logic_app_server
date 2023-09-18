package server

import (
	"fmt"
	"os"
)

// ReadQueryFile reads a SQL file and returns the contents as a string
func ReadQueryFile(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %s", filename, err))
	}
	return string(bytes)
}
