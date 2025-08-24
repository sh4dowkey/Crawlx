# Requires administrator privileges to run.
# To run, right-click the file and select "Run with PowerShell"
# Or run from an elevated PowerShell terminal: .\install.ps1

$binaryName = "crawlx.exe"
$sourcePath = ".\dist\$binaryName"
$destPath = "C:\Windows\System32\$binaryName"

# --- Main Installation Script ---

Function Get-CrawlxBanner {
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

Write-Host ""
Write-Host "     --- Downloading ---" -ForegroundColor Cyan
Write-Host ""
Write-Host "        Please wait..." -ForegroundColor Yellow
Write-Host ""

# Simulate a short delay for effect
Start-Sleep -Seconds 1

Write-Host (Get-CrawlxBanner) -ForegroundColor Cyan
Write-Host "            CRAWLX CLI INSTALLER" -ForegroundColor Cyan
Write-Host ""
Write-Host "Starting installation of '$binaryName'..."
Write-Host ""

# --- Step 1: Check for binary ---
Write-Host ">> Step 1 of 3: Verifying binary location..." -ForegroundColor Yellow
if (-not (Test-Path $sourcePath)) {
    Write-Host "[FAIL] Binary not found at '$sourcePath'." -ForegroundColor Red
    Write-Host "       The installation cannot proceed." -ForegroundColor Red
    Write-Host ""
    
    # --- Guidance Section ---
    Write-Host "--- Guide: How to Fix This ---" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "It seems the project has not been built yet. To build the executable file, you need to follow these steps:"
    Write-Host ""
    
    Write-Host "Step 1: Install Go" -ForegroundColor Green
    Write-Host "        - Go to https://go.dev/dl/ to download and install the latest version of Go."
    Write-Host "        - Follow the installer instructions."
    Write-Host ""
    
    Write-Host "Step 2: Build the Project" -ForegroundColor Green
    Write-Host "        - Open your terminal (e.g., PowerShell) and navigate to the project's root folder."
    Write-Host "        - Run the command: " -NoNewline
    Write-Host "go build -o ./dist/crawlx.exe ./cmd/" -ForegroundColor Yellow
    Write-Host ""
    
    Write-Host "Step 3: Run the Installer Again" -ForegroundColor Green
    Write-Host "        - After the build is complete, run this script one more time to finish the installation."
    Write-Host ""
    Write-Host "-------------------------------" -ForegroundColor Cyan
    exit
}
Write-Host "[OK] Binary found." -ForegroundColor Green
Write-Host ""

# --- Step 2: Copy binary ---
Write-Host ">> Step 2 of 3: Copying file to system directory..." -ForegroundColor Yellow
try {
    Copy-Item -Path $sourcePath -Destination $destPath -Force -ErrorAction Stop
    Write-Host "[OK] File copied successfully." -ForegroundColor Green
} catch {
    Write-Host "[FAIL] Installation failed: $_" -ForegroundColor Red
    Write-Host "       Tip: Run this script as an administrator to gain the required permissions." -ForegroundColor Yellow
    exit
}
Write-Host ""

# --- Step 3: Verify installation ---
Write-Host ">> Step 3 of 3: Verifying final installation..." -ForegroundColor Yellow
if (Test-Path $destPath) {
    Write-Host "[OK] Installation complete! 'crawlx.exe' is ready." -ForegroundColor Green
} else {
    Write-Host "[FAIL] Verification failed. File not found at destination." -ForegroundColor Red
}
Write-Host ""
Write-Host "You can now run 'crawlx.exe' from any folder." -ForegroundColor Cyan


