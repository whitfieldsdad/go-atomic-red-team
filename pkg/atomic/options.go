package atomic

import (
	"os"
)

var (
	DefaultAtomicsDir = os.ExpandEnv("$ATOMICS_DIR")
)

type TestOptions struct {
	InputArguments map[string]interface{} `json:"input_arguments" yaml:"input_arguments"`
}

func NewTestOptions() *TestOptions {
	return &TestOptions{
		InputArguments: make(map[string]interface{}),
	}
}
