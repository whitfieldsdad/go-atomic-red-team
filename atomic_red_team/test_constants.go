package atomic_red_team

import (
	"os"
	"path/filepath"
)

var (
	TestTempDir = filepath.Join(os.TempDir(), "791ca085-0b77-465c-ac9f-3fe69bf222ab")
)
