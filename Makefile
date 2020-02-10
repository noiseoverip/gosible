build:
	echo "Building"

build:
	go build

test-integration:
	go test --tags integration ./...

test:
	go test ./...

benchmark:


