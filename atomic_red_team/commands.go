package atomic_red_team

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Command struct {
	Command     string `json:"command"`
	CommandType string `json:"command_type"`
}

func NewCommand(command, commandType string) (*Command, error) {
	return &Command{
		Command:     command,
		CommandType: commandType,
	}, nil
}

func (command Command) Execute(ctx context.Context, opts *CommandOptions) (*ExecutedCommand, error) {
	return ExecuteCommand(ctx, command.Command, command.CommandType, opts)
}

type ExecutedCommand struct {
	Id               string    `json:"id"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Command          Command   `json:"command"`
	ExitCode         int       `json:"exit_code"`
	Subprocess       Process   `json:"subprocess"`
	RelatedProcesses []Process `json:"related_processes"`
}

func (result ExecutedCommand) GetProcesses() []Process {
	var processes []Process
	processes = append(processes, result.Subprocess)
	processes = append(processes, result.RelatedProcesses...)
	return processes
}

func (result ExecutedCommand) GetDuration() time.Duration {
	return result.EndTime.Sub(result.StartTime)
}

func ExecuteCommand(ctx context.Context, command, commandType string, opts *CommandOptions) (*ExecutedCommand, error) {
	if opts == nil {
		opts = NewCommandOptions()
	}
	argv, err := wrapCommand(command, commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap command")
	}
	startTime := time.Now()
	subprocess, err := executeArgv(ctx, argv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}
	endTime := time.Now()

	var relatedProcesses []Process
	if opts.IncludeParentProcesses {
		relatedProcesses, err = GetProcessAncestors(subprocess.PID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get related processes")
		}
	}
	executedCommand := &ExecutedCommand{
		Id:               NewUUID4(),
		StartTime:        startTime,
		EndTime:          endTime,
		Command:          Command{Command: command, CommandType: commandType},
		ExitCode:         *subprocess.ExitCode,
		Subprocess:       *subprocess,
		RelatedProcesses: relatedProcesses,
	}
	return executedCommand, nil
}

func executeArgv(ctx context.Context, argv []string) (*Process, error) {
	path, err := exec.LookPath(argv[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to find command")
	}
	cmd := exec.Command(path, argv[1:]...)
	cmd.SysProcAttr = getSysProcAttrs()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Execute the command.
	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	// Collect information about the subprocess.
	pid := cmd.Process.Pid
	process, err := GetProcess(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to collect process metadata")
	}
	if process.Argv == nil || process.CommandLine == "" {
		process.Argv = argv
		process.CommandLine = strings.Join(argv, " ")
	}

	// Wait for the command to complete.
	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for command to exit")
	}
	process.Stdout = stdout.String()
	process.Stderr = stderr.String()

	exitCode := cmd.ProcessState.ExitCode()
	process.ExitCode = &exitCode
	return process, nil
}

var (
	WindowsPowerShell = "powershell"
	PowerShellCore    = "pwsh"
	PowerShell        = getPowerShellCommandType()
	CommandPrompt     = "command_prompt"
	Sh                = "sh"
	Bash              = "bash"
)

var (
	commandShims = map[string][]string{
		WindowsPowerShell: {"powershell", "-ExecutionPolicy", "Bypass", "-Command"},
		PowerShellCore:    {"pwsh", "-Command"},
		CommandPrompt:     {"cmd", "/c"},
		Sh:                {"sh", "-c"},
		Bash:              {"bash", "-c"},
	}
)

var (
	DefaultCommandType = getDefaultCommandType()
)

func getDefaultCommandType() string {
	if runtime.GOOS == "windows" {
		return CommandPrompt
	}
	return Bash
}

func getPowerShellCommandType() string {
	if runtime.GOOS == "windows" {
		return WindowsPowerShell
	}
	return PowerShellCore
}

func wrapCommand(command, commandType string) ([]string, error) {
	argv := commandShims[commandType]
	if argv == nil {
		return nil, errors.Errorf("invalid command type: %s", commandType)
	}
	argv = append(argv, command)
	return argv, nil
}
