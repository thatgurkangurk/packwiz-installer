package cmd

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thatgurkangurk/packwiz-installer/core"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install [flags] URL",
	Aliases: []string{"i"},
	Short:   "Install and update a packwiz modpack",
	Args:    exactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// args
		packUrl, err := url.ParseRequestURI(args[0])
		if err != nil {
			return fmt.Errorf("the install command requires URL of 'pack.toml'")
		}
		// flags
		var (
			hformat string
			hhash   string
		)
		if cmd.Flag("hash").Value.String() != "" {
			var ok bool
			hformat, hhash, ok = parseHashFlag(cmd.Flag("hash").Value.String())
			if !ok {
				return fmt.Errorf("invalid --hash format <HashFormat>:<Hash>")
			}
		}

		repo := core.NewRepository(packUrl, hformat, hhash)
		err = repo.Load(cmd.Context())
		if err != nil {
			return err
		}
		pack, err := core.NewPack(repo)
		if err != nil {
			return err
		}
		inst, err := core.NewLocalInstaller(pack, cmd.Flag("dir").Value.String())
		if err != nil {
			return err
		}

		fmt.Println("URL:", packUrl)
		fmt.Println("Dir:", inst.BaseDir)

		updates, err := inst.Install(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Println(updates.String())
		fmt.Println("Done.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().String("hash", "", `Hash of 'pack.toml' in the form of "<format>:<hash>" e.g. "sha256:abc012..."`)
	installCmd.Flags().StringP("dir", "d", ".", "Directory to install the modpack to")
}

func parseHashFlag(s string) (format string, hash string, ok bool) {
	if s == "" {
		return "", "", true
	}

	h := strings.Split(s, ":")
	if !(len(h) >= 2 && h[0] != "" && h[1] != "") {
		return "", "", false
	}

	format = h[0]
	hash = h[1]
	if !slices.Contains(core.PreferredHashList, format) {
		return "", "", false
	}

	return format, hash, true
}
