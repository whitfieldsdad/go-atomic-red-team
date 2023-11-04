package atomic_red_team

import (
	"github.com/google/uuid"
)

// NewUUID4 returns a type 4 UUID.
func NewUUID4() string {
	return uuid.New().String()
}
