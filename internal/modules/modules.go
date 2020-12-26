package modules

import (
	"ansiblego/internal/transport"
	"fmt"
)

type Context struct {
	PlaybookDir  string
	InventoryDir string
}

// Module represents Ansible module interface
type Module interface {
	Run(ctx Context, host *Host) *ModuleExecResult
}

type Host struct {
	Vars      map[string]interface{}
	Transport transport.Transport
}

func ErrorModuleConfig(text string, args ...interface{}) *ModuleExecResult {
	return &ModuleExecResult{Result: false, StdOut: "", StdErr: fmt.Sprintf(text, args...)}
}

// ModuleExecResult holds any module result
type ModuleExecResult struct {
	Result bool
	StdOut string
	StdErr string
}

func (r *ModuleExecResult) String() string {
	return fmt.Sprintf("Result:%v\n\tStdout:\n%s\n\tStdErr:\n%s\n", r.Result, r.StdOut, r.StdErr)
}
