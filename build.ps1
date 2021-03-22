param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o payloads/pneuma-windows.exe main.go;
$env:GOOS='darwin';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o payloads/pneuma-darwin main.go;
$env:GOOS='linux';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o payloads/pneuma-linux main.go;