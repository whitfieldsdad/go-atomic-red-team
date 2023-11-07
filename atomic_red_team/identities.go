package atomic_red_team

import (
	"os/user"

	"github.com/denisbrodbeck/machineid"
)

const (
	AppId = "b0227e37-e096-4963-ac1a-e07cac092d4c"
)

type EndpointIdentities struct {
	HostId  string `json:"host_id"`
	UserId  string `json:"user_id"`
	AgentId string `json:"agent_id"`
}

// GetEndpointIdentities returns the host ID, user ID, and agent ID for the current endpoint.
func GetEndpointIdentities() (*EndpointIdentities, error) {
	hostId, err := getHostId()
	if err != nil {
		return nil, err
	}
	userId, err := getUserId(hostId)
	if err != nil {
		return nil, err
	}
	agentId := calculateAgentId(hostId, userId)
	return &EndpointIdentities{
		HostId:  hostId,
		UserId:  userId,
		AgentId: agentId,
	}, nil
}

func getHostId() (string, error) {
	hostId, err := machineid.ProtectedID(AppId)
	if err != nil {
		return "", err
	}
	m := map[string]string{
		"host_id": hostId,
	}
	return NewUUID5FromMap(AppId, m), nil
}

func getUserId(hostId string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	m := map[string]string{
		"user_id": user.Uid,
		"host_id": hostId,
	}
	return NewUUID5FromMap(AppId, m), nil
}

func calculateAgentId(hostId, userId string) string {
	m := map[string]string{
		"host_id": hostId,
		"user_id": userId,
	}
	return NewUUID5FromMap(AppId, m)
}
