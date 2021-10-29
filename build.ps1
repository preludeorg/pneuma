param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
go build -ldflags="-s -w -buildid= -X main.randomHash=$randomHash" -o payloads/pneuma-windows.exe main.go;
$env:GOOS='darwin';
go build -ldflags="-s -w -buildid= -X main.randomHash=$randomHash" -o payloads/pneuma-darwin main.go;
$env:GOOS='linux';
go build -ldflags="-s -w -buildid= -X main.randomHash=$randomHash" -o payloads/pneuma-linux main.go;

$env:GOOS='windows';
$env:CGO_ENABLED=1;
go build --buildmode=c-shared --ldflags="-s -w -buildid= -X main.randomHash=${1}" -o payloads/pneuma-windows.dll library/library.go;