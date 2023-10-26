package pdf

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func runFlatten(t *testing.T) (string, error) {
	t.Log("\nMake sure that projectPath/assets/testFileForFlat.pdf exists")

	current, err := os.Getwd()
	parentPath := filepath.Dir(current)
	if err != nil {
		return "", err
	}

	// magick wand init
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.SetOption("density", "1000")
	mw.SharpenImage(0, 0)
	mw.SetImageFormat("jpg")
	mw.SetImageCompressionQuality(99)

	got, err := Flatten(mw, 0, parentPath+"/assets/testFileForFlat.pdf")

	if err != nil {
		return got, err
	}

	return got, nil
}

func TestFlattenReturnsValidPath(t *testing.T) {
	got, err := runFlatten(t)
	if err != nil {
		t.Error(err)
	}
	// check if file is a pdf
	if filepath.Ext(got) != ".pdf" {
		t.Errorf("Flatten() should return a .pdf file path, got %s", got)
	}
	// cleanup
	if err := os.Remove(got); err != nil {
		t.Error(err)
	}
}

func TestFlattenCreatesFile(t *testing.T) {
	got, err := runFlatten(t)
	t.Log(got)
	if err != nil {
		t.Error(err)
	}
	// check if file exists
	if _, err := os.Stat(got); errors.Is(err, os.ErrNotExist) {
		t.Error("Flatten() did not create a file at expected path")
	}

	// cleanup
	if err := os.Remove(got); err != nil {
		t.Error(err)
	}

}
