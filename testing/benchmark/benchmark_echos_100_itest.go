package basic

import (
	"testing"
)

func TestEchos100GO(t *testing.T) {
	RunGosible(t, &BenchmarkConfig{PlaybookName: "test_echos_100.yaml"})
}

func TestEchos100Ansible(t *testing.T) {
	RunAnsible(t,  &BenchmarkConfig{PlaybookName: "test_echos_100.yaml"})
}

