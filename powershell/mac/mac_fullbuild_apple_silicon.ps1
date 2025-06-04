# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Apple Silicon (macOS arm64)..."
Write-Host "Removing old executables..."
Remove-Item ./dist/mac/ios_shortcut_rpc_serv_arm64mac -Force

$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o ./dist/mac/ios_shortcut_rpc_serv_arm64mac main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Building iOS Shortcut RPC Server for Apple Silicon (macOS arm64) complete."

Write-Host "Removing old Python executables..."
Remove-Item ./python/set_rpc -Force
Remove-Item ./python/del_rpc -Force

Write-Host "Building Python executables..."
make-exe ./python/set_rpc.py --output ./python/
make-exe ./python/del_rpc.py --output ./python/

Write-Host "Building Python executables complete."
Write-Host "Removing Temporary Files..."
Write-Host "Please press A key"
Remove-Item ./build/ -Force
Remove-Item ./set_rpc.spec -Force
Remove-Item ./del_rpc.spec -Force
Write-Host "Building complete for Apple Silicon (macOS arm64) platform."
Write-Host "Server run command: ./dist/mac/ios_shortcut_rpc_serv_arm64mac"