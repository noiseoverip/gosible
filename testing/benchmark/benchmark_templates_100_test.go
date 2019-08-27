package basic

import (
	"testing"
)

func TestTemplates100GO(t *testing.T) {
	RunGosible(t, &BenchmarkConfig{PlaybookName: "test_templates_100.yaml"})
}

func TestTemplates100Ansible(t *testing.T) {
	RunAnsible(t,  &BenchmarkConfig{PlaybookName: "test_templates_100.yaml"})
}

