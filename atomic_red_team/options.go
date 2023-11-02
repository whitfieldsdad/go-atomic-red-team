package atomic_red_team

import (
	"os"

	"github.com/pkg/errors"
)

const (
	IncludeParentProcesses = true
)

type TestOptions struct {
	AtomicsDir     string                 `json:"atomics_dir"`
	InputArguments map[string]interface{} `json:"input_arguments"`
}

func NewTestOptions(atomicsDir string) (*TestOptions, error) {
	if atomicsDir == "" {
		atomicsDir = DefaultAtomicsDir
	}
	if _, err := os.Stat(atomicsDir); os.IsNotExist(err) {
		return nil, errors.Errorf("atomics directory does not exist: %s", atomicsDir)
	}
	return &TestOptions{
		AtomicsDir:     atomicsDir,
		InputArguments: make(map[string]interface{}),
	}, nil
}

type ProcessOptions struct {
	IncludeParentProcesses bool `json:"include_parent_processes"`
}

func NewProcessOptions() *ProcessOptions {
	return &ProcessOptions{
		IncludeParentProcesses: IncludeParentProcesses,
	}
}

type FileOptions struct {
	HashingOptions *HashingOptions `json:"hashing_options"`
}

func NewFileOptions() *FileOptions {
	return &FileOptions{}
}

const (
	IncludeMD5        = true
	IncludeSHA1       = true
	IncludeSHA256     = true
	IncludeSHA512     = true
	IncludeXXH64      = true
	IncludeBLAKE3_256 = true
)

type HashingOptions struct {
	MD5        bool `json:"md5"`
	SHA1       bool `json:"sha1"`
	SHA256     bool `json:"sha256"`
	SHA512     bool `json:"sha512"`
	XXH64      bool `json:"xxh64"`
	BLAKE3_256 bool `json:"blake3_256"`
}

func NewHashingOptions() *HashingOptions {
	return &HashingOptions{
		MD5:        IncludeMD5,
		SHA1:       IncludeSHA1,
		SHA256:     IncludeSHA256,
		SHA512:     IncludeSHA512,
		XXH64:      IncludeXXH64,
		BLAKE3_256: IncludeBLAKE3_256,
	}
}
