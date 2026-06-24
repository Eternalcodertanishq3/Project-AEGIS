# Download and extract Kiwix tools for Windows
Write-Host "AEGIS Kiwix Engine Auto-Installer"
Write-Host "================================="
Write-Host "This script will download the official Kiwix Tools for Windows to fix missing DLLs."

$url = "https://download.kiwix.org/release/kiwix-tools/kiwix-tools_win-x86_64-3.8.1.zip"
$zipPath = "kiwix-tools.zip"
$targetDir = "$PSScriptRoot\sidecars\kiwix-serve\windows"

try {
    Write-Host "`nDownloading Kiwix tools... (This may take a minute)"
    Invoke-WebRequest -Uri $url -OutFile $zipPath

    Write-Host "Extracting..."
    Expand-Archive -Path $zipPath -DestinationPath "kiwix-tools-temp" -Force

    Write-Host "Installing to $targetDir..."
    if (!(Test-Path $targetDir)) {
        New-Item -ItemType Directory -Force -Path $targetDir | Out-Null
    }
    
    # Copy all DLLs and the EXE
    Copy-Item -Path "kiwix-tools-temp\*\*.dll" -Destination $targetDir -Force
    Copy-Item -Path "kiwix-tools-temp\*\kiwix-serve.exe" -Destination $targetDir -Force

    Write-Host "`nSuccess! The Kiwix engine is now fully installed and portable."
} catch {
    Write-Host "`nError downloading or installing Kiwix tools: $_" -ForegroundColor Red
    Write-Host "You can manually download the Kiwix Tools for Windows from https://download.kiwix.org/release/kiwix-tools/"
    Write-Host "and extract the DLLs into sidecars/kiwix-serve/windows/"
} finally {
    if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
    if (Test-Path "kiwix-tools-temp") { Remove-Item "kiwix-tools-temp" -Recurse -Force }
}

Write-Host "`nPress any key to exit..."
$Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") | Out-Null
