package atomic_red_team

import "os"

var (
	DefaultAtomicsDir = os.ExpandEnv("$ATOMICS_DIR")
)

type TestOptions struct {
	InputArguments map[string]interface{} `json:"input_arguments"`
}

func NewTestOptions(atomicsDir string) *TestOptions {
	if atomicsDir == "" {
		atomicsDir = DefaultAtomicsDir
	}
	return &TestOptions{
		InputArguments: make(map[string]interface{}),
	}
}
