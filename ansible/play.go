package ansible

type Play struct {
	Hosts string `yaml:"hosts"`
	Tasks []Task `yaml:"tasks",omitempty`
	Roles []string `yaml:"roles",omitempty`
}