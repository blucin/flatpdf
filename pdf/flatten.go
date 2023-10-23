package pdf

import (
	"errors"
	"fmt"
	"path/filepath"
)

// Flatten takes a PDF file path and flattens it to make it read-only.
// It returns the path of the flattened PDF file.
func Flatten(filePathStr string) (string, error) {
	fileExtension := filepath.Ext(filePathStr)

	if fileExtension != ".pdf" {
		return "", errors.New("Flatten only works with .pdf files")
	}

	// TODO: Implement PDF flattening logic
	flattenedPath := fmt.Sprintf("%s_flat", filePathStr)
	return flattenedPath, nil
}
