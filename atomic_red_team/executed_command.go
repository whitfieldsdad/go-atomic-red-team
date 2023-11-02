package atomic_red_team

import "time"

type ExecutedCommand struct {
	Id      string    `json:"id"`
	Time    time.Time `json:"time"`
	Command string    `json:"command"`
	Process *Process  `json:"process"`
	Error   string    `json:"error,omitempty"`
}

func (ec ExecutedCommand) GetArtifacts() []Artifact {
	return []Artifact{ec.Process}
}
