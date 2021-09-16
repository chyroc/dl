#!/usr/bin/env bash

set -ex
go run ./.github/cmd/generate_website/main.go $1
