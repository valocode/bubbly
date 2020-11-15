#!/bin/bash

# hack: Let the DB boot up
sleep 1

go test -race -v --tags=integration ./integration
