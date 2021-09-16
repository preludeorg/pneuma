#!/bin/bash
GOOS=darwin go build -ldflags="-s -w -X main.randomHash=${1}" -o pneuma-darwin main.go;
GOOS=linux go build -ldflags="-s -w -X main.randomHash=${1}" -o pneuma-linux main.go;
GOOS=windows go build -ldflags="-s -w -X main.randomHash=${1}" -o pneuma-windows.exe main.go;
