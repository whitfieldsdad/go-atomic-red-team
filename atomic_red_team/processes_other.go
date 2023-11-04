//go:build !windows
// +build !windows

package atomic_red_team

import "os"

func isElevated() (bool, error) {
	return os.Geteuid() == 0, nil
}
