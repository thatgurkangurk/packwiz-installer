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
