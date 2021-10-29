#!/bin/bash
GOOS=darwin garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-darwin main.go;
GOOS=linux garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-linux main.go;
GOOS=windows garble build -ldflags="-X main.randomHash=${1}" -o payloads/pneuma-windows.exe main.go;

command -v x86_64-w64-mingw32-gcc &> /dev/null && GOOS=windows CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-windows.dll library/library.go;

if [[ $(uname) == "Darwin" ]]
then
  command -v x86_64-linux-musl-gcc &> /dev/null && GOOS=linux CC=x86_64-linux-musl-gcc CGO_ENABLED=1 garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-linux.so library/library.go;
  CGO_ENABLED=1 garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-darwin.dylib library/library.go;
else
  CGO_ENABLED=1 garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-linux.so library/library.go;
  command -v o64-clang &> /dev/null && GOOS=darwin CC=o64-clang CGO_ENABLED=1 garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-linux.so library/library.go;
fi