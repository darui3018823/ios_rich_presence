Remove-Item ./dist/win/ios_shortcut_rpc_serv_amd64win.exe -Force
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o ./dist/win/ios_shortcut_rpc_serv_amd64win.exe main.go
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Remove-Item ./python./set_rpc.exe -Force
Remove-Item ./python./del_rpc.exe -Force
make-exe ./python/set_rpc.py --output ./python/
make-exe ./python/del_rpc.py --output ./python/
Remove-Item ./build/ -Force
Remove-Item ./set_rpc.spec -Force
Remove-Item ./del_rpc.spec -Force