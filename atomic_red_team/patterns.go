package atomic_red_team

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
)

var (
	CaseSensitive = false
)

// IsRegExp checks if the provided pattern is a regular expression (e.g. ^T\d{4}$)
func IsRegExp(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

// IsGlobPattern checks if the provided pattern is a glob pattern (e.g. *.yaml)
func IsGlobPattern(pattern string) bool {
	_, err := glob.Compile(pattern)
	return err == nil
}

// StringMatchesGlobPattern checks if a string matches a glob pattern (e.g. *.yaml).
func StringMatchesGlobPattern(value, pattern string) bool {
	if !CaseSensitive {
		value = strings.ToLower(value)
		pattern = strings.ToLower(pattern)
	}
	g, err := glob.Compile(pattern)
	if err != nil {
		log.Warnf("Could not compile glob pattern: %s", pattern)
		return false
	}
	return g.Match(value)
}

// StringMatchesRegExp checks if a string matches a regular expression (e.g. ^T\d{4}$).
func StringMatchesRegExp(value, pattern string) bool {
	if !CaseSensitive {
		value = strings.ToLower(value)
		pattern = strings.ToLower(pattern)
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Warnf("Could not compile regular expression: %s", pattern)
		return false
	}
	return r.MatchString(value)
}

// StringMatchesPattern checks if a string matches a regular expression or glob pattern (e.g. *.yaml).
func StringMatchesPattern(value, pattern string) bool {
	if pattern == "" {
		return true
	}
	if IsRegExp(pattern) {
		return StringMatchesRegExp(value, pattern)
	}
	if IsGlobPattern(pattern) {
		return StringMatchesGlobPattern(value, pattern)
	}
	log.Warnf("Unsupported pattern - pattern is not a regular expression or glob pattern: %s", pattern)
	return false
}

// StringMatchesAnyPattern checks if a string matches any of the provided glob patterns (e.g. *.yaml).
func StringMatchesAnyPattern(value string, patterns []string) bool {
	for _, pattern := range patterns {
		if StringMatchesPattern(value, pattern) {
			return true
		}
	}
	return false
}

// AnyStringMatchesPattern checks if any of the provided strings matches a glob pattern (e.g. *.yaml).
func AnyStringMatchesPattern(values []string, pattern string) bool {
	for _, value := range values {
		if StringMatchesPattern(value, pattern) {
			return true
		}
	}
	return false
}

// AnyStringMatchesAnyPattern checks if any of the provided strings matches any of the provided glob patterns (e.g. *.yaml).
func AnyStringMatchesAnyPattern(values, patterns []string) bool {
	for _, value := range values {
		if StringMatchesAnyPattern(value, patterns) {
			return true
		}
	}
	return false
}
