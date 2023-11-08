package atomic_red_team

import (
	"os"
	"path/filepath"
)

type File struct {
	Path     string  `json:"path" yaml:"path"`
	Filename string  `json:"filename" yaml:"filename"`
	Size     *int    `json:"size,omitempty" yaml:"size,omitempty"`
	Hashes   *Hashes `json:"hashes,omitempty" yaml:"hashes,omitempty"`
}

func NewFile(path string) *File {
	filename := filepath.Base(path)
	return &File{
		Path:     path,
		Filename: filename,
	}
}

func GetFile(path string) (*File, error) {
	file := NewFile(path)

	info, err := os.Stat(path)
	if err != nil {
		return file, err
	}
	sz := int(info.Size())
	file.Size = &sz
	file.Hashes, _ = GetFileHashes(path)
	return file, nil
}
