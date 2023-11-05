package atomic_red_team

import "time"

type Event struct {
	Id         string    `json:"id"`
	Time       time.Time `json:"time"`
	HostId     string    `json:"host_id"`
	UserId     string    `json:"user_id"`
	AgentId    string    `json:"agent_id"`
	ObjectId   string    `json:"object_id"`
	ObjectType string    `json:"object_type"`
	EventType  string    `json:"event_type"`
}

type StatusEvent struct {
	Event
	StatusType string `json:"status_type"`
}

func NewProcessEvent(hostId, userId, agentId, objectId, statusType string) StatusEvent {
	return StatusEvent{
		Event: Event{
			Id:         NewUUID4(),
			Time:       time.Now(),
			HostId:     hostId,
			UserId:     userId,
			AgentId:    agentId,
			ObjectId:   objectId,
			ObjectType: "process",
			EventType:  "status",
		},
		StatusType: statusType,
	}
}
