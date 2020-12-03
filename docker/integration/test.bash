#!/bin/bash

# hack: Let the DB boot up
sleep 1

go test -race -v --tags=integration ./integration

# in the future we can create a report
# go get -u github.com/jstemmer/go-junit-report
# go test -race -v --tags=integration ./integration -bench . -count 5 | tee go-junit-report >./integration/testdata/integration-report.xml
