package cmd

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/atomic_red_team"
)

const (
	OutputFormatJson  = "json"
	OutputFormatYaml  = "yaml"
	OutputFormatPlain = "plain"
)

var (
	lineSeparator = strings.Repeat("-", 80)
)

var testsCmd = &cobra.Command{
	Use:   "tests",
	Short: "Tests",
}

var listTestsCmd = &cobra.Command{
	Use:   "list",
	Short: "List tests",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		outputFormat, _ := flags.GetString("output-format")
		tests, err := listTests(cmd.Flags())
		if err != nil {
			log.Errorf("Failed to list tests: %s", err)
			return
		}
		for _, test := range tests {
			printTest(test, outputFormat)
		}
	},
}

var countTestsCmd = &cobra.Command{
	Use:   "count",
	Short: "Count tests",
	Run: func(cmd *cobra.Command, args []string) {
		tests, err := listTests(cmd.Flags())
		if err != nil {
			log.Errorf("Failed to list tests: %s", err)
			return
		}
		fmt.Println(len(tests))
	},
}

var executeTestsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		outputFormat, _ := flags.GetString("output-format")
		atomicsDir := getAtomicsDir(flags)
		opts := getTestOptions(flags)

		tests, err := listTests(flags)
		if err != nil {
			log.Errorf("Failed to list tests: %s", err)
			return
		}
		ctx := context.Background()

		var results []atomic_red_team.TestResult
		for _, test := range tests {
			result, err := test.Run(ctx, atomicsDir, opts)
			if err != nil {
				log.Fatalf("Failed to execute test %s: %s", test.Name, err)
			}
			results = append(results, *result)
		}
		for _, result := range results {
			printTestResult(result, outputFormat)
		}
	},
}

var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Test dependencies",
}

var listDependenciesCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		outputFormat, _ := flags.GetString("output-format")
		tests, err := listTests(flags)
		if err != nil {
			log.Errorf("Failed to list tests: %s", err)
			return
		}
		for _, test := range tests {
			for _, dependency := range test.Dependencies {
				printTestDependency(test, dependency, outputFormat)
			}
		}
	},
}

// TODO
var countDependenciesCmd = &cobra.Command{
	Use:   "count",
	Short: "List dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		outputFormat, _ := flags.GetString("output-format")
		tests, err := listTests(flags)
		if err != nil {
			log.Errorf("Failed to list tests: %s", err)
			return
		}
		total := 0
		for _, test := range tests {
			total += len(test.Dependencies)
		}
		if outputFormat == OutputFormatPlain {
			fmt.Println(total)
		} else if outputFormat == OutputFormatJson || outputFormat == OutputFormatYaml {
			m := map[string]int{
				"total": total,
			}
			if outputFormat == OutputFormatJson {
				PrintJson(m)
			} else {
				PrintYaml(m)
			}
		} else {
			log.Fatalf("Unknown output format: %s", outputFormat)
		}
	},
}

// TODO
var checkDependenciesCmd = &cobra.Command{
	Use:   "check",
	Short: "Check dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		panic("not implemented")
	},
}

// TODO
var resolveDependenciesCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		panic("not implemented")
	},
}

func listTests(flags *pflag.FlagSet) ([]atomic_red_team.Test, error) {
	atomicsDir, _ := flags.GetString("atomics-dir")
	filter := getCommandLineFilter(flags)
	return atomic_red_team.ReadTests(atomicsDir, filter)
}

func getAtomicsDir(flags *pflag.FlagSet) string {
	atomicsDir, _ := flags.GetString("atomics-dir")
	if atomicsDir == "" {
		atomicsDir = atomic_red_team.DefaultAtomicsDir
	}
	return atomicsDir
}

func getTestOptions(flags *pflag.FlagSet) *atomic_red_team.TestOptions {
	opts := atomic_red_team.NewTestOptions()
	opts.CommandOptions.IncludeParentProcesses, _ = flags.GetBool("include-parent-processes")
	return opts
}

func getCommandLineFilter(flags *pflag.FlagSet) *atomic_red_team.TestFilter {
	f := &atomic_red_team.TestFilter{}
	f.Ids, _ = flags.GetStringSlice("id")
	f.Names, _ = flags.GetStringSlice("name")
	f.Descriptions, _ = flags.GetStringSlice("description")
	f.AttackTechniqueIds, _ = flags.GetStringSlice("attack-technique-id")
	f.ExecutorTypes, _ = flags.GetStringSlice("executor-type")
	f.ElevationRequired, _ = getNullableBool("elevation-required", flags)
	f.Platforms, _ = flags.GetStringSlice("platform")
	matchPlatform, _ := flags.GetBool("match-platform")
	if len(f.Platforms) == 0 && matchPlatform {
		f.Platforms = []string{runtime.GOOS}
	}
	return f
}

func printTest(test atomic_red_team.Test, outputFormat string) {
	if outputFormat == OutputFormatPlain {
		printTestPlain(test)
		fmt.Println(lineSeparator)
	} else if outputFormat == OutputFormatJson {
		PrintJson(test)
	} else if outputFormat == OutputFormatYaml {
		PrintYaml(test)
	} else {
		log.Fatalf("Unknown output format: %s", outputFormat)
	}
}

