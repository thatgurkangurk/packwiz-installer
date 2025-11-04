// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func exactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == n {
			return nil
		}
		return fmt.Errorf(
			"%q requires exactly %d %s.\n",
			cmd.CommandPath(),
			n,
			pluralize("argument", n),
		)
	}
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
