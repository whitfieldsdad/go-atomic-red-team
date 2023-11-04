//go:build windows
// +build windows

package atomic_red_team

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

func isElevated() (bool, error) {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)

	if err != nil {
		return false, errors.Wrap(err, "failed to allocate and initialize SID")
	}
	token := windows.Token(0)
	defer token.Close()

	member, err := token.IsMember(sid)
	if err != nil {
		return false, errors.Wrap(err, "failed to check token membership")
	}
	if member {
		return true, nil
	}
	elevated := token.IsElevated()
	return elevated, nil
}
