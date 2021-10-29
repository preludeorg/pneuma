param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
$env:GOOS='windows';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-windows.exe main.go;
$env:GOOS='darwin';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-darwin main.go;
$env:GOOS='linux';
garble build -ldflags="-X main.randomHash=$randomHash" -o payloads/pneuma-linux main.go;

$env:GOOS='windows';
$env:CGO_ENABLED=1;
garble build --buildmode=c-shared --ldflags="-X main.randomHash=${1}" -o payloads/pneuma-windows.dll library/library.go;