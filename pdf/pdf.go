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

func worker(id int, mw *imagick.MagickWand, outputDir string, jobs <-chan string, results chan<- string) error {
	for job := range jobs {
		fmt.Printf("Worker %d started job %s\n", id, job)
		flattened, err := Flatten(mw, id, job, outputDir)
		if err != nil {
			return err
		}
		results <- flattened
	}
	return nil
}

type FlattenPDFOptions struct {
	ImageDensity int // recommended: 200-300 (low) or 600 (medium) or 1200 (high)
	ImageQuality int // recommended: 0-100
	Threads      int // default: 0
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

	if outputDir == "" {
		outputDir = "."
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

	// threading disabled
	if opts.Threads == 0 {
		for i, filePathStr := range filePaths {
			out, err := Flatten(mw, i, filePathStr, outputDir)
			if err != nil {
				return nil, err
			} else {
				outputFiles = append(outputFiles, out)
			}
		}
		return outputFiles, nil
	}

	// threading enabled
	jobs := make(chan string, len(filePaths))
	results := make(chan string, len(filePaths))

	for w := 1; w <= opts.Threads; w++ {
		go worker(w, mw, outputDir, jobs, results)
	}

	for j := 1; j <= len(filePaths); j++ {
		jobs <- filePaths[j-1]
	}
	close(jobs)

	for a := 1; a <= len(filePaths); a++ {
		out := <-results
		outputFiles = append(outputFiles, out)
	}
	return outputFiles, nil
}
