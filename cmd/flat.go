package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(flatCmd)
}

var flatCmd = &cobra.Command{
	Use:   "flat [pdf files]",
	Short: "Flats the pdf files passed as arguments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flat called")
	},
}
