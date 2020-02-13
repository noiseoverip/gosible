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
	go test -v --count=1 --tags benchmark ./testing/benchmark/...

test-all: test test-integration
