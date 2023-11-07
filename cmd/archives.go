package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/atomic_red_team"
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
		inputPaths := args
		outputPath, _ := cmd.Flags().GetString("output-path")
		password, _ := cmd.Flags().GetString("password")

		var err error
		if password != "" {
			err = atomic_red_team.CreateEncryptedTarballFile(outputPath, inputPaths, password)
		} else {
			err = atomic_red_team.CreateTarballFile(outputPath, inputPaths)
		}
		if err != nil {
			log.Fatalf("Failed to create archive: %s", err)
		}
		log.Infof("Created archive: %s", outputPath)
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
