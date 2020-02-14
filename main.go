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

var inventoryPath = flag.String("i", "", "Path to inventory")
var verbosity = flag.Int("v", 0, "Verbosity")
var benchmarkTest = flag.Bool("b", false, "Benchmark test")

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

	context := pkg.Context{
		PlaybookFilePath:  path.Join(cwd, playbookPath),
		InventoryFilePath: path.Join(cwd, *inventoryPath),
	}
	r := &pkg.Runner{Context: &context}
	err = r.Run()
	if err != nil {
		return fmt.Errorf("runner error: %v", err)
	}
	return nil
}

func runBenchmark() {
	log.Printf("Benchmark START")

	errGosible := benchmark.RunGosible(&benchmark.BenchmarkConfig{
		PlaybookName: "test_echos_10.yaml",
		ExpectedMaxDurationSec: 2,
		Verbose: false})
	if errGosible != nil {
		panic(errGosible)
	}

	errAnsible := benchmark.RunAnsible(&benchmark.BenchmarkConfig{
		PlaybookName: "test_echos_10.yaml",
		ExpectedMaxDurationSec: 10,
		Verbose: false,
	})
	if errAnsible != nil {
		panic(errAnsible)
	}

	log.Printf("Benchmark DONE")
}
//
// Usage: ansiblego -i inventory site.yml
//
func main() {
	flag.Parse()
	if *benchmarkTest {
		runBenchmark()
	} else {
		err := run()
		if err != nil {
			logging.Info(err.Error())
			os.Exit(1)
		}
	}
}
