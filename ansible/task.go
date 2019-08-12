package ansible

import (
	"ansiblego/modules"
	"fmt"
)

type Task struct {
	Name string `yaml:"name"`
	Module modules.Module
	When string		// unr-rendered 'when' attribute which controls if task will be executed or not
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var all map[string]interface{}
	err := unmarshal(&all)
	if err != nil {
		return err
	}
	for key, value := range all {
		switch key {
		case "when":
			t.When = fmt.Sprintf("%v", value) // This will make sure it is a string
		case "command":
			t.Module = modules.LoadCommand(map[string]string{ "stdin": value.(string) })
		case "name":
			t.Name = value.(string)
		}
		// TODO: default should probably module lookup form supported modules...
	}
	return nil
}
