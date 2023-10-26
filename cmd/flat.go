package cmd

import (
	"fmt"

	"github.com/blucin/flatpdf/pdf"
	"github.com/spf13/cobra"
)

var ImageDensity int
var ImageQuality int
var OutputDir string

func init() {
	rootCmd.AddCommand(flatCmd)
	flatCmd.PersistentFlags().IntVarP(&ImageDensity, "density", "d", 600, "image density")
	flatCmd.PersistentFlags().IntVarP(&ImageQuality, "quality", "q", 99, "image quality (0-100)")
	flatCmd.PersistentFlags().StringVarP(&OutputDir, "output", "o", "", "output directory")
}

var flatCmd = &cobra.Command{
	Use:   "flat [pdf files]",
	Short: "Flats the pdf files passed as arguments",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFiles, err := pdf.FlattenPDF(args, OutputDir, pdf.FlattenPDFOptions{
			ImageDensity: ImageDensity,
			ImageQuality: ImageQuality,
		})

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		for _, file := range outputFiles {
			fmt.Printf("Flattened file: %s\n", file)
		}
	},
}
