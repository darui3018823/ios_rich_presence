Remove-Item ./dist/other/analysis_header_amd64win.exe -Force
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./dist/other/analysis_header_amd64win.exe ./.temp/analysis_header.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH