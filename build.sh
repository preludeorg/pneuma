#!/bin/bash
GOOS=darwin go build -ldflags="-s -w -X main.randomHash=${1}" -o payloads/pneuma-darwin main.go;
GOOS=linux go build -ldflags="-s -w -X main.randomHash=${1}" -o payloads/pneuma-linux main.go;
GOOS=windows go build -ldflags="-s -w -X main.randomHash=${1}" -o payloads/pneuma-windows.exe main.go;
