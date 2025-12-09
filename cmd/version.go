package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of history-sanitizer",
	Long:  `All software has versions. This is history-sanitizer's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("history-sanitizer version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

