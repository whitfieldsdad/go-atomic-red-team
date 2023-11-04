package atomic_red_team

import (
	"os/user"
	"time"

	"github.com/pkg/errors"
)

type User struct {
	Id       string    `json:"id"`
	Time     time.Time `json:"time"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	UID      string    `json:"uid"`
	GID      string    `json:"gid"`
	GroupIds []string  `json:"group_ids"`
	HomeDir  string    `json:"home_dir"`
}

func GetUser() (*User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup user")
	}
	gids, err := u.GroupIds()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup GIDs")
	}
	return &User{
		Id:       NewUUID4(),
		Time:     time.Now(),
		Username: u.Username,
		Name:     u.Name,
		UID:      u.Uid,
		GID:      u.Gid,
		HomeDir:  u.HomeDir,
		GroupIds: gids,
	}, nil
}
