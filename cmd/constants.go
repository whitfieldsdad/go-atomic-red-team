package cmd

const (
	JsonOutputFormat 			     = "json"
	JsonlOutputFormat                = "jsonl"
	PlainOutputFormat                = "plain"
	BriefOutputFormat                = "brief"
	AttackNavigatorLayerOutputFormat = "attack-navigator-layer"
)

const (
	DefaultOutputFormat = JsonlOutputFormat
)

var (
	OutputFormats = []string{JsonlOutputFormat, PlainOutputFormat, BriefOutputFormat}
)

var (
	listCommandAliases  = []string{"ls", "l"}
	countCommandAliases = []string{"total", "n"}
)
