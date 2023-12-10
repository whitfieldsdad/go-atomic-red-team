package atomic

import (
	"errors"
	"time"

	"github.com/whitfieldsdad/go-building-blocks/pkg/bb"
)

type TestResult struct {
	Id               string                       `json:"id" yaml:"id"`
	Time             time.Time                    `json:"time" yaml:"time"`
	Test             Test                         `json:"test" yaml:"test"`
	ExecutedCommands []bb.ExecutedCommand         `json:"executed_commands" yaml:"executed_commands"`
	Dependencies     []DependencyResolutionResult `json:"dependencies,omitempty" yaml:"dependencies"`
}

func NewTestResult(testId string, test Test, executedCommands []bb.ExecutedCommand) (*TestResult, error) {
	if testId == "" {
		return nil, errors.New("missing test ID")
	}
	if executedCommands == nil {
		return nil, errors.New("missing executed commands")
	}
	var startTime *time.Time
	for _, executedCommand := range executedCommands {
		if startTime == nil || executedCommand.StartTime.Before(*startTime) {
			startTime = &executedCommand.StartTime
		}
	}
	return &TestResult{
		Id:               bb.NewUUID4(),
		Time:             *startTime,
		Test:             test,
		ExecutedCommands: executedCommands,
	}, nil
}

func (result TestResult) GetProcesses() []bb.Process {
	var processes []bb.Process
	for _, executedCommand := range result.ExecutedCommands {
		processes = append(processes, executedCommand.GetProcesses()...)
	}
	return processes
}

func (result TestResult) GetCommands() []bb.Command {
	var commands []bb.Command
	for _, executedCommand := range result.ExecutedCommands {
		commands = append(commands, executedCommand.Command)
	}
	return commands
}
