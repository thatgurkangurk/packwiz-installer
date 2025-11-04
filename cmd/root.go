// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thatgurkangurk/packwiz-installer/pkg/build"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "packwiz-installer",
	Short: "installs your packwiz modpacks automatically",
	Long: `this is a tool that allows you to automatically install/update your modpacks made with packwiz. (for clients and servers!)
	
hopefully this can replace packwiz-installer.jar and packwiz-installer-bootstrap.jar, and be faster.`,
	Version: build.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
