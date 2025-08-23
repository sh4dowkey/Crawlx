# Requires administrator privileges to run
# To run, right-click the file and select "Run with PowerShell"
# Or run from an elevated PowerShell terminal: .\install.ps1

$binaryName = "crawlx.exe"
$sourcePath = ".\dist\$binaryName"
$destPath = "C:\Windows\System32\$binaryName"

Write-Host "🔧 Installing $binaryName to $destPath..."

# Check if the pre-built binary exists
if (-not (Test-Path $sourcePath)) {
    Write-Host "❌ Error: Pre-built binary not found at $sourcePath. Please build the project first." -ForegroundColor Red
    exit
}

# Copy the binary
Write-Host "📂 Copying pre-built binary..."
try {
    Copy-Item -Path $sourcePath -Destination $destPath -Force
} catch {
    Write-Host "❌ Installation failed: $_" -ForegroundColor Red
    Write-Host "Please run this script as an administrator." -ForegroundColor Red
    exit
}

Write-Host "✅ Installation complete. You can now run '$binaryName' from any folder!" -ForegroundColor Green