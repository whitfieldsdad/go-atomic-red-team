package atomic_red_team

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"

	"github.com/cespare/xxhash/v2"
	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/mapstructure"
	"github.com/zeebo/blake3"
)

const (
	FileReadBufferSize int32 = 1000000
)

type Hashes struct {
	MD5        string `json:"md5"`
	SHA1       string `json:"sha1"`
	SHA256     string `json:"sha256"`
	SHA512     string `json:"sha512"`
	XXH64      string `json:"xxh64"`
	BLAKE3_256 string `json:"blake3_256"`
}

func (h Hashes) List() []string {
	hashes := []string{
		h.MD5,
		h.SHA1,
		h.SHA256,
		h.SHA512,
		h.XXH64,
		h.BLAKE3_256,
	}
	var setHashes []string
	for _, hash := range hashes {
		if hash != "" {
			setHashes = append(setHashes, hash)
		}
	}
	return setHashes
}

func GetFileHashes(path string, opts *HashingOptions) (*Hashes, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	sz := info.Size()
	log.Debugf("Hashing %s (size: %s)", path, humanize.Bytes(uint64(sz)))

	defer f.Close()
	return GetReaderHashes(f, opts)
}

func GetReaderHashes(rd io.Reader, opts *HashingOptions) (*Hashes, error) {
	if opts == nil {
		opts = NewHashingOptions()
	}

	writers := map[string]io.Writer{}
	if opts.MD5 {
		writers["md5"] = md5.New()
	}
	if opts.SHA1 {
		writers["sha1"] = sha1.New()
	}
	if opts.SHA256 {
		writers["sha256"] = sha256.New()
	}
	if opts.SHA512 {
		writers["sha512"] = sha512.New()
	}
	if opts.XXH64 {
		writers["xxh64"] = xxhash.New()
	}
	if opts.BLAKE3_256 {
		writers["blake3_256"] = blake3.New()
	}

	reader := bufio.NewReaderSize(rd, os.Getpagesize())

	writerList := []io.Writer{}
	for _, writer := range writers {
		writerList = append(writerList, writer)
	}
	multiWriter := io.MultiWriter(writerList...)
	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return nil, err
	}

	hashes := map[string]string{}
	for name, writer := range writers {
		hashes[name] = fmt.Sprintf("%x", writer.(hash.Hash).Sum(nil))
	}
	result := &Hashes{}
	err = mapstructure.Decode(hashes, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetMD5(b []byte) string {
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}

func GetSHA1(b []byte) string {
	h := sha1.Sum(b)
	return fmt.Sprintf("%x", h)
}

func GetSHA256(b []byte) string {
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h)
}

// GetSha256 returns the SHA256 hash of a byte slice.
func GetSHA512(b []byte) string {
	h := sha512.Sum512(b)
	return fmt.Sprintf("%x", h)
}

func GetXXH64(b []byte) string {
	h := xxhash.Sum64(b)
	return strconv.FormatUint(h, 16)
}

func GetFileMD5(path string) (string, error) {
	return getFileHash(path, md5.New())
}

func GetFileSHA1(path string) (string, error) {
	return getFileHash(path, sha1.New())
}

func GetFileSHA256(path string) (string, error) {
	return getFileHash(path, sha256.New())
}

func GetFileSHA512(path string) (string, error) {
	return getFileHash(path, sha512.New())
}

func GetFileXXH64(path string) (string, error) {
	return getFileHash(path, xxhash.New())
}

func getFileHash(path string, h hash.Hash) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
