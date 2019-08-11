package main

import (
	"ansiblego/runner"
	"flag"
	"fmt"
	"os"
	"path"
)

var inventoryPath = flag.String("i", "", "Path to inventory")


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

	r := &runner.Runner{PlaybookFilePath: path.Join(cwd, playbookPath), InventoryFilePath: path.Join(cwd, *inventoryPath) }
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
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}