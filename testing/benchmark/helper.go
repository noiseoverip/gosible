package benchmark

import (
	"ansiblego/pkg"
	"ansiblego/pkg/logging"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type BenchmarkConfig struct {
	PlaybookName string
	ExpectedMaxDurationSec int64
	Verbose bool
}

func RunGosible(config *BenchmarkConfig) error {
	if !config.Verbose {
		// Disable gosible logging, leave only error logs
		logging.L = logging.NewGosibleSilentLogger()
	}
	r := pkg.Runner{
		Context: &pkg.Context{
			InventoryFilePath: path.Join(resourcePath(), "hosts"),
			PlaybookFilePath:  path.Join(resourcePath(), config.PlaybookName),
		},
		Strategy: &pkg.SequentialExecuter{},
	}
	start := time.Now().Unix()
	err := r.Run()
	if err != nil {
		return err
	}
	duration := time.Now().Unix()-start
	log.Printf("gosible %s: %d seconds", config.PlaybookName, duration)
	if duration > config.ExpectedMaxDurationSec {
		log.Fatalf("\t Expected was: %d", config.ExpectedMaxDurationSec)
	}
	return nil
}

func RunAnsible(config *BenchmarkConfig) error {
	playbookPath := path.Join(resourcePath(), config.PlaybookName)
	cmd := exec.Command("ansible-playbook", "-i", path.Join(resourcePath(), "hosts"), playbookPath)
	if config.Verbose {
		log.Printf("\t %s %s", cmd.Path, strings.Join(cmd.Args, " ") )
	}
	cmd.Stderr = os.Stderr
	if config.Verbose {
		cmd.Stdout = os.Stdout
	}
	start := time.Now().Unix()
	if err := cmd.Run(); err != nil {
		return err
	}
	duration := time.Now().Unix()-start
	log.Printf("ansible %s: %d seconds", config.PlaybookName, duration)
	if duration > config.ExpectedMaxDurationSec {
		log.Fatalf("\t Expected was: %d", config.ExpectedMaxDurationSec)
	}
	return nil
}

func resourcePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, "testing", "benchmark", "files")
}
