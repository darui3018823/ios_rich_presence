# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Apple Silicon (macOS arm64)..."
Write-Host "Removing old executables..."
Remove-Item ./dist/mac/ios_shortcut_rpc_serv_arm64mac -Force

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o ./dist/mac/ios_shortcut_rpc_serv_arm64mac main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Build for Apple Silicon (macOS arm64) complete."
Write-Host "Server run command: ./dist/mac/ios_shortcut_rpc_serv_arm64mac"