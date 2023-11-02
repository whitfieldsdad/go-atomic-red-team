package atomic_red_team

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

type Command struct {
	Command     string `json:"command,omitempty"`
	CommandType string `json:"command_type,omitempty"`
}

func NewCommand(command, commandType string) (*Command, error) {
	commandType, err := parseCommandType(commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command type")
	}
	return &Command{
		Command:     command,
		CommandType: commandType,
	}, nil
}

func (c *Command) Execute(ctx context.Context) (*ExecutedCommand, error) {
	return ExecuteCommand(ctx, c.Command, c.CommandType)
}

func ExecuteCommand(ctx context.Context, command, commandType string) (*ExecutedCommand, error) {
	argv, err := prepareCommand(command, commandType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap command")
	}
	return ExecuteCommandArgs(ctx, argv)
}

func ExecuteCommandArgs(ctx context.Context, argv []string) (*ExecutedCommand, error) {
	startTime := time.Now()
	process, err := executeCommand(ctx, argv)
	var errorString string
	if err != nil {
		errorString = err.Error()
	}
	if process != nil {
		startTime = process.Time
	}
	executedCommand := &ExecutedCommand{
		Id:      NewUUID4(),
		Time:    startTime,
		Command: strings.Join(argv, " "),
		Process: process,
		Error:   errorString,
	}
	return executedCommand, nil
}

func executeCommand(ctx context.Context, argv []string) (*Process, error) {
	commandLine := strings.Join(argv, " ")
	log.Infof("Executing command: %s", commandLine)

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.SysProcAttr = getSysProcAttrs()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	// Collect basic information about the subprocess.
	pid := cmd.Process.Pid
	ppid := os.Getpid()

	log.Infof("Started process (PID: %d, PPID: %d)", pid, ppid)

	process, err := GetProcess(pid)
	if err != nil {
		process = &Process{
			Time:        time.Now(),
			PID:         pid,
			PPID:        ppid,
			User:        CurrentUser,
			Argv:        argv,
			CommandLine: commandLine,
		}
	}

	// Wait for the process to exit.
	log.Infof("Waiting for process to exit (PID: %d, PPID: %d)", pid, ppid)
	err = cmd.Wait()
	if err != nil {
		log.Warnf("Failed to wait for command to exit (PID: %d, PPID: %d): %s", pid, ppid, err)
		return nil, errors.Wrap(err, "failed to wait for command to exit")
	}
	exitCode := cmd.ProcessState.ExitCode()
	log.Infof("Process exited (PID: %d, PPID: %d, exit code: %d)", pid, ppid, exitCode)

	// Update the process with the exit code and duration.
	now := time.Now()
	process.ExitCode = &exitCode
	process.ExitTime = &now

	// Update the process with the stdout and stderr.
	process.Stdout = stdout.String()
	process.Stderr = stderr.String()
	return process, nil
}
