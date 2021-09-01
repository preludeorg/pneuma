param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-windows.exe main.go;
$env:GOOS='darwin';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-darwin main.go;
$env:GOOS='linux';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-linux main.go;