package main

import (
	"ansiblego/testing/benchmark"
	"ansiblego/testing/tools"
	"flag"
	"fmt"
	"log"
	"os"
)

var verbosity int
var targetHost string

func main() {
	// Benchmark CLI interface
	flag.IntVar(&verbosity, "v", 0, "Verbosity level")
	flag.StringVar(&targetHost, "host", "", "Target host IP address")
	flag.Parse()

	if err := runBenchmark(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func runBenchmark() error {
	log.Printf("Benchmark START")

	tests := []*benchmark.BenchmarkConfig{
		{PlaybookName: "test_echos_10.yaml", ExpectedMaxDurationSec: 20, Verbose: verbosity, TargetHostAddr: targetHost},
		{PlaybookName: "test_echos_100.yaml", ExpectedMaxDurationSec: 20, Verbose: verbosity, TargetHostAddr: targetHost},
		{PlaybookName: "test_templates_10.yaml", ExpectedMaxDurationSec: 40, Verbose: verbosity, TargetHostAddr: targetHost},
	}

	hostsFile := "testing/benchmark/files/hosts"
	err := tools.RenderHostsFile("testing/benchmark/files/hosts_template", hostsFile, targetHost)
	panicIfError(err)
	defer os.RemoveAll(hostsFile)

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

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
