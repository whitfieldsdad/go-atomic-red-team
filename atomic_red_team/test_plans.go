package atomic_red_team

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type TestPlan interface {
	GetTestFilters() []TestFilter
	GetTestOptions() TestOptions
}

type BulkTestPlan struct {
	AtomicsDir     string                 `json:"atomics_dir,omitempty"`
	Tests          []TestFilter           `json:"tests,omitempty"`
	InputArguments map[string]interface{} `json:"input_arguments,omitempty"`
}

func (p BulkTestPlan) GetTestFilters() []TestFilter {
	var filters []TestFilter
	filters = append(filters, p.Tests...)
	return removeEmptyTestFilters(filters)
}

func (p BulkTestPlan) GetTestOptions() TestOptions {
	return TestOptions{
		AtomicsDir:     p.AtomicsDir,
		InputArguments: p.InputArguments,
	}
}

// MultiTestPlan provides a way to run multiple tests with different input arguments and an optional set of global arguments.
type MultiTestPlan struct {
	AtomicsDir     string                 `json:"atomics_dir,omitempty"`
	Tests          []testReference        `json:"tests,omitempty"`
	InputArguments map[string]interface{} `json:"input_arguments,omitempty"`
}

func (p MultiTestPlan) GetTestFilters() []TestFilter {
	var filters []TestFilter
	for _, test := range p.Tests {
		filters = append(filters, test.GetTestFilter())
	}
	return removeEmptyTestFilters(filters)
}

func (p MultiTestPlan) GetTestOptions() TestOptions {
	return TestOptions{
		AtomicsDir:     p.AtomicsDir,
		InputArguments: p.InputArguments,
	}
}

type testReference struct {
	Id                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	Platforms         []string `json:"platforms,omitempty"`
	ElevationRequired *bool    `json:"elevation_required,omitempty"`
	AttackTechniqueId string   `json:"attack_technique_id,omitempty"`
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
func ReadTestPlan(path string) (TestPlan, error) {
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

func ParseTestPlan(data map[string]interface{}) (TestPlan, error) {
	if isAttackNavigatorLayer(data) {
		layer, err := ParseAttackNavigatorLayer(data)
		if err != nil {
			return nil, err
		}
		return layer.ToTestPlan(), nil
	}
	return parseTestPlan(data)
}

func parseTestPlan(data map[string]interface{}) (TestPlan, error) {
	var plan TestPlan
	var err error

	// Parse the test plan as either a multi-test plan or as a bulk test plan.
	testPlanFields := getMapKeySet(data)
	testReferenceFields, _, testFilterFields := diffStructFieldSets(testReference{}, TestFilter{})

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

func getMapKeySet(m map[string]interface{}) mapset.Set[string] {
	keys := mapset.NewSet[string]()
	for k := range m {
		keys.Add(k)
	}
	return keys
}

// diffStructFieldSets returns the triple: (a - b), (a ∩ b), (b - a)
func diffStructFieldSets(a, b interface{}) (mapset.Set[string], mapset.Set[string], mapset.Set[string]) {
	sa := getStructFieldSet(a)
	sb := getStructFieldSet(b)
	si := sa.Intersect(sb)
	sa = sa.Difference(si)
	sb = sb.Difference(si)
	return sa, si, sb
}

func getStructFields(i interface{}) []string {
	var fields []string
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

func getStructFieldSet(i interface{}) mapset.Set[string] {
	fields := mapset.NewSet[string]()
	for _, field := range getStructFields(i) {
		fields.Add(field)
	}
	return fields
}
