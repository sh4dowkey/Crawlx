param(
    [switch]$Uninstall,
    [switch]$Help,
    [switch]$Force
)

# CrawlX Installation Script for Windows
# Version: 2.0
# Requires Administrator privileges

#Requires -RunAsAdministrator

$Config = @{
    BinaryName = "crawlx.exe"
    SourcePath = ".\dist\crawlx.exe"
    InstallDir = "$env:LOCALAPPDATA\Programs\crawlx"
    LegacyPath = "C:\Windows\System32\crawlx.exe"
    ScriptVersion = "2.0"
}

$DestPath = "$($Config.InstallDir)\$($Config.BinaryName)"

# Logging functions
function Write-LogInfo { param([string]$Message); Write-Host "[INFO] $Message" -ForegroundColor Blue }
function Write-LogSuccess { param([string]$Message); Write-Host "[SUCCESS] $Message" -ForegroundColor Green }
function Write-LogWarning { param([string]$Message); Write-Host "[WARNING] $Message" -ForegroundColor Yellow }
function Write-LogError { param([string]$Message); Write-Host "[ERROR] $Message" -ForegroundColor Red }
function Write-LogStep { param([string]$Message); Write-Host ">>> $Message" -ForegroundColor Cyan }

function Show-Banner {
    Write-Host "`n================================================================================" -ForegroundColor Cyan
    Write-Host "  CrawlX Installation Script v$($Config.ScriptVersion)" -ForegroundColor Cyan
    Write-Host "  Fast, concurrent web crawler built in Go" -ForegroundColor White
    Write-Host "  https://github.com/sh4dowkey/crawlx" -ForegroundColor White
    Write-Host "================================================================================`n" -ForegroundColor Cyan
}

function Test-SystemRequirements {
    Write-LogStep "Checking system requirements"
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-LogError "PowerShell 5.0 or higher required"
        throw "Incompatible PowerShell version"
    }
    
    # Check if running as Administrator
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) {
        Write-LogError "Administrator privileges required"
        throw "Run PowerShell as Administrator"
    }
    
    Write-LogSuccess "System requirements verified"
}

function Find-Binary {
    Write-LogStep "Detecting binary file"
    
    if (Test-Path $Config.SourcePath) {
        $fileInfo = Get-Item $Config.SourcePath
        $fileSizeMB = [math]::Round($fileInfo.Length / 1MB, 2)
        Write-LogSuccess "Found binary: $($Config.SourcePath) (${fileSizeMB}MB)"
        return $true
    }
    
    Write-LogError "Binary not found at: $($Config.SourcePath)"
    Write-LogError "Expected structure:"
    Write-Host "  folder/" -ForegroundColor Gray
    Write-Host "    ├── dist/" -ForegroundColor Gray
    Write-Host "    │   └── crawlx.exe" -ForegroundColor Gray
    Write-Host "    └── scripts/" -ForegroundColor Gray
    Write-Host "        └── setup.ps1" -ForegroundColor Gray
    throw "Binary not found"
}

function Remove-ExistingInstallations {
    Write-LogStep "Checking for existing installations"
    
    $removed = @()
    
    # Remove current installation
    if (Test-Path $DestPath) {
        Remove-Item $DestPath -Force -ErrorAction SilentlyContinue
        $removed += $DestPath
    }
    
    # Remove legacy installation
    if (Test-Path $Config.LegacyPath) {
        Remove-Item $Config.LegacyPath -Force -ErrorAction SilentlyContinue
        $removed += $Config.LegacyPath
    }
    
    if ($removed.Count -gt 0) {
        Write-LogInfo "Removed existing installations: $($removed.Count)"
    } else {
        Write-LogInfo "No existing installation found"
    }
}

function Install-Binary {
    Write-LogStep "Installing binary"
    
    # Create install directory
    if (-not (Test-Path $Config.InstallDir)) {
        New-Item -Path $Config.InstallDir -ItemType Directory -Force | Out-Null
    }
    
    # Copy binary
    Copy-Item -Path $Config.SourcePath -Destination $DestPath -Force
    Write-LogSuccess "Binary installed to: $DestPath"
}

function Update-SystemPath {
    Write-LogStep "Updating system PATH"
    
    $currentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$($Config.InstallDir)*") {
        $newPath = "$currentPath;$($Config.InstallDir)"
        [System.Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        $env:PATH += ";$($Config.InstallDir)"
        Write-LogSuccess "Added to PATH"
        Write-LogWarning "Restart terminal for PATH changes to take effect"
    } else {
        Write-LogInfo "Already in PATH"
    }
}

