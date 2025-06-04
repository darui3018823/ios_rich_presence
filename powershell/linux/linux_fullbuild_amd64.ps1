# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "Building iOS Shortcut RPC Server for Linux (amd64)..."
Write-Host "Removing old executables..."
Remove-Item ./dist/linux/ios_shortcut_rpc_serv_amd64linux -Force

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ./dist/linux/ios_shortcut_rpc_serv_amd64linux main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Write-Host "Building iOS Shortcut RPC Server for Linux (amd64) complete."

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
Write-Host "Building complete for Linux (amd64) platform."
Write-Host "Server run command: ./dist/linux/ios_shortcut_rpc_serv_amd64linux"