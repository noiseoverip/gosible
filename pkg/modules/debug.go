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

func (d *Debug) Run(_ Context, _ transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	value := "NULL"
	if v, ok := vars[d.Var]; ok {
		b, err := yaml.Marshal(v)
		if err != nil {
			logging.Debug("ERROR: %v", err)
			return &ModuleExecResult{Result: false}
		}
		value = string(b)
	}
	logging.Info("\t %s %s", d.Var, value)
	return &ModuleExecResult{Result: true}
}
