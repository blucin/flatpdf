package pdf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func validatePDF(filePath string) error {
	if filepath.Ext(filePath) != ".pdf" {
		return fmt.Errorf("file '%s' is not a pdf", filePath)
	}
	if !fileExists(filePath) {
		return fmt.Errorf("pdf file '%s' does not exist", filePath)
	}
	return nil
}

type FlattenPDFOptions struct {
	ImageDensity int // recommended: 200-300 (low) or 600 (medium) or 1200 (high)
	ImageQuality int // recommended: 0-100
	OutputDir    string
}

// FlattenPDF takes a slice of PDF file paths, validates them and
// flattens them to make them read-only.
// It returns the paths of the flattened PDF files.
func FlattenPDF(filePaths []string, outputDir string, opts FlattenPDFOptions) ([]string, error) {
	for _, filePath := range filePaths {
		err := validatePDF(filePath)
		if err != nil {
			return nil, err
		}
	}

	if opts.ImageDensity == 0 {
		opts.ImageDensity = 600
	}

	if opts.ImageQuality > 100 || opts.ImageQuality < 0 {
		return nil, fmt.Errorf("image quality must be between 0 and 100")
	}

	if opts.OutputDir == "" {
		opts.OutputDir = "."
	}

	// TODO: detect pdf paper size (A4, A5, etc.)

	var outputFiles = make([]string, 0)

	// magick wand init
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.SetOption("density", fmt.Sprintf("%d", opts.ImageDensity))
	mw.SharpenImage(0, 0)
	mw.SetImageFormat("jpg")
	mw.SetImageCompressionQuality(uint(opts.ImageQuality))

	for i, filePathStr := range filePaths {
		out, err := Flatten(mw, i, filePathStr)
		if err != nil {
			return nil, err
		} else {
			outputFiles = append(outputFiles, out)
		}
	}

	return outputFiles, nil
}
