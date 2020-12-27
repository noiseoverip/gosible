build:
	go build

.PHONY: test
test:
	go test ./...

.PHONY: testint
testint:
	go test --count 1 --tags integration ./testing/integration/...

.PHONY: bench
bench:
	# Run few different playbooks with gosible and ansible for comparison
	go run cmd/gosible-bench/main.go --host ${HOST}

.PHONY: format
format: vendor
	gofmt -s -w .
	golangci-lint run --fix

test-all: test testint bench
