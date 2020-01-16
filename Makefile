build:
	echo "Building"

test-integration:
	go test ansiblego/testing/basic

test: test-integration

