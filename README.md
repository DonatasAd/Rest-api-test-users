\$ Build with power shell

go to root project folder
$env:GOOS = "linux"
$env:CGO_ENABLED = "0"
$env:GOARCH = "amd64"
go build -o build\main main.go
cd build
~\Go\Bin\build-lambda-zip.exe -output main.zip main
