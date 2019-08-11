package modules

import (
	"ansiblego/transport"
)

// Module represents Ansible module interface
type Module interface {
	Run(transport transport.Transport, vars map[string]interface{}) *ModuleExecResult
}