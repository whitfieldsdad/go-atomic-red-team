package atomic_red_team

import (
	"os/user"

	"github.com/pkg/errors"
)

type User struct {
	Name     string   `json:"name" yaml:"name"`
	Username string   `json:"username" yaml:"username"`
	UID      string   `json:"uid" yaml:"uid"`
	GID      string   `json:"gid" yaml:"gid"`
	GroupIds []string `json:"group_ids" yaml:"group_ids"`
	HomeDir  string   `json:"home_dir" yaml:"home_dir"`
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
