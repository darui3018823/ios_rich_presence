# build.ps1 - iOS ShortCut DiscordRPC Server, All-in-One Build Script
# Recommended PowerShell Core Version: 7.5.0 or later

# 2025 iOS ShortCut DiscordRPC Server: darui3018823 All rights reserved.
# All works created by darui3018823 associated with this repository are the intellectual property of darui3018823.
# Packages and other third-party materials used in this repository are subject to their respective licenses and copyrights.

param (
    [switch]$Full,
    [switch]$Auto
)

function Get-Platform {
    $os = ""
    $arch = ""

    if ($IsWindows) { $os = "windows" }
    elseif ($IsMacOS) { $os = "darwin" }
    elseif ($IsLinux) { $os = "linux" }
    else {
        Write-Host "❌ Unsupported OS platform." -ForegroundColor Red
        exit 1
    }

    switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { $arch = "amd64" }
        "ARM64" { $arch = "arm64" }
        "x86"   { $arch = "386" }
        default {
            Write-Host "❌ Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)" -ForegroundColor Red
            exit 1
        }
    }

    return @{ os = $os; arch = $arch }
}

function Read-WithAuto($label, $options, $defaultDetector) {
    Write-Host "$label (Enterキーで自動検出)" -ForegroundColor Yellow
    foreach ($opt in $options.GetEnumerator()) {
        Write-Host "[$($opt.Key)] $($opt.Value)"
    }
    $choice = Read-Host "番号を入力してください"

    if ([string]::IsNullOrWhiteSpace($choice)) {
        return $defaultDetector.Invoke()
    }

    if ($options.ContainsKey($choice)) {
        return $options[$choice]
    }

    Write-Host "❌ 無効な入力です。" -ForegroundColor Red
    exit 1
}

function Build-GoServer($os, $arch) {
    $output = "./dist/$os/ios_shortcut_rpc_serv_${arch}${os}.exe"
    Write-Host "`n🔨 Building Go Server for $os $arch..." -ForegroundColor Green
    Remove-Item $output -Force -ErrorAction SilentlyContinue

    $env:GOOS = $os
    $env:GOARCH = $arch
    go build -o $output main.go
    Remove-Item Env:GOOS
    Remove-Item Env:GOARCH

    Write-Host "✅ Go Server build complete: $output" -ForegroundColor Green
}

function Build-PythonTools {
    Write-Host "`n🐍 Building Python executables..." -ForegroundColor Green
    Remove-Item ./python/set_rpc.exe, ./python/del_rpc.exe -Force -ErrorAction SilentlyContinue

    make-exe ./python/set_rpc.py --output ./python/
    make-exe ./python/del_rpc.py --output ./python/

    Remove-Item ./build/ -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item ./set_rpc.spec, ./del_rpc.spec -Force -ErrorAction SilentlyContinue

    Write-Host "✅ Python build complete." -ForegroundColor Green
}

# ====== 実行開始 ======

Write-Host "`n💡 iOS Shortcut RPC Server Build Tool" -ForegroundColor Cyan

if ($Auto) {
    $platform = Get-Platform
    $os = $platform.os
    $arch = $platform.arch
} else {
    $defaultOs = { (Get-Platform).os }
    $defaultArch = { (Get-Platform).arch }

    $os = Read-WithAuto "対象OSを選択してください" ([ordered]@{
        "1" = "windows"
        "2" = "linux"
        "3" = "darwin"
    }) $defaultOs

    $arch = Read-WithAuto "アーキテクチャを選択してください" ([ordered]@{
        "1" = "amd64"
        "2" = "arm64"
    }) $defaultArch

    $defaultFull = { $false }
    $fullChoice = Read-WithAuto "Full Build（Python等含む）を行いますか？" ([ordered]@{
        "1" = $true
        "2" = $false
    }) $defaultFull


    $Full = $fullChoice
}

Build-GoServer $os $arch

if ($Full) {
    Build-PythonTools
}

Write-Host "`n🎉 ビルド完了: $os / $arch" -ForegroundColor Cyan
Write-Host "▶ 実行コマンド: ./dist/$os/ios_shortcut_rpc_serv_${arch}${os}.exe" -ForegroundColor Yellow
