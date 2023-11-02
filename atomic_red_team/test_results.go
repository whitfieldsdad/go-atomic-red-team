package atomic_red_team

import (
	"errors"
	"time"
)

type TestResult struct {
	Id               string            `json:"id"`
	Time             time.Time         `json:"time"`
	TestId           string            `json:"test_id"`
	ExecutedCommands []ExecutedCommand `json:"executed_commands"`
	Error            string            `json:"error,omitempty"`
}

func NewTestResult(testId string) (*TestResult, error) {
	if testId == "" {
		return nil, errors.New("missing test ID")
	}
	return &TestResult{
		Id:     NewUUID4(),
		TestId: testId,
	}, nil
}
