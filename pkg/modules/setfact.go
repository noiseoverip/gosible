package modules

import (
	"ansiblego/pkg/transport"
)

type SetHostFact struct {
	Facts map[string]interface{}
}

func LoadSetHostFact(args map[string]interface{}) Module {
	return &SetHostFact{Facts: args}
}

func(c *SetHostFact) Run(ctx Context, transport transport.Transport, vars map[string]interface{}) *ModuleExecResult {
	for k,v := range c.Facts {
		vars[k] = v
	}
	return &ModuleExecResult{ Result: true }
}
