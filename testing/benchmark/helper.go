package basic

import (
	"ansiblego/runner"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"testing"
)
type BenchmarkConfig struct {
	PlaybookName string
}

func RunGosible(t *testing.T, config *BenchmarkConfig) {
	setup(t)
	wd, _ := os.Getwd()
	r := runner.Runner{ InventoryFilePath: path.Join(wd, "hosts"), PlaybookFilePath: path.Join(wd, config.PlaybookName) }
	err := r.Run()
	assert.NoError(t, err)
}

func RunAnsible(t *testing.T, config *BenchmarkConfig) {
	setup(t)
	cmd := exec.Command("ansible-playbook",  "-i", "hosts", config.PlaybookName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		t.Errorf("%v", err)
		t.Fatalf("Failed")
	}
}

func setup(t *testing.T) {
	err := os.Chdir("files")
	assert.Nil(t, err)
	dir, err := os.Getwd()
	assert.Nil(t, err)
	t.Logf("current dir %s", dir)
}
