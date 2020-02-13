package ansible

import (
	"ansiblego/pkg/modules"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Task struct {
	Name     string `yaml:"name"`
	Module   modules.Module
	When     string		// raw 'when' attribute which controls if task will be executed or not
	Register string	// variable name to register tasks result to
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var all map[string]interface{}
	err := unmarshal(&all)
	if err != nil {
		return err
	}
	for key, value := range all {
		switch key {
		case "name":
			t.Name = value.(string)
		case "when":
			t.When = fmt.Sprintf("%v", value) // This will make sure it is a string
		case "register":
			t.Register = fmt.Sprintf("%v", value)
		case "debug":
			params := map[string]interface{}{}
			err := mapstructure.Decode(value, &params)
			if err != nil {
				return err
			}
			t.Module = modules.LoadDebug(params)
		case "command":
			t.Module = modules.LoadCommand(map[string]string{ "stdin": value.(string) })
		case "template":
			params := map[string]string{}
			err := mapstructure.Decode(value, &params)	// TODO: this might be slow, need to investigate
			if err != nil {
				return err
			}
			t.Module = modules.NewTemplate(params)
		case "copy":
			params := map[string]string{}
			err := mapstructure.Decode(value, &params)
			if err != nil {
				return err
			}
			t.Module = modules.LoadCopy(params)
		case "assert":
			params := map[string]string{}
			err := mapstructure.Decode(value, &params)
			if err != nil {
				return err
			}
			t.Module = modules.LoadAssert(params)
		case "set_fact":
			params := map[string]interface{}{}
			err := mapstructure.Decode(value, &params)
			if err != nil {
				return err
			}
			t.Module = modules.LoadSetHostFact(params)
		}
		// TODO: default should probably module lookup form supported modules...
	}
	return nil
}
