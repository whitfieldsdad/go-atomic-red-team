package atomic_red_team

import (
	"reflect"

	mapset "github.com/deckarep/golang-set/v2"
)

// diffStructFields returns the triple: (a - b), (a ∩ b), (b - a)
func diffStructFields(a, b interface{}) (mapset.Set[string], mapset.Set[string], mapset.Set[string]) {
	sa := getStructFields(a)
	sb := getStructFields(b)
	si := sa.Intersect(sb)
	sa = sa.Difference(si)
	sb = sb.Difference(si)
	return sa, si, sb
}

// getStructFields returns a list of all struct field names.
func getStructFields(i interface{}) mapset.Set[string] {
	fields := mapset.NewSet[string]()
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i).Name
		fields.Add(field)
	}
	return fields
}

// getMapKeys returns the keys in a map.
func getMapKeys(m map[string]interface{}) mapset.Set[string] {
	keys := mapset.NewSet[string]()
	for k := range m {
		keys.Add(k)
	}
	return keys
}
