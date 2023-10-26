package pdf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// Flatten takes a PDF file path & imagick MagickWand instance
// and flattens it to make it read-only.
// It returns the path of the flattened PDF file with _flat appended to the filename.
func Flatten(mw *imagick.MagickWand, index int, filePathStr string) (string, error) {
	fileExtension := filepath.Ext(filePathStr)

	if fileExtension != ".pdf" {
		return "", errors.New("flatten only works with .pdf files")
	}

	flattenedPath := fmt.Sprintf("%s_flat%s", filePathStr[:len(filePathStr)-len(fileExtension)], fileExtension)

	// convert PDF to images
	tempDir, err := os.MkdirTemp("", "flatpdf_temp")
	if err != nil {
		return "", errors.New("failed to create temp dir")
	}
	imagePaths, err := convertPDFToImages(mw, filePathStr, tempDir)
	fmt.Fprintf(os.Stderr, "imagePaths: %v\n", imagePaths)
	if err != nil {
		return "", errors.New("failed to convert pdf to images")
	}

	// convert images to PDF
	pdfConf := &pdfcpu.Import{
		PageDim:  types.PaperSize["A4"],
		PageSize: "A4",
		Pos:      types.Center,
		Scale:    1,
		InpUnit:  types.POINTS,
	}

	api.ImportImagesFile(imagePaths, flattenedPath, pdfConf, nil)

	// cleanup: delete temp image files and temp dir
	if err := os.RemoveAll(tempDir); err != nil {
		return "", errors.New("failed to delete temp dir")
	}

	return flattenedPath, nil
}

// convertPDFToImages converts a PDF file to a slice of image file paths.
func convertPDFToImages(mw *imagick.MagickWand, inputPath string, outputPath string) ([]string, error) {
	if err := mw.ReadImage(inputPath); err != nil {
		return nil, errors.New("failed to load pdf")
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_BACKGROUND); err != nil {
		return nil, errors.New("failed to set image alpha channel")
	}

	if err := mw.SetImageBackgroundColor(imagick.NewPixelWand()); err != nil {
		return nil, errors.New("failed to set image background color")
	}

	// convert pdf to images
	var imagePaths []string
	for i := uint(0); i < mw.GetNumberImages(); i++ {
		mw.SetIteratorIndex(int(i))
		tempName := fmt.Sprintf("flatpdf_temp_%d.jpg", i)
		tempPath := filepath.Join(outputPath, tempName)
		if err := mw.WriteImage(tempPath); err != nil {
			return nil, errors.New("failed to write images")
		}
		imagePaths = append(imagePaths, tempPath)
	}

	return imagePaths, nil
}
