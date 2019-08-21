package modules

import (
	"ansiblego/transport"
	"fmt"
)

type SetHostFact struct {
	FactName string			// name
	FactValue interface{}	// value
}

func LoadSetHostFact(args map[string]interface{}) Module {
	if len(args) > 2 {
		panic(fmt.Errorf("Too many attributes"))
	}
	// TODO: find how to best extract fidrst map key not knowing it's value and do it in one line
	m := &SetHostFact{}
	for k,v := range args {
		m.FactName = k
		m.FactValue = v
	}
	return &SetHostFact{FactName: args["name"].(string), FactValue: args["value"]}
}

func(c *SetHostFact) Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	vars[c.FactName] = c.FactValue
	return &ModuleExecResult{ Result: true }
}
