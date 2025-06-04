Write-Host "Removing old Python executables..."
Remove-Item ./python/set_rpc.exe -Force
Remove-Item ./python/del_rpc.exe -Force

Write-Host "Building Python executables..."
make-exe ./python/set_rpc.py --output ./python/
make-exe ./python/del_rpc.py --output ./python/

Write-Host "Building Python executables complete."
Write-Host "Removing Temporary Files..."
Write-Host "Please press A key"
Remove-Item ./build/ -Force
Remove-Item ./set_rpc.spec -Force
Remove-Item ./del_rpc.spec -Force
Write-Host "Temporary Files removed."
Write-Host "Python executables build complete."