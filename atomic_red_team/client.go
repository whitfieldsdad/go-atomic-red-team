package atomic_red_team

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
	"gopkg.in/yaml.v3"
)

func ReadTests(directory, password string, filter *TestFilter) ([]Test, error) {
	file, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return readTestsFromDirectory(directory, filter)
	} else {
		if strings.HasSuffix(directory, ".yaml") || strings.HasSuffix(directory, ".yml") {
			return readTestsFromYamlFile(directory, filter)
		} else if strings.HasSuffix(directory, ".tar.gz") || strings.HasSuffix(directory, ".tar.gz.age") {
			return readTestsFromTarballFile(directory, password, filter)
		} else {
			return nil, errors.Errorf("unsupported file type: %s", directory)
		}
	}
	return nil, nil
}

func readTestsFromDirectory(directory string, filter *TestFilter) ([]Test, error) {
	var attackTechniqueIds []string
	if filter != nil {
		attackTechniqueIds = filter.AttackTechniqueIds
	}
	var paths []string
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		if len(attackTechniqueIds) > 0 && !containsAny(path, attackTechniqueIds) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})

	var tests []Test
	for _, path := range paths {
		testsFromFile, err := readTestsFromYamlFile(path, filter)
		if err != nil {
			log.Warnf("Failed to read tests from file: %s", err)
			continue
		}
		tests = append(tests, testsFromFile...)
	}
	return tests, nil
}

func readTestsFromYamlFile(path string, filter *TestFilter) ([]Test, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}
	return decodeAndFilterTests(data, filter)
}

func readTestsFromTarballFile(path string, password string, filter *TestFilter) ([]Test, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	var tarReader *tar.Reader
	if password != "" {
		tarReader, err = ReadEncryptedTarball(file, password)
	} else {
		tarReader, err = ReadTarball(file)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to read tarball")
	}
	fs := tarfs.New(tarReader)
	afs := &afero.Afero{Fs: fs}

	var tests []Test
	afs.Walk("/", func(path string, info os.FileInfo, err error) error {
		if StringMatchesPattern(path, "*T*/T*.yaml") {
			blob, err := afs.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "failed to read file")
			}
			testsFromFile, err := decodeAndFilterTests(blob, filter)
			if err != nil {
				log.Warnf("Failed to read tests from file: %s", err)
				return nil
			}
			tests = append(tests, testsFromFile...)
		}
		return nil
	})
	return tests, nil
}

func decodeAndFilterTests(data []byte, filter *TestFilter) ([]Test, error) {
	tests, err := decodeTests(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode tests")
	}
	tests = filterTests(tests, filter)
	return tests, nil
}

func decodeTests(data []byte) ([]Test, error) {
	var bundle TestBundle
	err := yaml.Unmarshal(data, &bundle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}
	attackTechniqueId := bundle.GetAttackTechniqueId()
	attackTechniqueName := bundle.GetAttackTechniqueName()

	tests := bundle.AtomicTests
	for i, test := range tests {
		test.AttackTechniqueId = attackTechniqueId
		test.AttackTechniqueName = attackTechniqueName

		for j, dependency := range test.Dependencies {
			dependency.executorName = test.DependencyExecutorName
			test.Dependencies[j] = dependency
		}
		tests[i] = test
	}
	return tests, nil
}

func filterTests(tests []Test, filter *TestFilter) []Test {
	var filteredTests []Test
	for _, test := range tests {
		if test.IsManual() {
			continue
		}
		if filter != nil && !filter.Matches(test) {
			continue
		}
		filteredTests = append(filteredTests, test)
	}
	return filteredTests
}
