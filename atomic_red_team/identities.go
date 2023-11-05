package atomic_red_team

import (
	"log"
	"os/user"

	"github.com/denisbrodbeck/machineid"
)

const (
	AppId = "b0227e37-e096-4963-ac1a-e07cac092d4c"
)

func GetHostId() string {
	hostId, err := machineid.ProtectedID(AppId)
	if err != nil {
		log.Fatal(err)
	}
	m := map[string]string{
		"host_id": hostId,
	}
	return NewUUID5FromMap(AppId, m)
}

func GetUserId() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	m := map[string]string{
		"user_id": user.Uid,
	}
	return NewUUID5FromMap(AppId, m)
}

func GetAgentId() string {
	hostId := GetHostId()
	userId := GetUserId()
	return calculateAgentId(hostId, userId)
}

func calculateAgentId(hostId, userId string) string {
	m := map[string]string{
		"host_id": hostId,
		"user_id": userId,
	}
	return NewUUID5FromMap(AppId, m)
}
