package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blucin/flatpdf/pdf"
	"github.com/spf13/cobra"
)

var inputFiles []string
var outputDir string
var imageDensity int
var imageQuality int
var requiresFlatSuffix bool

func init() {
	rootCmd.AddCommand(flatCmd)
	flatCmd.PersistentFlags().StringSliceVarP(&inputFiles, "input", "i", []string{}, "Input PDF file(s)")
	flatCmd.PersistentFlags().IntVarP(&imageDensity, "density", "d", 600, "image density")
	flatCmd.PersistentFlags().IntVarP(&imageQuality, "quality", "q", 99, "image quality (0-100) for jpeg encoding")
	flatCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory relative to the current path. Does not add _flat suffix on save so it will override any input files present in the output directory")
	flatCmd.MarkPersistentFlagRequired("input")
}

var flatCmd = &cobra.Command{
	Use:   "flat [pdf files]",
	Short: "Flats the pdf files passed as arguments",
	Run: func(cmd *cobra.Command, args []string) {
		currDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %s", err)
			return
		}
		if outputDir == "" || outputDir == "." {
			requiresFlatSuffix = true
		}
		outputDir = filepath.Join(currDir, outputDir)
		ok, err := pdf.Exists(outputDir)
		if err != nil {
			log.Fatalf("Error checking output directory: %s", err)
			return
		}
		if !ok {
			log.Fatalf("Output directory does not exist")
			return
		}
		if err := pdf.ValidatePDFPaths(inputFiles); err != nil {
			log.Fatalf("Error validating pdf paths: %s", err)
			return
		}

		outputFiles, err := pdf.FlattenPDF(pdf.FlattenPDFOptions{
			InputFiles:         inputFiles,
			OutputDir:          outputDir,
			ImageDensity:       imageDensity,
			ImageQuality:       imageQuality,
			RequiresFlatSuffix: requiresFlatSuffix,
		})
		if err != nil {
			log.Fatalf("Error flattening pdf: %s", err)
			return
		}
		for _, file := range outputFiles {
			log.Printf("File saved to: %s", file)
		}
	},
}
