package atomic

import "fmt"

type TestBundle struct {
	Name            string `json:"display_name" yaml:"display_name"`
	AttackTechnique string `json:"attack_technique" yaml:"attack_technique"`
	AtomicTests     []Test `json:"atomic_tests" yaml:"atomic_tests"`
}

func (t *TestBundle) GetAttackTechniqueId() string {
	return t.AttackTechnique
}

func (t *TestBundle) GetAttackTechniqueName() string {
	return t.Name
}

func (t *TestBundle) DisplayName() string {
	return fmt.Sprintf("%s: %s", t.AttackTechnique, t.Name)
}
