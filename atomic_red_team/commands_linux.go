//go:build !windows && !js && !darwin
// +build !windows,!js,!darwin

package atomic_red_team

import "syscall"

func getSysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
