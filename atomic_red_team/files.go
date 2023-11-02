package atomic_red_team

import (
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"
)

type File struct {
	Id         string          `json:"id"`
	Time       time.Time       `json:"time"`
	Path       string          `json:"path"`
	Name       string          `json:"name"`
	Size       int64           `json:"size"`
	Hashes     *Hashes         `json:"hashes"`
	Timestamps *FileTimestamps `json:"timestamps"`
}

func (f File) GetArtifactType() ArtifactType {
	return FileArtifactType
}

// FileTimestamps contains the MACb timestamps of a file.
type FileTimestamps struct {
	ModifyTime *time.Time `json:"modify_time"`
	AccessTime *time.Time `json:"access_time"`
	ChangeTime *time.Time `json:"change_time"`
	BirthTime  *time.Time `json:"birth_time"`
}

func GetFile(path string, opts *FileOptions) (*File, error) {
	if opts == nil {
		opts = NewFileOptions()
	}
	file := &File{
		Id:   NewUUID4(),
		Time: time.Now(),
		Path: path,
		Name: filepath.Base(path),
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file.Size = info.Size()
	file.Hashes, _ = GetFileHashes(path, opts.HashingOptions)
	file.Timestamps, _ = GetFileTimestamps(path)
	return file, nil
}

// GetFileTimestamps returns the MACb timestamps of a file.
func GetFileTimestamps(path string) (*FileTimestamps, error) {
	st, err := times.Stat(path)
	if err != nil {
		return nil, err
	}
	m := st.ModTime()
	a := st.AccessTime()
	timestamps := &FileTimestamps{
		ModifyTime: &m,
		AccessTime: &a,
	}
	if st.HasChangeTime() {
		c := st.ChangeTime()
		timestamps.ChangeTime = &c
	}
	if st.HasBirthTime() {
		b := st.BirthTime()
		timestamps.BirthTime = &b
	}
	return timestamps, nil
}
