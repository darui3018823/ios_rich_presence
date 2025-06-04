# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Intel Chipset (macOS amd64)..."
Write-Host "Removing old executables..."
Remove-Item ./dist/mac/ios_shortcut_rpc_serv_amd64mac -Force

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o ./dist/mac/ios_shortcut_rpc_serv_amd64mac main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Build for Intel Chipset (macOS amd64) complete."
Write-Host "Server run command: ./dist/mac/ios_shortcut_rpc_serv_amd64mac"