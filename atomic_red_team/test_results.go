package atomic_red_team

import (
	"errors"
	"time"
)

type TestResult struct {
	Id               string                       `json:"id" yaml:"id"`
	Time             time.Time                    `json:"time" yaml:"time"`
	Test             Test                         `json:"test" yaml:"test"`
	ExecutedCommands []ExecutedCommand            `json:"executed_commands" yaml:"executed_commands"`
	Dependencies     []DependencyResolutionResult `json:"dependencies,omitempty" yaml:"dependencies"`
}

func NewTestResult(testId string, test Test, executedCommands []ExecutedCommand) (*TestResult, error) {
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
		Id:               NewUUID4(),
		Time:             *startTime,
		Test:             test,
		ExecutedCommands: executedCommands,
	}, nil
}

func (result TestResult) GetProcesses() []Process {
	var processes []Process
	for _, executedCommand := range result.ExecutedCommands {
		processes = append(processes, executedCommand.GetProcesses()...)
	}
	return processes
}

func (result TestResult) GetCommands() []Command {
	var commands []Command
	for _, executedCommand := range result.ExecutedCommands {
		commands = append(commands, executedCommand.Command)
	}
	return commands
}
