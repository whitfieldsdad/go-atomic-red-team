package atomic_red_team

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

type FileFilter struct {
	FilenamePatterns []string `json:"filenames" yaml:"filenames"`
	PathPatterns     []string `json:"paths" yaml:"paths"`
}

func (f FileFilter) Matches(path string, info os.FileInfo) bool {
	if len(f.FilenamePatterns) > 0 {
		if !StringMatchesAnyPattern(info.Name(), f.FilenamePatterns) {
			return false
		}
	}
	if len(f.PathPatterns) > 0 {
		if !StringMatchesAnyPattern(path, f.PathPatterns) {
			return false
		}
	}
	return true
}

type FileStat struct {
	Path string
	Info os.FileInfo
}

func GetFS() afero.Afero {
	fs := afero.NewOsFs()
	return afero.Afero{Fs: fs}
}

func GetTarFS(path, password string) (*afero.Afero, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	var reader io.Reader
	reader = file

	// If a password was provided, decrypt the tarball.
	if password != "" {
		id, err := age.NewScryptIdentity(password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create identity")
		}
		reader, err = age.Decrypt(file, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt tarball")
		}
	}

	// If the tarball is GZIP compressed, decompress it.
	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tar.gz.age") {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create gzip reader")
		}
	}

	// Create a tar filesystem from the tarball.
	fs := tarfs.New(tar.NewReader(reader))
	afs := &afero.Afero{Fs: fs}
	return afs, nil
}

func FindFiles(afs afero.Afero, filter *FileFilter) ([]File, error) {
	var results []File
	afs.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filter != nil && !filter.Matches(path, info) {
			return nil
		}
		result := NewFile(path)
		results = append(results, *result)
		return nil
	})
	return results, nil
}

func FindFilesInTarball(path, password string, filter *FileFilter) ([]File, error) {
	afs, err := GetTarFS(path, password)
	if err != nil {
		return nil, err
	}
	return FindFiles(*afs, filter)
}
