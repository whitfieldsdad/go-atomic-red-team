package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/atomic_red_team"
)

var processesCmd = &cobra.Command{
	Use:     "processes",
	Aliases: []string{"ps"},
	Short:   "Processes",
}

var listProcessesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List processes",
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := getProcesses(cmd.Flags())
		if err != nil {
			log.Fatalf("Failed to list processes: %s\n", err)
		}
		for _, process := range processes {
			PrintJson(process)
		}
	},
}

var countProcessesCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"n", "total"},
	Short:   "Count processes",
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := getProcesses(cmd.Flags())
		if err != nil {
			log.Fatalf("Failed to list processes: %s\n", err)
		}
		total := len(processes)
		fmt.Println(total)
	},
}

func getProcesses(flags *pflag.FlagSet) ([]atomic_red_team.Process, error) {
	pids, _ := flags.GetIntSlice("pid")
	ppids, _ := flags.GetIntSlice("ppid")
	names, _ := flags.GetStringSlice("name")
	executablePaths, _ := flags.GetStringSlice("path")
	executableNames, _ := flags.GetStringSlice("filename")
	executableHashes, _ := flags.GetStringSlice("hash")

	// Build the process filter.
	processFilter := &atomic_red_team.ProcessFilter{
		PIDs:             pids,
		PPIDs:            ppids,
		Names:            names,
		ExecutablePaths:  executablePaths,
		ExecutableNames:  executableNames,
		ExecutableHashes: executableHashes,
	}
	return atomic_red_team.GetProcesses(nil, processFilter)
}

func init() {
	rootCmd.AddCommand(processesCmd)

	flags := pflag.FlagSet{}
	flags.IntSliceP("pid", "", []int{}, "PIDs")
	flags.IntSliceP("ppid", "", []int{}, "PPIDs")
	flags.StringSliceP("name", "", []string{}, "Process names")
	flags.StringSliceP("path", "", []string{}, "Executable paths")
	flags.StringSliceP("filename", "", []string{}, "Executable filenames")
	flags.StringSliceP("hash", "", []string{}, "Executable hashes")

	listProcessesCmd.Flags().AddFlagSet(&flags)
	countProcessesCmd.Flags().AddFlagSet(&flags)
	processesCmd.AddCommand(listProcessesCmd, countProcessesCmd)
}
