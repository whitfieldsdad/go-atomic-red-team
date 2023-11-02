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
		layer.Techniques = append(layer.Techniques, struct {
			TechniqueID string `json:"techniqueID"`
			Enabled     bool   `json:"enabled"`
		}{
			TechniqueID: id,
			Enabled:     true,
		})
	}
	return &layer, nil
}

type AttackNavigatorLayer struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Techniques []struct {
		TechniqueID string `json:"techniqueID"`
		Enabled     bool   `json:"enabled"`
	} `json:"techniques"`
}

func (layer AttackNavigatorLayer) GetSelectedTechniqueIDs() []string {
	var ids []string
	for _, technique := range layer.Techniques {
		if technique.Enabled {
			ids = append(ids, technique.TechniqueID)
		}
	}
	return ids
}

func (layer AttackNavigatorLayer) ToTestPlan() TestPlan {
	testFilter := TestFilter{
		AttackTechniqueIds: layer.GetSelectedTechniqueIDs(),
	}
	return BulkTestPlan{
		Tests: []TestFilter{testFilter},
	}
}

func ParseAttackNavigatorLayer(data map[string]interface{}) (*AttackNavigatorLayer, error) {
	var layer AttackNavigatorLayer

	// Set `enabled` to true for all techniques that don't have it set.
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

// isAttackNavigatorLayer returns true if the map contains techniques.[].techniqueID.
func isAttackNavigatorLayer(data map[string]interface{}) bool {
	return true
}
