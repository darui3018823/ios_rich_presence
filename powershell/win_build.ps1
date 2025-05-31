mkdir ./dist/win/
Remove-Item ./dist/win/ios_shortcut_rpc_serv.exe -Force
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./dist/win/ios_shortcut_rpc_serv.exe main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH