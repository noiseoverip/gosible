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
	# Run ansiblego benchmark test
	go run main.go benchmark


.PHONY: format
format: vendor
	gofmt -s -w .
	golangci-lint run --fix

test-all: test testint bench
