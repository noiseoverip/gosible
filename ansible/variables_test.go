package ansible

import (
	"fmt"
	"testing"
)

func TestLoadGroupVars(t *testing.T) {
	vars, err := load_group_vars("./testfiles")
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Log(fmt.Printf("%s", vars))
	// TODO: this is working, add asserts and it will be OK for now
}