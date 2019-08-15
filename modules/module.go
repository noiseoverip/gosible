package modules

import (
	"ansiblego/transport"
	"fmt"
)

// Module represents Ansible module interface
type Module interface {
	Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult
}

func ErrorModuleConfig(text string, args ...interface{}) *ModuleExecResult {
	return &ModuleExecResult{Result: false, StdOut: "", StdErr: fmt.Sprintf(text, args...)}
}
