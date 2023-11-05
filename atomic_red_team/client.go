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

func GetTests(path string, filter *TestFilter) ([]Test, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat path")
	}
	if info.IsDir() {
		return getTestsFromDirectory(path, filter)
	} else if strings.HasSuffix(path, ".tar") {
		return getTestsFromTarFile(path, filter)
	} else if strings.HasSuffix(path, ".tar.gz") {
		return getTestsFromTarballFile(path, filter)
	} else {
		log.Fatalf("Unsupported file type: %s", path)
	}
	return nil, nil
}

func getTestsFromDirectory(directory string, filter *TestFilter) ([]Test, error) {
	paths, err := findTests(directory, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find tests")
	}
	var tests []Test
	for _, path := range paths {
		testsFromFile, err := ReadTestsFromFile(path, filter)
		if err != nil {
			log.Warnf("Failed to read tests from file: %s", err)
			continue
		}
		tests = append(tests, testsFromFile...)
	}
	return tests, nil
}

func getTestsFromTarFile(path string, filter *TestFilter) ([]Test, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	tfs := tarfs.New(tar.NewReader(file))
	afs := &afero.Afero{Fs: tfs}

	var tests []Test

	afs.Walk("/", func(path string, info os.FileInfo, err error) error {
		if StringMatchesPattern(path, "*T*/T*.yaml") {
			blob, err := afs.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "failed to read file")
			}
			testsFromFile, err := decodeTests(blob)
			if err != nil {
				log.Warnf("Failed to decode tests from file: %s", err)
				return nil
			}
			tests = append(tests, filterTests(testsFromFile, filter)...)
		}
		return nil
	})
	return tests, nil
}

func getTestsFromTarballFile(path string, filter *TestFilter) ([]Test, error) {
	return nil, nil
}

func findTests(directory string, attackTechniqueIds []string) ([]string, error) {
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
	return paths, nil
}

func ReadTestsFromFile(path string, filter *TestFilter) ([]Test, error) {
	if strings.HasSuffix(path, ".yaml") {
		return readTestsFromYamlFile(path, filter)
	}
	return readTestsFromYamlFile(path, filter)
}

func readTestsFromYamlFile(path string, filter *TestFilter) ([]Test, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}
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
