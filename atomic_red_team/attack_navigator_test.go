package atomic_red_team

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestLayerPath = "../data/attack-navigator/groups/FIN7.json"
)

func TestGetTestLayer(t *testing.T) {
	_, err := os.Stat(TestLayerPath)
	assert.Nil(t, err)
}

func TestGetSelectedTechniques(t *testing.T) {
	layer, err := ReadAttackNavigatorLayer(TestLayerPath)
	assert.Nil(t, err)

	techniques := layer.Techniques
	assert.Equal(t, 59, len(techniques))

	selectedTechniqueIds := layer.GetSelectedTechniqueIDs()
	assert.Equal(t, 39, len(selectedTechniqueIds))

	expectedTechniqueIds := []string{"T1005", "T1008", "T1021.001", "T1021.004", "T1021.005", "T1027.010", "T1036.004", "T1036.005", "T1047", "T1053.005", "T1059.001", "T1059.003", "T1059.005", "T1059.007", "T1071.004", "T1078", "T1091", "T1102.002", "T1105", "T1113", "T1125", "T1204.001", "T1204.002", "T1210", "T1218.005", "T1486", "T1497.002", "T1543.003", "T1546.011", "T1547.001", "T1553.002", "T1558.003", "T1559.002", "T1566.001", "T1566.002", "T1567.002", "T1571", "T1583.001", "T1587.001"}
	assert.Equal(t, expectedTechniqueIds, selectedTechniqueIds)
}

func TestToTestPlan(t *testing.T) {
	layer, err := ReadAttackNavigatorLayer(TestLayerPath)
	assert.Nil(t, err)

	plan := layer.ToTestPlan()
	assert.NotNil(t, plan)

	filters := plan.GetTestFilters()
	assert.Equal(t, 1, len(filters))

	filter := filters[0]
	assert.Equal(t, 39, len(filter.AttackTechniqueIds))

	expectedTechniqueIds := []string{"T1005", "T1008", "T1021.001", "T1021.004", "T1021.005", "T1027.010", "T1036.004", "T1036.005", "T1047", "T1053.005", "T1059.001", "T1059.003", "T1059.005", "T1059.007", "T1071.004", "T1078", "T1091", "T1102.002", "T1105", "T1113", "T1125", "T1204.001", "T1204.002", "T1210", "T1218.005", "T1486", "T1497.002", "T1543.003", "T1546.011", "T1547.001", "T1553.002", "T1558.003", "T1559.002", "T1566.001", "T1566.002", "T1567.002", "T1571", "T1583.001", "T1587.001"}
	assert.Equal(t, expectedTechniqueIds, filter.AttackTechniqueIds)
}
