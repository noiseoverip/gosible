package modules

import "fmt"

// ModuleExecResult holds any module result
type ModuleExecResult struct {
	Result bool
	StdOut string
	StdErr string
}

func(r *ModuleExecResult) String() string {
	return fmt.Sprintf("Result:%v\nStdout:\n%s\nStdErr:\n%s\n", r.Result, r.StdOut, r.StdErr)
}