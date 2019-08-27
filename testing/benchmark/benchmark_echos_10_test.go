package basic

import (
	"testing"
)

func TestEchos10GO(t *testing.T) {
	RunGosible(t, &BenchmarkConfig{PlaybookName: "test_echos_10.yaml"})
}

func TestEchos10Ansible(t *testing.T) {
	RunAnsible(t,  &BenchmarkConfig{PlaybookName: "test_echos_10.yaml"})
}

