package modules

import (
	"ansiblego/internal/templating"
)

type Assert struct {
	That       string
	FailMsg    string
	SuccessMsg string
}

func LoadAssert(args map[string]string) Module {
	return &Assert{That: args["that"], FailMsg: args["success_msg"], SuccessMsg: args["success_msg"]}
}

func (a *Assert) Run(ctx Context, host *Host) *ModuleExecResult {
	// Since variable change during runtime, we have to render args at the point of execution
	renderedCondition, err := templating.Assert(a.That, host.Vars)
	if err != nil {
		return &ModuleExecResult{Result: false, StdOut: "", StdErr: err.Error()}
	}
	message := "Assertion failed"
	if a.FailMsg != "" {
		message = a.FailMsg
	}
	if renderedCondition {
		message = "Assertion passed"
		if a.SuccessMsg != "" {
			message = a.SuccessMsg
		}
	}
	return &ModuleExecResult{Result: renderedCondition, StdOut: message, StdErr: ""}
}
