package template

type RecordingRule struct {
	Groups []Groups `yaml:"groups"`
}
type Rules struct {
	Expr   string `yaml:"expr"`
	Record string `yaml:"record"`
}
type Groups struct {
	Name  string  `yaml:"name"`
	Rules []Rules `yaml:"rules"`
}