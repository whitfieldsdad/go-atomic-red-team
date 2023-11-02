package atomic_red_team

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
)

// NewUUID4 returns a type 4 UUID.
func NewUUID4() string {
	return uuid.New().String()
}

// NewUUID5 returns a type 5 UUID.
func NewUUID5(data map[string]interface{}) string {
	blob, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON data: %s\n", err)
	}
	canonicalBlob, err := jcs.Transform(blob)
	if err != nil {
		log.Fatalf("Failed to canonicalize JSON data: %s\n", err)
	}
	return uuid.NewSHA1(uuid.Nil, canonicalBlob).String()
}
