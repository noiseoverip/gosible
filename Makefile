build:
	echo "Building"

build:
	go build

test:
	go test ./...

test-integration:
	go test --count 1 --tags integration ./testing/integration/...

test-benchmark:
	# Run ansiblego benchmark test
	go run main.go -b

test-all: test test-integration
