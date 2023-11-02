package atomic_red_team

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ClientOptions struct {
	AtomicsDir string `json:"atomic_red_team_dir"`
}

type Client struct {
	Options ClientOptions `json:"options"`
}

func NewClient(atomicsDir string) (*Client, error) {
	if atomicsDir == "" {
		return nil, errors.New("path to atomic-red-team/atomics directory is required")
	}
	return &Client{
		Options: ClientOptions{
			AtomicsDir: atomicsDir,
		},
	}, nil
}

func getAttackTechniqueIdsFromTestFilters(testFilters []TestFilter) []string {
	var attackTechniqueIds []string
	for _, testFilter := range testFilters {
		for _, attackTechniqueId := range testFilter.AttackTechniqueIds {
			if !slices.Contains(attackTechniqueIds, attackTechniqueId) {
				attackTechniqueIds = append(attackTechniqueIds, attackTechniqueId)
			}
		}
	}
	return attackTechniqueIds
}

func (c Client) ListTests(testFilters []TestFilter) ([]Test, error) {
	attackTechniqueIds := getAttackTechniqueIdsFromTestFilters(testFilters)

	var tests []Test
	testPaths, err := c.ListTestPaths(attackTechniqueIds)
	if err != nil {
		return nil, err
	}
	for _, testPath := range testPaths {
		testsFromFile, err := c.readTestsFromFile(testPath, testFilters)
		if err != nil {
			log.Warnf("Failed to read tests from file %s - %s", testPath, err)
			continue
		}
		log.Debugf("Read %d tests from %s", len(testsFromFile), testPath)
		tests = append(tests, testsFromFile...)
	}
	return tests, nil
}

func (c Client) ListTestPaths(attackTechniqueIds []string) ([]string, error) {
	var testPaths []string
	filepath.Walk(c.Options.AtomicsDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !filenameContainsAttackTechniqueId(info.Name(), attackTechniqueIds) {
			return nil
		}
		testPaths = append(testPaths, path)
		return nil
	})
	return testPaths, nil
}

func filenameContainsAttackTechniqueId(filename string, attackTechniqueIds []string) bool {
	if len(attackTechniqueIds) == 0 {
		g := glob.MustCompile("*T*.yaml")
		return g.Match(filename)
	}
	for _, attackTechniqueId := range attackTechniqueIds {
		g := glob.MustCompile(fmt.Sprintf("*%s*.yaml", attackTechniqueId))
		if g.Match(filename) {
			return true
		}
	}
	return false
}

func (c Client) readBundleFromFile(path string, testFilters []TestFilter) (*TestBundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bundle TestBundle
	err = yaml.Unmarshal(data, &bundle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}
	// Set the attack technique ID and name for each test.
	for i := 0; i < len(bundle.AtomicTests); i++ {
		test := &bundle.AtomicTests[i]
		test.AttackTechniqueId = bundle.GetAttackTechniqueId()
		test.AttackTechniqueName = bundle.GetAttackTechniqueName()
	}

	var tests []Test
	for _, test := range bundle.AtomicTests {
		if test.Executor.Name != "manual" && test.MatchesAnyFilter(testFilters) {

			// Set the executor for each dependency.
			for i := 0; i < len(test.Dependencies); i++ {
				dependency := &test.Dependencies[i]
				dependency.executorType = test.Executor.Name
			}
			tests = append(tests, test)
		}
	}
	bundle.AtomicTests = tests
	return &bundle, nil
}

func (c Client) readTestsFromFile(path string, testFilters []TestFilter) ([]Test, error) {
	bundle, err := c.readBundleFromFile(path, testFilters)
	if err != nil {
		return nil, err
	}
	return bundle.AtomicTests, nil
}
