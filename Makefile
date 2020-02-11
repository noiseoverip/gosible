build:
	echo "Building"

build:
	go build

test-integration:
	# count=1 is idiomatic way to disable cashing, all that needs to be done is to set an unknown flag
	go test --count 1 --tags integration ./...

test:
	go test --count=1  ./...

benchmark:


