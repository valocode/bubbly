BIN=./build/bubbly
KIND_CLUSTER_NAME=bubbly

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

## runing a dev instance
.PHONY: dev
dev:
	skaffold dev

## integration testing

.PHONY: kind-cleanup
kind-cleanup:
	kind delete cluster --name ${KIND_CLUSTER_NAME}

.PHONY: kind-bootstrap
kind-bootstrap:
	# create the kind cluster
	kind create cluster --name ${KIND_CLUSTER_NAME} --config kind-config.yaml

.PHONY: development
development:
	# --trigger manual means you need to press ENTER to trigger a re-build/deploy
	skaffold dev --trigger manual

.PHONY: integration
integration:
	# --force so that the k8s Job gets re-applied
	skaffold run -p integration --force
	# print the logs and follow until complete
	kubectl logs --follow --tail=-1 --selector job-name=bubbly-integration

# Project is CI-enabled with Github Actions. You can run CI locally
# using act (https://github.com/nektos/act). 
# There are some caveats, but the following target should work:
act: 
	act -P ubuntu-latest=golang:latest --env-file act.env -j simple
	