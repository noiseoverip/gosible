package modules

import (
	"ansiblego/internal/templating"
	"strings"
)

// Command implements module interface and executes CLI commands on transport layer
// Pipes are not supported at this point
type Command struct {
	Input string
}

func LoadCommand(args map[string]string) Module {
	return &Command{Input: args["stdin"]}
}

func (c *Command) Run(_ Context, host *Host) *ModuleExecResult {
	// Since variables change during runtime, we have to render args at the point of execution
	renderedArgs, err := templating.TemplateExec(c.Input, host.Vars)
	if err != nil {
		return &ModuleExecResult{Result: false, StdOut: "", StdErr: err.Error()}
	}
	cmd := strings.Split(renderedArgs, " ")
	resultCode, stdout, stderr, err := host.Transport.Exec(cmd[0], cmd[1:]...)
	if err != nil {
		return &ModuleExecResult{Result: false, StdOut: "", StdErr: err.Error()}
	}
	return &ModuleExecResult{Result: resultCode == 0, StdOut: stdout, StdErr: stderr}
}
