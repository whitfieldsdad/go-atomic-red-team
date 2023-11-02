package atomic_red_team

import "os"

var (
	DefaultAtomicsDir = os.ExpandEnv("$ATOMICS_DIR")
)
