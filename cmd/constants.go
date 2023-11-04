package cmd

type OutputFormat string

const (
	JsonlOutputFormat OutputFormat = "jsonl"
	PlainOutputFormat OutputFormat = "plain"
	BriefOutputFormat OutputFormat = "brief"
)
