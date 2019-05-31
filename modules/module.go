package modules

import "ansiblego/transport"

// Module represents Ansible module interface
type Module interface {
	//  TODO: module interface should provide means of
	Run(transport transport.Transport) *ModuleExecResult
}