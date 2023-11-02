package atomic_red_team

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// FindAtomicsDir is used to locate a copy of the redcanaryco/atomic-red-team git repository on the local system.
func FindAtomicsDir() (string, error) {
	if DefaultAtomicsDir != "" {
		return DefaultAtomicsDir, nil
	}
	log.Infof("Searching for atomics directory...")
	startTime := time.Now()
	dir, err := findAtomicsDir()
	if err != nil {
		return "", err
	}
	endTime := time.Now()
	log.Infof("Found atomics directory: %s (took %.2f seconds)", dir, endTime.Sub(startTime).Seconds())
	return dir, nil
}

func findAtomicsDir() (string, error) {
	for _, root := range getAutoDiscoveryPaths() {
		log.Infof("Searching for atomics directory in: %s", root)
		if root == "" {
			continue
		}
		file, err := os.Open(root)
		if err != nil {
			continue
		}
		info, err := file.Stat()
		if err != nil {
			continue
		}
		if info.IsDir() {
			dir := ""
			filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					return nil
				}
				if info.Name() == "atomics" {
					dir = path
					return filepath.SkipDir
				}
				return nil
			})
			if dir != "" {
				return dir, nil
			}
		}
	}
	return "", errors.New("atomics directory not found")
}

func getAutoDiscoveryPaths() []string {
	paths := []string{}

	// Add the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Warn("Failed to identify current working directory: %s\n", err)
	} else {
		paths = append(paths, cwd)
	}

	// Add the user's home directory (i.e. ~, $HOME, %USERPROFILE%, etc.)
	homedir, err := homedir.Dir()
	if err != nil {
		log.Warn("Failed to locate home directory: %s\n", err)
	} else {
		paths = append(paths, homedir)
	}

	// Add platform-specific directories
	if runtime.GOOS == "windows" {
		drive := os.ExpandEnv("$SYSTEMDRIVE")
		if drive != "" {
			subdirs := []string{"Program Files", "Program Files (x86)", "ProgramData", "Users", "Windows"}
			for _, subdir := range subdirs {
				paths = append(paths, filepath.Join(drive, subdir))
			}
		}
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		paths = append(paths, "/opt", "/usr/", "/var", "/etc", "/tmp", "/private")
	}
	return paths
}
