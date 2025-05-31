mkdir ./dist/win/
Remove-Item ./dist/win/ios_shortcut_rpc_serv_armwin.exe -Force
$env:GOOS = "windows"
$env:GOARCH = "arm"
go build -o ./dist/win/ios_shortcut_rpc_serv_armwin.exe main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH