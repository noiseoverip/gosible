package main

import (
	"ansiblego/pkg"
	"ansiblego/pkg/logging"
	"ansiblego/testing/benchmark"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func runPlaybook(inventory string, playbook string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r := pkg.NewRunner(path.Join(cwd, inventory), path.Join(cwd, playbook))
	return r.Run()
}

func main() {
	// Playbook CLI interface
	playbookCommand := flag.NewFlagSet("playbook", flag.ExitOnError)
	inventoryPath := playbookCommand.String("i", "", "Path to inventory")
	pVerbosity := playbookCommand.Int("v", 0, "Verbosity level")
	playbookCommand.Usage = func() {
		fmt.Println("Usage: playbook [options] playbook.yml")
		playbookCommand.PrintDefaults()
	}

	// Benchmark CLI interface
	benchmarkCommand := flag.NewFlagSet("benchmark", flag.ExitOnError)
	bVerbosity := benchmarkCommand.Int("v", 0, "Verbosity level")

	var usage = func() {
		fmt.Println("gosible [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("\tplaybook - run ansible playbook")
		fmt.Println("\tbenchmark - run benchmark to compare gosible vs ansible. Currently only works from inside source repo")
	}

	if len(os.Args) == 1 {
		usage()
	}

	switch os.Args[1] {
	case "playbook":
		playbookCommand.Parse(os.Args[2:])
	case "benchmark":
		benchmarkCommand.Parse(os.Args[2:])
	case "help":
		fallthrough
	case "--help":
		usage()
	default:
		fmt.Printf("Invalid subcommand %s.\n", os.Args[1])
	}

	if playbookCommand.Parsed() {
		if playbookCommand.NArg() < 1 {
			fmt.Println("Error: Last argument should be path to playbook")
			fmt.Println()
			playbookCommand.Usage()
			os.Exit(1)
		}
		if *inventoryPath == "" {
			fmt.Println("Error: path to inventory must be provided")
			fmt.Println()
			playbookCommand.Usage()
			os.Exit(1)
		}

		if *pVerbosity > 0 {
			logging.Global = logging.NewGosibleVerboseLogger(*pVerbosity)
		}

		logging.Global = logging.NewGosibleVerboseLogger(*pVerbosity)

		if err := runPlaybook(*inventoryPath, playbookCommand.Arg(0)); err != nil {
			fmt.Printf("Failure during playbook execution: %s", err)
			os.Exit(1)
		}
	}

	if benchmarkCommand.Parsed() {
		if err := runBenchmark(*bVerbosity); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}
}

func runBenchmark(verbosity int) error {
	log.Printf("Benchmark START")

	tests := []*benchmark.BenchmarkConfig{
		{PlaybookName: "test_echos_10.yaml", ExpectedMaxDurationSec: 20, Verbose: verbosity},
		{PlaybookName: "test_echos_100.yaml", ExpectedMaxDurationSec: 20, Verbose: verbosity},
		{PlaybookName: "test_templates_10.yaml", ExpectedMaxDurationSec: 40, Verbose: verbosity},
	}

	for _, tool := range []func(c *benchmark.BenchmarkConfig) error{
		benchmark.RunGosible,
		benchmark.RunAnsible,
	} {
		for _, tt := range tests {
			err := tool(tt)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
