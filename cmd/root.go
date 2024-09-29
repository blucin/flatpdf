package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flatpdf",
	Short: "A pdf flattener to make them read-only",
	Long: `flatpdf is a tool to flatten PDF files, making them read-only.

You can pass multiple PDF files as input. By default, the flattened PDF files 
will be saved in the same directory as the originals, with the suffix '_flat' added 
to the file name. 

Alternatively, you can use the '-o' flag to specify an output directory where 
all flattened PDFs will be saved. Note that if a file with the same name already exists 
in the output directory, it will be overwritten.

Use the '-h' flag to view all available options.
`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of flatpdf",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flatpdf v0.1")
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
