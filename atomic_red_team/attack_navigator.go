package atomic_red_team

import (
	"github.com/mitchellh/mapstructure"
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
	Name       string                     `json:"name"`
	Domain     string                     `json:"domain"`
	Techniques []AttackNavigatorTechnique `json:"techniques"`
}

type AttackNavigatorTechnique struct {
	TechniqueID string `json:"techniqueID"`
	Enabled     bool   `json:"enabled"`
	Color       string `json:"color,omitempty"`
}

func (layer AttackNavigatorLayer) GetSelectedTechniqueIDs() []string {
	var ids []string
	for _, technique := range layer.Techniques {
		if !technique.Enabled {
			continue
		}
		if technique.TechniqueID == "" || technique.Color == "" {
			continue
		}
	}
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
