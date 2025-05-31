# Recomended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Windows AMD64..."
Write-Host "Removing old executables..."
Remove-Item ./dist/win/ios_shortcut_rpc_serv_amd64win.exe -Force

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./dist/win/ios_shortcut_rpc_serv_amd64win.exe main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Building iOS Shortcut RPC Server for Windows AMD64 complete."
Write-Host "Server run command: ./dist/win/ios_shortcut_rpc_serv_amd64win.exe"