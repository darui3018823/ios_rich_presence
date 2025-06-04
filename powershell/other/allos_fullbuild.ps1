# Recommended PowerShell Core Version: 7.5.0 or later

Write-Host "This script will build binaries for the following devices in one batch:"
Write-Host " - Windows (amd64)"
Write-Host " - Windows (arm64)"
Write-Host " - macOS (Apple Silicon arm64)"
Write-Host " - macOS (Intel amd64)"
Write-Host " - Linux (amd64)"
Write-Host " - Linux (arm64)"
Write-Host " - Python executables"
Write-Host ""

$answer = Read-Host "Do you want to start the full build? (y/n)"
if ($answer -ne "y" -and $answer -ne "Y") {
    Write-Host "Build cancelled."
    exit
}

# Windows amd64
Write-Host "Building for Windows (amd64)..."
& ./powershell/win/win_build_amd64.ps1

# Windows arm64
Write-Host "Building for Windows (arm64)..."
& ./powershell/win/win_build_arm64.ps1

# macOS Apple Silicon
Write-Host "Building for macOS (Apple Silicon arm64)..."
& ./powershell/mac/mac_build_apple_silicon.ps1

# macOS Intel
Write-Host "Building for macOS (Intel amd64)..."
& ./powershell/mac/mac_build_intel_chipset.ps1

# Linux amd64
Write-Host "Building for Linux (amd64)..."
& ./powershell/linux/linux_build_amd64.ps1

# Linux arm64
Write-Host "Building for Linux (arm64)..."
& ./powershell/linux/linux_build_arm64.ps1

# Python executables
Write-Host "Building Python executables..."
& ./powershell/other/python_build.ps1

Write-Host ""
Write-Host "All builds complete!"