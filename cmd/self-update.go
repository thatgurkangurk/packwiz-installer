package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thatgurkangurk/packwiz-installer/core"
)

// selfUpdateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "updates the cli",
	Run: func(cmd *cobra.Command, args []string) {
		core.Update()
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selfUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selfUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
