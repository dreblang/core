#!/bin/bash

go test -cover -coverpkg=./... -coverprofile=all.cov -v ./...
go tool cover -html all.cov
