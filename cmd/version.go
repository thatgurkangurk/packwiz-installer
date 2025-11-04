// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
