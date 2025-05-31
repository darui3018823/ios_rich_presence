mkdir ./dist/win/
Remove-Item ./dist/other/analysis_header_amd64win.exe -Force
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./dist/win/analysis_header_amd64win.exe main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH