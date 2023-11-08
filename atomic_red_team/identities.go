package atomic_red_team

import (
	"github.com/denisbrodbeck/machineid"
)

const (
	AppId = "b0227e37-e096-4963-ac1a-e07cac092d4c"
)

type EndpointIdentities struct {
	HostId string `json:"host_id" yaml:"host_id"`
}

func GetEndpointIdentities() (*EndpointIdentities, error) {
	hostId, err := getHostId()
	if err != nil {
		return nil, err
	}
	return &EndpointIdentities{
		HostId: hostId,
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
