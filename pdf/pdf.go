package pdf

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var pool pdfium.Pool
var instance pdfium.Pdfium

type FlattenPDFOptions struct {
	InputFiles         []string // list of pdf files to process
	OutputDir          string   // optional directory to save processed files.
	ImageDensity       int      // recommended: 200-300 (low) or 600 (medium) or 1200 (high)
	ImageQuality       int      // recommended: 80-100, only for jpg encoding
	RequiresFlatSuffix bool     // whether to add _flat suffix to the output file name
}

// Renders a pdf file into images
func PdfToImages(filePath string, tempDir string, opts FlattenPDFOptions) (imagePaths []string, err error) {
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return nil, err
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return nil, err
	}

	// Create render requests for each page
	pageRequests := make([]requests.RenderPageInDPI, pageCount.PageCount)
	for i := 0; i < pageCount.PageCount; i++ {
		pageRequests[i] = requests.RenderPageInDPI{
			DPI: opts.ImageDensity,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    i,
				},
			},
		}
	}

	// Render the pages
	pagesRenders, err := instance.RenderPagesInDPI(&requests.RenderPagesInDPI{
		Pages: pageRequests,
	})
	if err != nil {
		return nil, err
	}
	defer pagesRenders.Cleanup()

	// Iterate over each rendered page
	singleImg := pagesRenders.Result.Image
	for i, page := range pagesRenders.Result.Pages {
		pageBounds := image.Rect(0, i*page.Height, page.Width, (i+1)*page.Height)
		pageImage := singleImg.SubImage(pageBounds).(*image.RGBA)
		savePath := filepath.Join(tempDir, fmt.Sprintf("page_%d.jpeg", i))
		f, err := os.Create(savePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		err = jpeg.Encode(f, pageImage, &jpeg.Options{Quality: opts.ImageQuality})
		if err != nil {
			return nil, err
		}
		imagePaths = append(imagePaths, savePath)
	}

	return imagePaths, nil
}

// Check pdfcpu for valid pageSize strings
func ImagesToPdf(imagePaths []string, outputPath string, pdfName string, pageSize string) (flattenedPath string, err error) {
	flattenedPath = filepath.Join(outputPath, pdfName)
	pdfConf := &pdfcpu.Import{
		PageDim:  types.PaperSize[pageSize],
		PageSize: pageSize,
		Pos:      types.Center,
		Scale:    1,
		InpUnit:  types.POINTS,
	}
	err = api.ImportImagesFile(imagePaths, flattenedPath, pdfConf, nil)
	if err != nil {
		return "", err
	}
	return flattenedPath, nil
}

// FlattenPDF takes a slice of PDF file paths, validates them and
// flattens them to make them read-only.
// It returns the paths of the flattened PDF files.
func FlattenPDF(opts FlattenPDFOptions) ([]string, error) {
	// pdfium init
	var err error
	poolConfig := webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	}
	pool, err = webassembly.Init(poolConfig)
	if err != nil {
		return nil, err
	}

	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		return nil, err
	}
	defer instance.Close()

	outputFiles := make([]string, 0)

    tempDir, err := os.MkdirTemp("", "flatpdf_temp")
	if err != nil {
		return nil, err
	}

	for _, inputPath := range opts.InputFiles {
        // pdf -> images
		imagePaths, err := PdfToImages(inputPath, tempDir, opts)
		if err != nil {
			return nil, err
		}

		// TODO: detect pdf paper size (A4, A5, etc.)
		pageSize := "A4"

		var flatPdfName string
        inputPdfName := filepath.Base(inputPath)[:len(filepath.Base(inputPath))-4]
		if opts.RequiresFlatSuffix {
			flatPdfName = inputPdfName + "_flat.pdf"
		} else {
			flatPdfName = inputPdfName + ".pdf"
		}

        // images -> pdf
		outputFile, err := ImagesToPdf(imagePaths, opts.OutputDir, flatPdfName, pageSize)
		outputFiles = append(outputFiles, outputFile)

        // clean up images
        DeleteFiles(imagePaths)
	}

    err = os.Remove(tempDir)
    if err != nil {
        return nil, err
    }

	return outputFiles, nil
}
