package main

import (
	"ansiblego/internal"
	"ansiblego/internal/logging"
	"flag"
	"fmt"
	"os"
	"path"
)

func runPlaybook(inventory string, playbook string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r := internal.NewRunner(path.Join(cwd, inventory), path.Join(cwd, playbook))
	return r.Run()
}

func main() {
	// Playbook CLI interface
	inventoryPath := flag.String("i", "hosts", "Path to inventory")
	pVerbosity := flag.Int("v", 0, "Verbosity level")

	flag.Parse()

	if *pVerbosity > 0 {
		logging.Global = logging.NewGosibleVerboseLogger(*pVerbosity)
	}

	logging.Global = logging.NewGosibleVerboseLogger(*pVerbosity)

	if len(flag.Args()) > 1 {
		fmt.Printf("Only one positional arg is accepted")
		os.Exit(1)
	}

	if err := runPlaybook(*inventoryPath, flag.Arg(0)); err != nil {
		fmt.Printf("Failure during playbook execution: %s", err)
		os.Exit(1)
	}
}
