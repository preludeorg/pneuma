param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o pneuma-windows.exe main.go;
$env:GOOS='darwin';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o pneuma-darwin main.go;
$env:GOOS='linux';
go build -ldflags="-s -w -X main.randomHash=$randomHash" -o pneuma-linux main.go;
