package templating

import (
	"ansiblego/pkg/logging"
	"fmt"
	"github.com/flosch/pongo2"
	"strings"
)

// Optimization: we could tag args which do not contain variables during initialization and then skip
// template execution for them during runtime
func TemplateExec(input string, vars map[string]interface{}) (output string, err error) {
	logging.Debug("templar input: %s\n", input)

	tpl, err := pongo2.FromString(input)
	if err != nil {
		panic(err)
	}
	output, err = tpl.Execute(vars)
	if err != nil {
		panic(err)
	}

	logging.Debug("templar output: %s\n", output)
	return output, nil
}

func Assert(condition string, vars map[string]interface{}) (bool, error) {
	conditional := fmt.Sprintf("{%% if %s %%} True {%% else %%} False {%% endif %%}", condition)
	if outRaw, err := TemplateExec(conditional, vars); err == nil {
		out := strings.TrimSpace(outRaw)
		return out == "True", nil
	} else {
		return false, fmt.Errorf("failed to evaluate condition: %v", err)
	}
}
