package main

import (
	"ansiblego/testing/benchmark"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Benchmark CLI interface
	bVerbosity := flag.Int("v", 0, "Verbosity level")

	flag.Parse()

	if err := runBenchmark(*bVerbosity); err != nil {
		fmt.Print(err)
		os.Exit(1)
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