func printTestPlain(test atomic_red_team.Test) {
	fmt.Printf("ID: %s\n", test.AutoGeneratedGuid)
	fmt.Printf("Name: %s\n", test.Name)
	fmt.Printf("ATT&CK technique ID: %s\n", test.AttackTechniqueId)
	fmt.Printf("ATT&CK technique name: %s\n", test.AttackTechniqueName)
	fmt.Println()
	fmt.Printf("Description: %s\n", strings.TrimRight(test.Description, "\n"))
	fmt.Println()
	fmt.Printf("Supported platforms: %s\n", strings.Join(test.SupportedPlatforms, ","))
	fmt.Printf("Command type: %s\n", test.Executor.Name)
	fmt.Printf("Requires elevation: %v\n", test.Executor.ElevationRequired)
	fmt.Println()
	fmt.Printf("Commands:\n\n%s\n", strings.TrimRight(test.Executor.Command, "\n"))
	if test.Executor.CleanupCommand != "" {
		fmt.Println()
		fmt.Printf("Cleanup commands:\n\n%s\n", strings.TrimRight(test.Executor.CleanupCommand, "\n"))
	}
}

func printTestDependency(test atomic_red_team.Test, dependency atomic_red_team.Dependency, outputFormat string) {
	if outputFormat == OutputFormatPlain {
		printTestDependencyPlain(test, dependency)
		fmt.Println(lineSeparator)
	} else if outputFormat == OutputFormatJson {
		PrintJson(dependency)
	} else if outputFormat == OutputFormatYaml {
		PrintYaml(dependency)
	} else {
		log.Fatalf("Unknown output format: %s", outputFormat)
	}
}

func printTestDependencyPlain(test atomic_red_team.Test, dependency atomic_red_team.Dependency) {
	fmt.Printf("Test ID: %s\n", test.AutoGeneratedGuid)
	fmt.Printf("Test name: %s\n", test.Name)
	fmt.Println()
	fmt.Printf("Description: %s", strings.TrimRight(test.Description, "\n"))
}

func printTestResult(result atomic_red_team.TestResult, outputFormat string) {
	if outputFormat == OutputFormatPlain {
		printTestResultPlain(result)
		fmt.Println(lineSeparator)
	} else if outputFormat == OutputFormatJson {
		PrintJson(result)
	} else if outputFormat == OutputFormatYaml {
		PrintYaml(result)
	} else {
		log.Fatalf("Unknown output format: %s", outputFormat)
	}
}

func printTestResultPlain(result atomic_red_team.TestResult) {
	fmt.Printf("Test ID: %s\n", result.Test.AutoGeneratedGuid)
	fmt.Printf("Test result ID: %s\n", result.Id)
	fmt.Printf("Time: %s\n", result.Time.Format(time.RFC3339))
	fmt.Println()
	fmt.Printf("Executed commands:\n\n")
	for _, command := range result.ExecutedCommands {
		fmt.Printf("%s\n", strings.TrimRight(command.Command.Command, "\n"))
	}
	fmt.Println()
	fmt.Printf("Processes:\n\n")
	for _, command := range result.ExecutedCommands {
		for _, process := range command.Processes {
			fmt.Printf("- %d,%d\n", process.PID, process.PPID)
		}
	}
	fmt.Println()
	fmt.Printf("Executables:\n\n")
	var paths []string
	for _, command := range result.ExecutedCommands {
		for _, process := range command.Processes {
			if process.Executable != nil {
				path := process.Executable.Path
				if path != "" && !slices.Contains(paths, path) {
					paths = append(paths, path)
				}
			}
		}
	}
	for _, path := range paths {
		fmt.Printf("- %s\n", path)
	}
}

func init() {

	// Add commands.
	rootCmd.AddCommand(testsCmd)
	testsCmd.AddCommand(listTestsCmd, countTestsCmd, executeTestsCmd)

	testsCmd.AddCommand(dependenciesCmd)
	dependenciesCmd.AddCommand(listDependenciesCmd, countDependenciesCmd, checkDependenciesCmd, resolveDependenciesCmd)

	// Add flags.
	flagset := pflag.FlagSet{}
	flagset.StringP("atomics-dir", "", atomic_red_team.DefaultAtomicsDir, "Path to atomic-red-team/atomics directory")
	flagset.StringP("output-format", "o", OutputFormatPlain, "Output format")

	flagset.StringSliceP("id", "", []string{}, "Test IDs")
	flagset.StringSliceP("name", "", []string{}, "Test names")
	flagset.StringSliceP("description", "", []string{}, "Test descriptions")
	flagset.StringSliceP("attack-technique-id", "", []string{}, "ATT&CK technique IDs")
	flagset.StringSliceP("attack-technique-name", "", []string{}, "ATT&CK technique names")
	flagset.StringSliceP("platform", "", []string{}, "Platforms")
	flagset.StringSliceP("plan", "p", []string{}, "Test plans")
	flagset.StringSliceP("executor-type", "t", []string{}, "Executor types")
	flagset.BoolP("elevation-required", "", false, "Elevation required")
	flagset.BoolP("match-platform", "", false, "Match platform")

	// Pass the same flags to all commands.
	listTestsCmd.Flags().AddFlagSet(&flagset)
	countTestsCmd.Flags().AddFlagSet(&flagset)
	executeTestsCmd.Flags().AddFlagSet(&flagset)
	listDependenciesCmd.Flags().AddFlagSet(&flagset)
	countDependenciesCmd.Flags().AddFlagSet(&flagset)
	checkDependenciesCmd.Flags().AddFlagSet(&flagset)

	// Add some flags specific to executing tests.
	executeTestsCmd.Flags().BoolP("include-parent-processes", "", atomic_red_team.IncludeParentProcesses, "Include parent processes")
}
