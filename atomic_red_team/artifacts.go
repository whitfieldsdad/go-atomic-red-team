package atomic_red_team

type ArtifactType string

const (
	ProcessArtifactType         ArtifactType = "process"
	FileArtifactType            ArtifactType = "file"
)

// Artifact is the interface that all artifacts must implement
type Artifact interface {
	GetArtifactType() ArtifactType
}
