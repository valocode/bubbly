BIN=./build/bubbly

all: build run-help

.PHONY: build
build:
	go build -o ${BIN}

run-help: 
	${BIN} -h

## testing

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

display-coverage: test-coverage
	go tool cover -html=coverage.txt

test-report:
	go test -coverprofile=coverage.txt -covermode=atomic -json ./... > test_report.json
