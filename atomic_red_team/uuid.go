package atomic_red_team

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
)

// NewUUID4 returns a type 4 UUID.
func NewUUID4() string {
	return uuid.New().String()
}

// NewUUID5 returns a type 5 UUID.
func NewUUID5(namespace string, blob []byte) string {
	return uuid.NewSHA1(uuid.MustParse(namespace), blob).String()
}

// NewUUID5FromMap returns a type 5 UUID from a map of contributing properties.
func NewUUID5FromMap(namespace string, m map[string]string) string {
	blob, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	canonicalizedBlob, err := jcs.Transform(blob)
	if err != nil {
		panic(err)
	}
	return NewUUID5(namespace, canonicalizedBlob)
}
