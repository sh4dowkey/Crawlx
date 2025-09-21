param([switch]$Uninstall, [switch]$Help)

#Requires -RunAsAdministrator

# Configuration
$Config = @{
    BinaryName = "crawlx.exe"
    SourcePath = ".\dist\crawlx.exe"
    InstallDir = "$env:LOCALAPPDATA\Programs\crawlx"
    LegacyPath = "C:\Windows\System32\crawlx.exe"
    Version = "2.0"
}

$DestPath = "$($Config.InstallDir)\$($Config.BinaryName)"

# Simple logging functions
function Write-Info    { param([string]$Msg) Write-Host "[INFO] $Msg" -ForegroundColor Blue }
function Write-Success { param([string]$Msg) Write-Host "[OK] $Msg" -ForegroundColor Green }
function Write-Warning { param([string]$Msg) Write-Host "[WARN] $Msg" -ForegroundColor Yellow }
function Write-Error   { param([string]$Msg) Write-Host "[ERROR] $Msg" -ForegroundColor Red }
function Write-Step    { param([string]$Msg) Write-Host "- $Msg" -ForegroundColor Cyan }

function Show-Header {
    Write-Host ""
    Write-Host "CrawlX Installation Script" -ForegroundColor White
    Write-Host "=========================" -ForegroundColor Gray
    Write-Host ""
}

function Test-Requirements {
    Write-Step "Checking system requirements"
    
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-Error "PowerShell 5.0+ required"
        exit 1
    }
    
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) {
        Write-Error "Administrator privileges required"
        Write-Host "Please run PowerShell as Administrator" -ForegroundColor Gray
        exit 1
    }
    
    Write-Success "System requirements met"
}

function Find-Binary {
    Write-Step "Locating binary file"
    
    if (Test-Path $Config.SourcePath) {
        $size = [math]::Round((Get-Item $Config.SourcePath).Length / 1MB, 1)
        Write-Success "Found crawlx.exe (${size}MB)"
    } else {
        Write-Error "Binary not found: $($Config.SourcePath)"
        Write-Host ""
        Write-Host "Expected folder structure:" -ForegroundColor Gray
        Write-Host "  ├── dist/" -ForegroundColor Gray
        Write-Host "  │   └── crawlx.exe" -ForegroundColor Gray
        Write-Host "  └── scripts/" -ForegroundColor Gray
        Write-Host "      └── setup.ps1" -ForegroundColor Gray
        exit 1
    }
}

function Remove-OldInstallations {
    Write-Step "Checking for existing installations"
    
    $removed = 0
    if (Test-Path $DestPath) {
        Remove-Item $DestPath -Force -ErrorAction SilentlyContinue
        $removed++
    }
    
    if (Test-Path $Config.LegacyPath) {
        Remove-Item $Config.LegacyPath -Force -ErrorAction SilentlyContinue
        $removed++
    }
    
    if ($removed -gt 0) {
        Write-Success "Removed $removed existing installation(s)"
    } else {
        Write-Success "No existing installations found"
    }
}

function Install-Binary {
    Write-Step "Installing to system directory"
    
    if (-not (Test-Path $Config.InstallDir)) {
        New-Item -Path $Config.InstallDir -ItemType Directory -Force | Out-Null
    }
    
    Copy-Item -Path $Config.SourcePath -Destination $DestPath -Force
    Write-Success "Binary installed to $($Config.InstallDir)"
}

function Update-PATH {
    Write-Step "Updating system PATH"
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$($Config.InstallDir)*") {
        $newPath = "$currentPath;$($Config.InstallDir)"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Added to PATH (restart terminal to apply)"
    } else {
        Write-Success "Already in PATH"
    }
}

function Test-Installation {
    Write-Step "Verifying installation"
    
    if (Test-Path $DestPath) {
        Write-Success "Installation completed successfully"
        Write-Host ""
        Write-Host "Installation Details:" -ForegroundColor White
        Write-Host "  Location: $DestPath" -ForegroundColor Gray
        Write-Host "  Added to PATH: Yes" -ForegroundColor Gray
    } else {
        Write-Error "Installation verification failed"
        exit 1
    }
}

function Show-NextSteps {
    Write-Host ""
    Write-Host "Next Steps:" -ForegroundColor White
    Write-Host "1. Open a new PowerShell window" -ForegroundColor Gray
    Write-Host "2. Test: crawlx -u https://example.com -d 1" -ForegroundColor Gray
    Write-Host "3. Get help: crawlx --help" -ForegroundColor Gray
    Write-Host ""
    Write-Host "To uninstall: .\scripts\setup.ps1 -Uninstall" -ForegroundColor Yellow
}

function Install-CrawlX {
    Show-Header
    
    try {
        Test-Requirements
        Find-Binary
        Remove-OldInstallations
        Install-Binary
        Update-PATH
        Test-Installation
        Show-NextSteps
    } catch {
        Write-Error "Installation failed: $($_.Exception.Message)"
        Write-Host ""
        Write-Host "Troubleshooting:" -ForegroundColor Yellow
        Write-Host "- Ensure you're running as Administrator" -ForegroundColor Gray
        Write-Host "- Check that crawlx.exe exists in .\dist\" -ForegroundColor Gray
        Write-Host "- Try restarting PowerShell" -ForegroundColor Gray
        exit 1
    }
}

function Uninstall-CrawlX {
    Write-Host ""
    Write-Host "CrawlX Uninstall" -ForegroundColor White
    Write-Host "================" -ForegroundColor Gray
    Write-Host ""
    
    Write-Step "Removing installation"
    
    $removed = @()
    if (Test-Path $DestPath) {
        Remove-Item $DestPath -Force
        $removed += "Binary"
    }
    
    if (Test-Path $Config.InstallDir) {
        $contents = Get-ChildItem $Config.InstallDir -ErrorAction SilentlyContinue
        if ($contents.Count -eq 0) {
            Remove-Item $Config.InstallDir -Force
            $removed += "Directory"
        }
    }
    
    if (Test-Path $Config.LegacyPath) {
        Remove-Item $Config.LegacyPath -Force -ErrorAction SilentlyContinue
        $removed += "Legacy files"
    }
    
    # Clean PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    $pathEntries = $currentPath -split ';' | Where-Object { $_ -ne $Config.InstallDir }
    $newPath = $pathEntries -join ';'
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    
    Write-Success "Uninstallation completed"
    Write-Host "Removed: $($removed -join ', ')" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Please restart your terminal" -ForegroundColor Yellow
}

function Show-Help {
    Write-Host ""
    Write-Host "CrawlX Setup Script" -ForegroundColor White
    Write-Host "===================" -ForegroundColor Gray
    Write-Host ""
    Write-Host "USAGE:" -ForegroundColor Yellow
    Write-Host "  .\setup.ps1           Install CrawlX"
    Write-Host "  .\setup.ps1 -Uninstall    Remove CrawlX"
    Write-Host "  .\setup.ps1 -Help         Show this help"
    Write-Host ""
    Write-Host "REQUIREMENTS:" -ForegroundColor Yellow
    Write-Host "  - Administrator privileges"
    Write-Host "  - PowerShell 5.0+"
    Write-Host "  - Windows 10+ (recommended)"
    Write-Host ""
}

# Main execution
if ($Help) { 
    Show-Help 
} elseif ($Uninstall) { 
    Uninstall-CrawlX 
} else { 
    Install-CrawlX 
}