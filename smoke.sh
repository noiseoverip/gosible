#!/bin/bash
# This is a quick smoke test which invokes a playbook from testing/basic package

go run main.go playbook -i testing/integration/files/hosts testing/integration/files/site.yaml -v 3