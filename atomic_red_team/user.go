package atomic_red_team

import (
	"os/user"

	"github.com/pkg/errors"
)

type User struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	UID      string   `json:"uid"`
	GID      string   `json:"gid"`
	GroupIds []string `json:"group_ids"`
	HomeDir  string   `json:"home_dir"`
}

func GetCurrentUser() (*User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup current user")
	}
	return getUserInfo(u)
}

func GetUser(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup user")
	}
	return getUserInfo(u)
}

func getUserInfo(u *user.User) (*User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup user")
	}
	gids, err := u.GroupIds()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup GIDs")
	}
	return &User{
		Username: u.Username,
		Name:     u.Name,
		UID:      u.Uid,
		GID:      u.Gid,
		HomeDir:  u.HomeDir,
		GroupIds: gids,
	}, nil
}
