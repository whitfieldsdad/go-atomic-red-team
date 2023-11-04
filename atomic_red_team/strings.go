package atomic_red_team

import "strings"

func containsAny(value string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(value, substring) {
			return true
		}
	}
	return false
}
