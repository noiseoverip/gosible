package basic

import (
	"ansiblego/pkg"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

type BenchmarkConfig struct {
	PlaybookName string
}

func RunGosible(t *testing.T, config *BenchmarkConfig) {
	start := time.Now().Nanosecond()
	wd, _ := os.Getwd()
	r := pkg.Runner{
		Context: &pkg.Context{
			InventoryFilePath: path.Join(wd, "files", "hosts"),
			PlaybookFilePath:  path.Join(wd, "files", config.PlaybookName),
		},
		Strategy: &pkg.SequentialExecuter{},
	}
	err := r.Run()
	assert.NoError(t, err)
	t.Logf("Duration %d", time.Now().Nanosecond()-start)
}

func RunAnsible(t *testing.T, config *BenchmarkConfig) {
	wd, _ := os.Getwd()
	hostsPath := path.Join(wd, "files", "hosts")
	playbookPath := path.Join(wd, "files", config.PlaybookName)
	cmd := exec.Command("ansible-playbook", "-i", hostsPath, playbookPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		t.Errorf("%v", err)
		t.Fatalf("Failed")
	}
}
