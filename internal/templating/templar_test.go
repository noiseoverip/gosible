package templating

import (
	"fmt"
	"github.com/flosch/pongo2"
	"strings"
	"testing"
)

func TestJinjaTemplateBasic(t *testing.T) {
	tpl, err := pongo2.FromString("Hello {{ name|capfirst }}!")
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(pongo2.Context{"name": "florian"})
	if err != nil {
		panic(err)
	}
	t.Logf("%s\n", out)
}

func TestJinjaTemplateCondition(t *testing.T) {
	conditional := fmt.Sprintf("{%% if %s %%} True {%% else %%} False {%% endif %%}", "testvar == 'florian' ")
	t.Log(conditional)
	tpl, err := pongo2.FromString(conditional)
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(pongo2.Context{"testvar": "florian"})
	if err != nil {
		panic(err)
	}
	if strings.TrimSpace(out) != "True" {
		t.Fatal()
	}
	t.Log(out)
}
