# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Linux (arm)..."
Write-Host "Removing old executables..."
Remove-Item ./dist/linux/ios_shortcut_rpc_serv_armlinux -Force

$env:GOOS = "linux"
$env:GOARCH = "arm"
go build -o ./dist/linux/ios_shortcut_rpc_serv_armlinux main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Build for Linux (arm) complete."
Write-Host "Server run command: ./dist/linux/ios_shortcut_rpc_serv_armlinux"