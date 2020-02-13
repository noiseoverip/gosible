package modules

import (
	"ansiblego/pkg/logging"
	"ansiblego/pkg/transport"
	"gopkg.in/yaml.v2"
)

type Debug struct {
	Var string
	//TODO: add msg with would be interpreted with jinja2
}

func LoadDebug(args map[string]interface{}) Module {
	return &Debug{Var: args["var"].(string)}
}

func(d *Debug) Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	value := "NULL"
	if v, ok := vars[d.Var]; ok {
		b, err := yaml.Marshal(v)
		if err != nil {
			logging.Info("ERROR: %v", err)
			return &ModuleExecResult{ Result: false }
		}
		value = string(b)
	}
	logging.Info(">> %s:\n%s\n>>\n", d.Var, value)
	return &ModuleExecResult{ Result: true }
}
