#!/usr/bin/env bash

set -ex
go run ./.github/cmd/generate_readme/main.go

cat README.md
