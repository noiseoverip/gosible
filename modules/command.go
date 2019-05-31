package modules

import "ansiblego/transport"

type Command struct {
	Input string
}

func LoadCommand(args map[string]string) Module {
	return &Command{Input: args["stdin"]}

}

func(c *Command) Run(transport transport.Transport) *ModuleExecResult {
	transport.Exec("echo", "LABAS")
	return &ModuleExecResult{ Result: true, StdOut: "AllGood"}
}
