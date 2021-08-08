#!/bin/bash
GOOS=darwin garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-darwin main.go;
GOOS=linux garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-linux main.go;
GOOS=windows garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-windows.exe main.go;
