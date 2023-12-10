package atomic

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type TestPlanInterface interface {
	GetTestFilters() []TestFilter
	GetTestOptions() TestOptions
}

type BulkTestPlan struct {
	Tests          []TestFilter           `json:"tests" yaml:"tests"`
	InputArguments map[string]interface{} `json:"input_arguments" yaml:"input_arguments"`
}

func (plan BulkTestPlan) GetTestFilters() []TestFilter {
	var filters []TestFilter
	filters = append(filters, plan.Tests...)
	return filters
}

func (plan BulkTestPlan) GetTestOptions() TestOptions {
	return TestOptions{
		InputArguments: plan.InputArguments,
	}
}

type TestPlan struct {
	Tests          []testReference        `json:"tests" yaml:"tests"`
	InputArguments map[string]interface{} `json:"input_arguments" yaml:"input_arguments"`
}

func (plan TestPlan) GetTestFilters() []TestFilter {
	var filters []TestFilter
	for _, test := range plan.Tests {
		filters = append(filters, test.GetTestFilter())
	}
	return filters
}

func (plan TestPlan) GetTestOptions() TestOptions {
	return TestOptions{
		InputArguments: plan.InputArguments,
	}
}

type testReference struct {
	Id                string   `json:"id" yaml:"id"`
	Name              string   `json:"name" yaml:"name"`
	Description       string   `json:"description" yaml:"description"`
	Platforms         []string `json:"platforms" yaml:"platforms"`
	ElevationRequired *bool    `json:"elevation_required" yaml:"elevation_required"`
	AttackTechniqueId string   `json:"attack_technique_id" yaml:"attack_technique_id"`
}

func (t testReference) GetTestFilter() TestFilter {
	f := TestFilter{}
	if t.Id != "" {
		f.Ids = []string{t.Id}
	}
	if t.Name != "" {
		f.Names = []string{t.Name}
	}
	if t.Description != "" {
		f.Descriptions = []string{t.Description}
	}
	if t.AttackTechniqueId != "" {
		f.AttackTechniqueIds = []string{t.AttackTechniqueId}
	}
	if len(t.Platforms) > 0 {
		f.Platforms = t.Platforms
	}
	f.ElevationRequired = t.ElevationRequired
	return f
}

// ReadTestPlan reads a test plan from a file.
func ReadTestPlan(path string) (TestPlanInterface, error) {
	log.Infof("Reading test plan: %s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, errors.Wrap(err, "JSON deserialization failed")
	}
	plan, err := ParseTestPlan(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse test plan")
	}
	return plan, nil
}

func ParseTestPlan(data map[string]interface{}) (TestPlanInterface, error) {
	if isAttackNavigatorLayer(data) {
		panic("not implemented")
	} else {
		return parseTestPlan(data)
	}
}

func isAttackNavigatorLayer(data map[string]interface{}) bool {
	_, ok := data["techniques"]
	if ok {
		for _, technique := range data["techniques"].([]interface{}) {
			m, _ := technique.(map[string]interface{})
			if _, ok := m["techniqueID"]; ok {
				return true
			}
		}
	}
	return false
}

func parseTestPlan(data map[string]interface{}) (TestPlanInterface, error) {
	var plan TestPlanInterface
	var err error

	// Parse the test plan as either a multi-test plan or as a bulk test plan.
	testPlanFields := getMapKeys(data)
	testReferenceFields, _, testFilterFields := diffStructFields(testReference{}, TestFilter{})

	if testPlanFields.IsSubset(testReferenceFields) {
		log.Info("Parsing test plan as a multi-test plan")
		err = mapstructure.Decode(data, &plan)
	} else if testPlanFields.IsSubset(testFilterFields) {
		log.Info("Parsing test plan as a bulk test plan")
		err = mapstructure.Decode(data, &plan)
	} else {
		return nil, errors.New("failed to determine test plan type")
	}
	return plan, err
}
