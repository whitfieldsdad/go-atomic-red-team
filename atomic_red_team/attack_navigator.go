package atomic_red_team

import (
	"encoding/json"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

const (
	enterpriseAttack = "enterprise-attack"
)

func NewAttackNavigatorLayer(techniqueIds []string) (*AttackNavigatorLayer, error) {
	layer := AttackNavigatorLayer{
		Name:   "Atomic Red Team",
		Domain: enterpriseAttack,
	}
	for _, id := range techniqueIds {
		layer.Techniques = append(layer.Techniques, AttackNavigatorTechnique{
			TechniqueID: id,
			Enabled:     true,
		})
	}
	return &layer, nil
}

type AttackNavigatorLayer struct {
	Name       string                     `json:"name" yaml:"name"`
	Domain     string                     `json:"domain" yaml:"domain"`
	Techniques []AttackNavigatorTechnique `json:"techniques" yaml:"techniques"`
}

type AttackNavigatorTechnique struct {
	TechniqueID string `json:"techniqueID" yaml:"techniqueID"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Color       string `json:"color,omitempty" yaml:"color,omitempty"`
}

func (layer AttackNavigatorLayer) GetSelectedTechniqueIDs() []string {
	var ids []string
	for _, technique := range layer.Techniques {
		if technique.TechniqueID == "" || technique.Color == "" {
			continue
		}
		if !technique.Enabled {
			continue
		}
		ids = append(ids, technique.TechniqueID)
	}
	slices.Sort(ids)
	return ids
}

func (layer AttackNavigatorLayer) ToTestPlan() TestPlanInterface {
	testFilter := TestFilter{
		AttackTechniqueIds: layer.GetSelectedTechniqueIDs(),
	}
	return BulkTestPlan{
		Tests: []TestFilter{testFilter},
	}
}

func ReadAttackNavigatorLayer(path string) (*AttackNavigatorLayer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	blob, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(blob, &m)
	if err != nil {
		return nil, err
	}
	return ParseAttackNavigatorLayer(m)
}

func ParseAttackNavigatorLayer(data map[string]interface{}) (*AttackNavigatorLayer, error) {
	var layer AttackNavigatorLayer
	techniques, _ := data["techniques"].([]interface{})
	for i, technique := range techniques {
		m, _ := technique.(map[string]interface{})
		if _, ok := m["enabled"]; !ok {
			m["enabled"] = true
			techniques[i] = m
		}
	}
	err := mapstructure.Decode(data, &layer)
	if err != nil {
		return nil, err
	}
	return &layer, nil
}
