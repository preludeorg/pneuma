param([String]$key="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
go build -ldflags="-s -w -X main.key=$key" -o payloads/pneuma-windows.exe main.go;
$env:GOOS='darwin';
go build -ldflags="-s -w -X main.key=$key" -o payloads/pneuma-darwin main.go;
$env:GOOS='linux';
go build -ldflags="-s -w -X main.key=$key" -o payloads/pneuma-linux main.go;