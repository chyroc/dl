#!/usr/bin/env bash

rm -rf ~/Downloads/dl-download.app || echo ""
mkdir -p ~/Downloads/dl-download.app/Contents/MacOS
go build -o ~/Downloads/dl-download.app/Contents/MacOS/dl-download main.go

open ~/Downloads/dl-download.app # Or click on the app in Finder
