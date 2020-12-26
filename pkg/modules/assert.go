package modules

import (
	"ansiblego/pkg/templating"
)

type Assert struct {
	That       string // that
	FailMsg    string // fail_msg
	SuccessMsg string // success_msg
}

func LoadAssert(args map[string]string) Module {
	return &Assert{That: args["that"], FailMsg: args["success_msg"], SuccessMsg: args["success_msg"]}
}

func (self *Assert) Run(ctx Context, host *Host) *ModuleExecResult {
	// Since variable change during runtime, we have to render args at the point of execution
	renderedCondition, err := templating.Assert(self.That, host.Vars)
	if err != nil {
		return &ModuleExecResult{Result: false, StdOut: "", StdErr: err.Error()}
	}
	message := "Assertion failed"
	if self.FailMsg != "" {
		message = self.FailMsg
	}
	if renderedCondition {
		message = "Assertion passed"
		if self.SuccessMsg != "" {
			message = self.SuccessMsg
		}
	}
	return &ModuleExecResult{Result: renderedCondition, StdOut: message, StdErr: ""}
}
