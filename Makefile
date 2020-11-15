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

## integration testing

.PHONY: integration-cleanup
integration-cleanup:
	docker-compose down

.PHONY: integration
integration: integration-cleanup
	docker-compose up --build --abort-on-container-exit --remove-orphans integration $${BUBBLY_PROVIDER}

.PHONY: storefront
storefront:
	docker-compose up --build --abort-on-container-exit --remove-orphans storefront $${BUBBLY_PROVIDER}

## local ci

# Project is CI-enabled with Github Actions. You can run CI locally
# using act (https://github.com/nektos/act). 
# There are some caveats, but the following target should work:
act: 
	act -P ubuntu-latest=golang:latest --env-file act.env -j simple
	