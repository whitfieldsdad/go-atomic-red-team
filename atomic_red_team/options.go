package atomic_red_team

import "os"

var (
	DefaultAtomicsDir = os.ExpandEnv("$ATOMICS_DIR")
)

const (
	IncludeParentProcesses = true
)

type TestOptions struct {
	InputArguments map[string]interface{} `json:"input_arguments" yaml:"input_arguments"`
	CommandOptions *CommandOptions        `json:"command_options" yaml:"command_options"`
}

func NewTestOptions() *TestOptions {
	return &TestOptions{
		InputArguments: make(map[string]interface{}),
		CommandOptions: NewCommandOptions(),
	}
}

type CommandOptions struct {
	IncludeParentProcesses bool `json:"include_parent_processes" yaml:"include_parent_processes"`
	IncludeCommandOutput   bool `json:"include_command_output" yaml:"include_command_output"`
}

func NewCommandOptions() *CommandOptions {
	return &CommandOptions{}
}
