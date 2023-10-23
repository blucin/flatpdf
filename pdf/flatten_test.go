package pdf

import (
	"os"
	"path/filepath"
	"testing"
)

// Make sure that projectPath/assets/testFileForFlat.pdf before running tests

func TestFlattenReturnsValidPath(t *testing.T) {
	current, err := os.Getwd()
	parentPath := filepath.Dir(current)
	if err != nil {
		t.Error(err)
	}

	got, err := Flatten(parentPath + "/assets/testFileForFlat.pdf")
	if err != nil {
		t.Error(err)
	}

	// check if file is a pdf
	if filepath.Ext(got) != ".pdf" {
		t.Errorf("Flatten() should return a .pdf file path, got %s", got)
	}
}
