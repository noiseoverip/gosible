package main

import (
	"ansiblego/logging"
	"ansiblego/runner"
	"flag"
	"fmt"
	"os"
	"path"
)

var inventoryPath = flag.String("i", "", "Path to inventory")
var verbosity = flag.Int("v", 0, "Verbosity")

func run() error {
	flag.Parse()
	if len(flag.Args()) < 1 {
		return fmt.Errorf("too few arguments provided, please provide path to playbook") //TODO: show to print to stderr
	}

	playbookPath := flag.Arg(0)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	context := runner.Context{
		PlaybookFilePath:  path.Join(cwd, playbookPath),
		InventoryFilePath: path.Join(cwd, *inventoryPath),
	}
	r := &runner.Runner{Context: &context}
	err = r.Run()
	if err != nil {
		return fmt.Errorf("runner error: %v", err)
	}
	return nil
}

//
// Usage: ansiblego -i inventory site.yml
//
func main() {
	err := run()
	if err != nil {
		logging.Info(err.Error())
		os.Exit(1)
	}
}
