package pdf

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidatePDFPaths validates the given paths to ensure they are valid PDF paths
func ValidatePDFPaths(paths []string) error {
	for _, p := range paths {
		if filepath.Ext(p) != ".pdf" {
			return fmt.Errorf("%s is not a valid pdf path", p)
		}
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return fmt.Errorf("file at %s path does not exist", p)
		}
	}
	return nil
}

// Exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Deletes files at the given paths
func DeleteFiles(paths []string) error {
    for _, p := range paths {
        err := os.Remove(p)
        if err != nil {
            return err
        }
    }
    return nil
}
