package modules

type SetHostFact struct {
	Facts map[string]interface{}
}

func LoadSetHostFact(args map[string]interface{}) Module {
	return &SetHostFact{Facts: args}
}

func (c *SetHostFact) Run(ctx Context, host *Host) *ModuleExecResult {
	for k, v := range c.Facts {
		host.Vars[k] = v
	}
	return &ModuleExecResult{Result: true}
}
