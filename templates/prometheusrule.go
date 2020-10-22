package template
type PrometheusRule struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
type Labels struct {
	Prometheus string `yaml:"prometheus"`
	Role       string `yaml:"role"`
}
type Metadata struct {
	Labels            Labels      `yaml:"labels"`
	Name              string      `yaml:"name"`
}
type Spec struct {
	Groups []Groups `yaml:"groups"`
}