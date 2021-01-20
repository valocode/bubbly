BIN=./build/bubbly

# env vars for running tests
#export BUBBLY_HOST=localhost
#export BUBBLY_PORT=80
export BUBBLY_STORE_PROVIDER=postgres
#export POSTGRES_ADDR=localhost:30000
#export POSTGRES_DATABASE=bubbly

all: build

.PHONY: build
build:
	go build -o ${BIN}

.PHONY: clean
clean:
	rm ${BIN}

## testing

test: test-unit

test-unit: build
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

display-coverage: test-coverage
	go tool cover -html=coverage.txt

test-report:
	go test -coverprofile=coverage.txt -covermode=atomic -json ./... > test_report.json

test-integration: build test-unit
	docker-compose up --build --abort-on-container-exit --remove-orphans integration $${BUBBLY_STORE_PROVIDER}

# TODO test-integration-docker-down or docker-compose run integration?

## runing a dev instance
#.PHONY: dev
#dev:
#	skaffold dev --trigger manual

## integration testing

#.PHONY: kind-cleanup
#kind-cleanup:
#	kind delete cluster --name ${KIND_CLUSTER_NAME}

#.PHONY: kind-bootstrap
#kind-bootstrap:
	# create the kind cluster
#	kind create cluster --name ${KIND_CLUSTER_NAME} --config kind-config.yaml

#.PHONY: integration
#integration:
#	go test ./integration -tags=integration

# Project is CI-enabled with Github Actions. You can run CI locally
# using act (https://github.com/nektos/act). 
# There are some caveats, but the following target should work:
act: 
	act -P ubuntu-latest=golang:latest --env-file act.env -j simple
	