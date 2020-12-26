package benchmark

import (
	"ansiblego/pkg"
	"ansiblego/pkg/logging"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	LOGGING_DEFAULT = iota // logs only from benchmark tool and errors from gosible/ansible
	LOGGING_BASIC          // normal logging from gosible/ansible
	LOGGING_VERBOSE        // verbose logging from gosible/ansible
)

type BenchmarkConfig struct {
	PlaybookName           string
	ExpectedMaxDurationSec int64
	Verbose                int
}

func RunGosible(config *BenchmarkConfig) error {
	switch config.Verbose {
	case LOGGING_DEFAULT:
		logging.Global = logging.NewGosibleSilentLogger()
	case LOGGING_BASIC:
		logging.Global = logging.NewGosibleDefaultLogger()
	case LOGGING_VERBOSE:
		logging.Global.SetVerbose(os.Stdout)
	}

	r := pkg.Runner{
		Context: &pkg.Context{
			InventoryFilePath: path.Join(resourcePath(), "hosts"),
			PlaybookFilePath:  path.Join(resourcePath(), config.PlaybookName),
		},
		Strategy: &pkg.ParalelExecutor{},
	}
	start := time.Now().Unix()
	err := r.Run()
	if err != nil {
		return err
	}
	duration := time.Now().Unix() - start
	log.Printf("gosible %s: %d seconds", config.PlaybookName, duration)
	if duration > config.ExpectedMaxDurationSec {
		log.Fatalf("\t Expected was: %d", config.ExpectedMaxDurationSec)
	}
	return nil
}

func RunAnsible(config *BenchmarkConfig) error {
	cmd := exec.Command("ansible-playbook", "-i", "hosts")
	cmd.Dir = resourcePath()
	ansibleVerbosity := 0
	if config.Verbose >= LOGGING_DEFAULT {
		cmd.Stderr = os.Stderr
	}
	if config.Verbose > LOGGING_DEFAULT {
		log.Printf("\t %s %s", cmd.Path, strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
	}
	if config.Verbose >= LOGGING_VERBOSE {
		ansibleVerbosity = 3
	}
	if ansibleVerbosity > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("-%s", strings.Repeat("v", ansibleVerbosity)))
	}
	// add path to playbook at the end
	cmd.Args = append(cmd.Args, config.PlaybookName)

	start := time.Now().Unix()
	if err := cmd.Run(); err != nil {
		return err
	}
	duration := time.Now().Unix() - start
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
