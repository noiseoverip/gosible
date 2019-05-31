package ansible

import (
	"ansiblego/modules"
)

type Task struct {
	Name string `yaml:"name"`
	Module modules.Module
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var all map[string]interface{}
	err := unmarshal(&all)
	if err != nil {
		return err
	}
	for key, value := range all {
		// TODO: key should be either top level attribute like "name", "with_items" or module name
		if key == "command" {
			t.Module = modules.LoadCommand(map[string]string{ "stdin": value.(string) })
		}
		if key == "name" {
			t.Name = value.(string)
		}
	}
	return nil
}
