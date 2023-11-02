package atomic_red_team

import "time"

type Event struct {
	Id             string            `json:"id"`
	Time           time.Time         `json:"time"`
	ObjectType     string            `json:"object_type"`
	EventType      string            `json:"event_type"`
	CorrelationIds []string          `json:"correlation_ids,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
}

type ProcessEventType string

const (
	ProcessEventTypeStart ProcessEventType = "start"
	ProcessEventTypeExit  ProcessEventType = "exit"
)

type LightweightProcessEvent struct {
	Event
	PID  int `json:"pid"`
	PPID int `json:"ppid"`
}

type ProcessEvent struct {
	Event
	Process Process `json:"process"`
}
