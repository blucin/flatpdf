package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flatpdf",
	Short: "A pdf flattener to make them read-only",
	Long: `flatpdf is a pdf flattener to make them read-only.
Pass pdf files as arguments to flat, new pdf files will be generated
with the same name but with the suffix '_flat'. Use the -h flag to 
see all available options.`,
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
