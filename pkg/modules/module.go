package modules

import (
	"ansiblego/pkg/transport"
	"fmt"
)

type Context struct {
	PlaybookDir string
	InventoryDir string
}

// Module represents Ansible module interface
type Module interface {
	Run(ctx Context, transport transport.Transport, vars map[string]interface{}) *ModuleExecResult
}

func ErrorModuleConfig(text string, args ...interface{}) *ModuleExecResult {
	return &ModuleExecResult{Result: false, StdOut: "", StdErr: fmt.Sprintf(text, args...)}
}
