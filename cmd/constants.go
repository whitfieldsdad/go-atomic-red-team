package cmd

import "strings"

type OutputFormat string

const (
	OutputFormatJson  = "json"
	OutputFormatYaml  = "yaml"
	OutputFormatPlain = "plain"
	OutputFormatBrief = "brief"
)

var (
	lineSeparator = strings.Repeat("-", 80)
)
