package atomic_red_team

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
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

type CommandOptions struct {
	IncludeParentProcesses bool `json:"include_parent_processes"`
}

func NewCommandOptions() *CommandOptions {
	return &CommandOptions{}
}

type ExecutedCommand struct {
	Id        string    `json:"id"`
	Time      time.Time `json:"time"`
	Command   *Command  `json:"command"`
	Processes []Process `json:"processes"`
}

func ExecuteCommand(ctx context.Context, command, commandType string, opts *CommandOptions) (*ExecutedCommand, error) {
	now := time.Now()
	executedCommand := &ExecutedCommand{
		Id:   NewUUID4(),
		Time: now,
		Command: &Command{
			Command:     command,
			CommandType: commandType,
		},
	}
	if opts == nil {
		opts = NewCommandOptions()
	}
	subprocess, err := executeCommand(ctx, command, commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}

	var processes []Process
	if opts.IncludeParentProcesses {
		ancestors, err := GetProcessAncestors(subprocess.PID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get process ancestors")
		}
		processes = append(processes, ancestors...)
	}
	processes = append(processes, *subprocess)
	executedCommand.Processes = processes
	return executedCommand, nil
}

func executeCommand(ctx context.Context, command, commandType string) (*Process, error) {
	var err error
	pid := os.Getpid()
	process, err := GetProcess(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get process")
	}
	process.User, err = GetUser()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	argv, err := prepareCommand(command, DefaultCommandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare command")
	}
	cmd, err := executeArgv(ctx, argv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command")
	}
	exitCode := cmd.ProcessState.ExitCode()
	process.ExitCode = &exitCode
	process.Stdout = cmd.Stdout.(*bytes.Buffer).String()
	process.Stderr = cmd.Stderr.(*bytes.Buffer).String()
	return process, nil
}

func executeArgv(ctx context.Context, argv []string) (*exec.Cmd, error) {
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

	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}
	pid := cmd.Process.Pid
	ppid := os.Getpid()
	log.Debugf("Executing command: %s %s (PID: %d, PPID: %d)", path, strings.Join(argv, " "), pid, ppid)

	log.Debugf("Waiting for command to exit (PID: %d, PPID: %d)", pid, ppid)
	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for command to exit")
	}
	log.Debugf("Command exited (PID: %d, PPID: %d, exit code: %d)", pid, ppid, cmd.ProcessState.ExitCode())
	return cmd, nil
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
	commandTypes = []string{PowerShell, CommandPrompt, Sh, Bash}
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

func prepareCommand(command, commandType string) ([]string, error) {
	commandType, err := parseCommandType(commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command type")
	}
	shim, ok := commandShims[commandType]
	if !ok {
		return nil, errors.Errorf("invalid command type: %s", commandType)
	}
	return append(shim, command), nil
}

func parseCommandType(commandType string) (string, error) {
	for _, validCommandType := range commandTypes {
		if commandType == validCommandType {
			return commandType, nil
		}
	}
	return "", errors.Errorf("invalid command type: %s", commandType)
}
