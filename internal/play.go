package internal

type Play struct {
	HostSelector string   `yaml:"hosts"`
	Tasks        []Task   `yaml:"tasks,omitempty"`
	Roles        []string `yaml:"roles,omitempty"`
}
