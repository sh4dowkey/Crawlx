param([switch]$Uninstall)

# Requires administrator privileges to run for PATH modification.
# To run, right-click the file and select "Run with PowerShell as Administrator"

# --- Configuration ---
$binaryName = "crawlx.exe"
$sourcePath = ".\dist\$binaryName"
$installDir = "$env:LOCALAPPDATA\Programs\crawlx" # Recommended install location
$destPath = "$installDir\$binaryName"
$oldSystem32Path = "C:\Windows\System32\$binaryName" # Path for the old, incorrect location

# --- Banner Function ---
Function Get-CrawlxBanner {
    # (Your awesome ASCII art banner remains the same)
    $banner = @"
======================================
                                                                        
 @@@@@@@  @@@@@@@    @@@@@@   @@@  @@@  @@@  @@@          @@@  @@@  
@@@@@@@@  @@@@@@@@  @@@@@@@@  @@@  @@@  @@@  @@@          @@@  @@@  
!@@       @@!  @@@  @@!  @@@  @@!  @@!  @@!  @@!          @@!  !@@  
!@!       !@!  @!@  !@!  @!@  !@!  !@!  !@!  !@!          !@!  @!!  
!@!       @!@!!@!   @!@!@!@!  @!!  !!@  @!@  @!!           !@@!@!   
!!!       !!@!@!    !!!@!!!!  !@!  !!!  !@!  !!!            @!!!    
:!!       !!: :!!   !!:  !!!  !!:  !!:  !!:  !!:           !: :!!   
:!:       :!:  !:!  :!:  !:!  :!:  :!:  :!:   :!:         :!:  !:!  
 ::: :::  ::   :::  ::   :::   :::: :: :::    :: ::::      ::  :::  
 :: :: :   :   : :   :   : :    :: :  : :    : :: : :      :   ::   
 
======================================
"@
    return $banner
}

# --- Main Installation Logic ---
Function Install-Crawlx {
    Write-Host (Get-CrawlxBanner) -ForegroundColor Cyan
    Write-Host "         CRAWLX CLI INSTALLER" -ForegroundColor Cyan
    Write-Host ""

    # --- Step 1: Clean up old versions ---
    Write-Host ">> Step 1 of 5: Checking for legacy installations..." -ForegroundColor Yellow
    if (Test-Path $oldSystem32Path) {
        Write-Host "[INFO] Found old version in System32. Removing it..." -ForegroundColor Cyan
        try {
            Remove-Item -Path $oldSystem32Path -Force -ErrorAction Stop
            Write-Host "[OK] Legacy version removed successfully." -ForegroundColor Green
        } catch {
            Write-Host "[FAIL] Could not remove legacy version at '$oldSystem32Path'. Please remove it manually." -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "[OK] No legacy installation found." -ForegroundColor Green
    }
    Write-Host ""

    # --- Step 2: Check for binary ---
    Write-Host ">> Step 2 of 5: Verifying binary location..." -ForegroundColor Yellow
    if (-not (Test-Path $sourcePath)) {
        Write-Host "[FAIL] Binary not found at '$sourcePath'." -ForegroundColor Red
        exit 1
    }
    Write-Host "[OK] Binary found." -ForegroundColor Green
    Write-Host ""

    # --- Step 3: Create installation directory ---
    Write-Host ">> Step 3 of 5: Creating installation directory..." -ForegroundColor Yellow
    if (-not (Test-Path $installDir)) {
        New-Item -Path $installDir -ItemType Directory | Out-Null
    }
    Write-Host "[OK] Directory '$installDir' is ready." -ForegroundColor Green
    Write-Host ""

    # --- Step 4: Copy binary ---
    Write-Host ">> Step 4 of 5: Copying file..." -ForegroundColor Yellow
    try {
        Copy-Item -Path $sourcePath -Destination $destPath -Force -ErrorAction Stop
        Write-Host "[OK] File copied successfully." -ForegroundColor Green
    } catch {
        Write-Host "[FAIL] Installation failed: $_" -ForegroundColor Red
        exit 1
    }
    Write-Host ""

    # --- Step 5: Add to User PATH ---
    Write-Host ">> Step 5 of 5: Adding to your PATH..." -ForegroundColor Yellow
    try {
        $currentUserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
        if (-not ($currentUserPath -split ';' -contains $installDir)) {
            $newUserPath = $currentUserPath + ";" + $installDir
            [System.Environment]::SetEnvironmentVariable("PATH", $newUserPath, "User")
            Write-Host "[OK] Installation path added to your User PATH." -ForegroundColor Green
            Write-Host "       Please restart your terminal for the changes to take effect." -ForegroundColor Yellow
        } else {
            Write-Host "[INFO] Path is already configured." -ForegroundColor Cyan
        }
    } catch {
        Write-Host "[FAIL] Could not modify PATH: $_" -ForegroundColor Red
        Write-Host "       Tip: Ensure you are running this script as an Administrator." -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host ""
    Write-Host "Installation Complete! You can now run 'crawlx' from a new terminal." -ForegroundColor Cyan
    Write-Host "To uninstall, run: powershell.exe -File setup.ps1 -Uninstall" -ForegroundColor Yellow
}

# --- Uninstallation Logic ---
Function Uninstall-Crawlx {
    # (Uninstallation logic remains the same)
    Write-Host "Starting uninstallation of crawlx..." -ForegroundColor Yellow
    
    $currentUserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentUserPath -split ';' -contains $installDir) {
        $newUserPath = ($currentUserPath -split ';') | Where-Object { $_ -ne $installDir } | ForEach-Object { "$_" }
        $newUserPath = $newUserPath -join ';'
        [System.Environment]::SetEnvironmentVariable("PATH", $newUserPath, "User")
        Write-Host "[OK] Removed from User PATH." -ForegroundColor Green
    }

    if (Test-Path $installDir) {
        Remove-Item -Path $installDir -Recurse -Force
        Write-Host "[OK] Removed installation directory." -ForegroundColor Green
    }

    if (Test-Path $oldSystem32Path) {
        Remove-Item -Path $oldSystem32Path -Force -ErrorAction SilentlyContinue
        Write-Host "[OK] Removed legacy System32 file." -ForegroundColor Green
    }

    Write-Host "Uninstallation complete. Please restart your terminal." -ForegroundColor Cyan
}

# --- Script Entry Point ---
if ($Uninstall) {
    Uninstall-Crawlx
} else {
    Install-Crawlx
}