package atomic_red_team

import (
	"errors"
	"time"
)

type TestResult struct {
	Id               string            `json:"id"`
	Time             time.Time         `json:"time"`
	Test             Test              `json:"test"`
	ExecutedCommands []ExecutedCommand `json:"executed_commands"`
}

func NewTestResult(testId string, test Test, executedCommands []ExecutedCommand) (*TestResult, error) {
	if testId == "" {
		return nil, errors.New("missing test ID")
	}
	if executedCommands == nil {
		return nil, errors.New("missing executed commands")
	}
	var startTime *time.Time
	for _, command := range executedCommands {
		if startTime == nil || command.Time.Before(*startTime) {
			startTime = &command.Time
		}
	}
	return &TestResult{
		Id:               NewUUID4(),
		Time:             *startTime,
		Test:             test,
		ExecutedCommands: executedCommands,
	}, nil
}
