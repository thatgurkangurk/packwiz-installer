package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thatgurkangurk/packwiz-installer/pkg/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit:  %s\nBuilt:   %s\n", build.Version, build.Commit, build.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
