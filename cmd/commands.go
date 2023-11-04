package cmd

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/atomic_red_team"
)

var commandsCmd = &cobra.Command{
	Use:     "commands",
	Aliases: []string{"commands", "c"},
	Short:   "Commands",
}

var executeCommandCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"exec", "execute", "x"},
	Short:   "Run a command",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := strings.Join(args, " ")
		commandType, err := cmd.Flags().GetString("command-type")
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		c, err := atomic_red_team.NewCommand(command, commandType)
		if err != nil {
			log.Fatal(err)
		}
		opts := getCommandOptions(cmd.Flags())
		result, err := c.Execute(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}
		PrintJson(result)
	},
}

func getCommandOptions(flags *pflag.FlagSet) *atomic_red_team.CommandOptions {
	opts := &atomic_red_team.CommandOptions{}
	if flags.Changed("include-parent-processes") {
		opts.IncludeParentProcesses, _ = flags.GetBool("include-parent-processes")
	}
	return opts
}

func init() {
	rootCmd.AddCommand(commandsCmd)
	commandsCmd.AddCommand(executeCommandCmd)

	executeCommandCmd.Flags().StringP("command-type", "t", atomic_red_team.DefaultCommandType, "Command type")
	executeCommandCmd.Flags().BoolP("include-parent-processes", "o", atomic_red_team.IncludeParentProcesses, "Include parent processes")
}
