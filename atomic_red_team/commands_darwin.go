//go:build darwin

package atomic_red_team

import "syscall"

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}