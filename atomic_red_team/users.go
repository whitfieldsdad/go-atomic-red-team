package atomic_red_team

import (
	"fmt"
	"os/user"
	"runtime"
)

var (
	CurrentUser = GetCurrentUser()
)

var (
	DefaultRegularUser = CurrentUser
	DefaultSudoUser    = GetDefaultSudoUser()
)

func GetCurrentUser() string {
	user, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Cannot lookup current user: %s", err.Error()))
	}
	return user.Username
}

func GetDefaultSudoUser() string {
	if runtime.GOOS == "windows" {
		return "Administrator"
	}
	return "root"
}