function Test-Installation {
    Write-LogStep "Verifying installation"
    
    if (-not (Test-Path $DestPath)) {
        throw "Installation verification failed"
    }
    
    Write-LogSuccess "Installation completed successfully!"
    Write-Host "  Binary location: $DestPath" -ForegroundColor White
    Write-Host "  Installation directory: $($Config.InstallDir)" -ForegroundColor White
}

function Show-UsageInfo {
    Write-Host "`n=== QUICK START ===" -ForegroundColor Green
    Write-Host "Basic usage:" -ForegroundColor White
    Write-Host "  crawlx -u https://example.com" -ForegroundColor Yellow
    Write-Host "Deep crawl:" -ForegroundColor White  
    Write-Host "  crawlx -u https://example.com -d 3 --verbose" -ForegroundColor Yellow
    Write-Host "`n=== NEXT STEPS ===" -ForegroundColor Green
    Write-Host "1. Open a new PowerShell window" -ForegroundColor White
    Write-Host "2. Try: crawlx -u https://example.com -d 1" -ForegroundColor Yellow
    Write-Host "3. Get help: crawlx --help" -ForegroundColor White
    Write-Host "`nTo uninstall: .\scripts\setup.ps1 -Uninstall`n" -ForegroundColor Gray
}

function Uninstall-CrawlX {
    Write-LogStep "Starting uninstallation"
    
    $removed = @()
    
    # Remove main installation
    if (Test-Path $DestPath) {
        Remove-Item $DestPath -Force
        $removed += $DestPath
    }
    
    # Remove installation directory if empty
    if (Test-Path $Config.InstallDir) {
        $contents = Get-ChildItem $Config.InstallDir -ErrorAction SilentlyContinue
        if ($contents.Count -eq 0) {
            Remove-Item $Config.InstallDir -Force
            $removed += $Config.InstallDir
        }
    }
    
    # Remove legacy installation
    if (Test-Path $Config.LegacyPath) {
        Remove-Item $Config.LegacyPath -Force -ErrorAction SilentlyContinue
        $removed += $Config.LegacyPath
    }
    
    # Clean PATH
    $currentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    $pathEntries = $currentPath -split ';' | Where-Object { $_ -ne $Config.InstallDir }
    $newPath = $pathEntries -join ';'
    [System.Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    
    Write-LogSuccess "Uninstallation completed"
    Write-LogInfo "Removed $($removed.Count) item(s)"
}

function Install-CrawlX {
    Show-Banner
    
    try {
        Test-SystemRequirements
        Find-Binary
        Remove-ExistingInstallations
        Install-Binary
        Update-SystemPath
        Test-Installation
        Show-UsageInfo
        
    } catch {
        Write-LogError "Installation failed: $($_.Exception.Message)"
        Write-LogError "Solutions:"
        Write-LogError "1. Run as Administrator"
        Write-LogError "2. Check binary exists in .\dist\"
        Write-LogError "3. Try with -Force flag"
        exit 1
    }
}

function Show-Help {
    Write-Host "CrawlX Installation Script v$($Config.ScriptVersion)" -ForegroundColor Cyan
    Write-Host "`nUSAGE:" -ForegroundColor Yellow
    Write-Host "  .\scripts\setup.ps1                Install CrawlX"
    Write-Host "  .\scripts\setup.ps1 -Uninstall     Remove CrawlX"  
    Write-Host "  .\scripts\setup.ps1 -Help          Show this help"
    Write-Host "`nREQUIREMENTS:" -ForegroundColor Yellow
    Write-Host "  - Administrator privileges"
    Write-Host "  - PowerShell 5.0+"
    Write-Host "  - Windows 10+ (recommended)"
    Write-Host "`nEXAMPLES:" -ForegroundColor Yellow
    Write-Host "  .\scripts\setup.ps1           # Install"
    Write-Host "  .\scripts\setup.ps1 -Force    # Force install"
    Write-Host "  .\scripts\setup.ps1 -Uninstall # Remove"
}

# Main execution
try {
    if ($Help) { Show-Help }
    elseif ($Uninstall) { Show-Banner; Uninstall-CrawlX }
    else { Install-CrawlX }
} catch {
    Write-LogError "Script failed: $($_.Exception.Message)"
    exit 1
}