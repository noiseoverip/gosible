package basic

import (
	"testing"
)

func TestTemplates10GO(t *testing.T) {
	RunGosible(t, &BenchmarkConfig{PlaybookName: "test_templates_10.yaml"})
}

func TestTemplates10Ansible(t *testing.T) {
	RunAnsible(t,  &BenchmarkConfig{PlaybookName: "test_templates_10.yaml"})
}

