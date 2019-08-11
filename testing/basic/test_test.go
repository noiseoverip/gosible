package basic

import (
	"ansiblego/runner"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRunLocalPlaybook(t *testing.T) {
	setup(t)
	wd, _ := os.Getwd()
	r := runner.Runner{ InventoryFilePath: path.Join(wd, "hosts"), PlaybookFilePath: path.Join(wd,"site.yaml") }
	err := r.Run()
	assert.NoError(t, err)
}

func setup(t *testing.T) {
	err := os.Chdir("files")
	assert.Nil(t, err)
	dir, err := os.Getwd()
	assert.Nil(t, err)
	t.Logf("current dir %s", dir)
}