package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var archivesCmd = &cobra.Command{
	Use:     "archives",
	Aliases: []string{"archive"},
	Short:   "Archives",
}

var createArchiveCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tarball",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(archivesCmd)
	archivesCmd.AddCommand(createArchiveCmd)

	flagset := pflag.FlagSet{}
	flagset.StringP("output-path", "o", "", "Output path")
	flagset.StringP("password", "p", "", "Password")
	_ = createArchiveCmd.MarkFlagRequired("output-path")
}
