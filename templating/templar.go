package templating

import (
	"fmt"
	"github.com/flosch/pongo2"
)

// Optimization: we could tag args which do not contain variables during initialization and then skip
// template execution for them during runtime
func TemplateExec(input string, vars map[string]interface{}) (output string,  err error) {
	fmt.Printf("templar input: %s\n", input)

	tpl, err := pongo2.FromString(input)
	if err != nil {
		panic(err)
	}
	output, err = tpl.Execute(pongo2.Context(vars))
	if err != nil {
		panic(err)
	}

	fmt.Printf("templar output: %s\n", output)
	return
}
