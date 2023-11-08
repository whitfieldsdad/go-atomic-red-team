package atomic_red_team

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type Hashes struct {
	MD5    string `json:"md5,omitempty" yaml:"md5,omitempty"`
	SHA1   string `json:"sha1,omitempty" yaml:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
}

func (h *Hashes) Empty() bool {
	return h.MD5 == "" && h.SHA1 == "" && h.SHA256 == ""
}

// GetFileHashes returns the MD5, SHA1, and SHA256 hashes of a file.
func GetFileHashes(path string) (*Hashes, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return GetReaderHashes(f)
}

// GetReaderHashes returns the MD5, SHA1, and SHA256 hashes of a reader.
func GetReaderHashes(rd io.Reader) (*Hashes, error) {
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()

	pagesize := os.Getpagesize()
	reader := bufio.NewReaderSize(rd, pagesize)
	multiWriter := io.MultiWriter(md5, sha1, sha256)
	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return nil, err
	}
	hashes := &Hashes{
		MD5:    fmt.Sprintf("%x", md5.Sum(nil)),
		SHA1:   fmt.Sprintf("%x", sha1.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256.Sum(nil)),
	}
	return hashes, nil
}
