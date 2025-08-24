#!/bin/bash

# Define the binary name and destination path
BINARY_NAME="crawlx"
SOURCE_PATH="./dist/$BINARY_NAME"
DEST_PATH="/usr/local/bin/$BINARY_NAME"

# ANSI Color Codes
COLOR_CYAN='\033[0;36m'
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[0;33m'
COLOR_RESET='\033[0m'

# --- Banners ---
get_banner() {
  cat << EOF
${COLOR_CYAN}
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
${COLOR_RESET}
EOF
}

# --- Main Installation Script ---
echo -e "$(get_banner)"
echo -e "${COLOR_CYAN}            CRAWLX CLI INSTALLER${COLOR_RESET}"
echo ""

# --- Step 1: Check for binary ---
echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 1 of 4: Verifying binary location..."
if [ ! -f "$SOURCE_PATH" ]; then
    echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Binary not found at '$SOURCE_PATH'."
    echo "       The installation cannot proceed."
    echo ""
    echo -e "${COLOR_YELLOW}--- Guide: How to Fix This ---${COLOR_RESET}"
    echo "It seems the project has not been built yet. To build the executable:"
    echo "1. Ensure you have Go installed (https://go.dev/dl/)"
    echo "2. Navigate to the project's root directory."
    echo "3. Run the command:"
    echo -e "   ${COLOR_GREEN}go build -o ./dist/crawlx ./cmd/${COLOR_RESET}"
    echo "4. After the build is complete, run this script again."
    echo "-------------------------------"
    exit 1
fi
echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Binary found."
echo ""

# --- Step 2: Copy binary ---
echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 2 of 4: Copying file to system directory..."
sudo cp "$SOURCE_PATH" "$DEST_PATH"
if [ $? -eq 0 ]; then
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} File copied successfully."
else
    echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Installation failed. Check permissions or try running with 'sudo'."
    exit 1
fi
echo ""

# --- Step 3: Grant executable permissions ---
echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 3 of 4: Granting executable permissions..."
sudo chmod +x "$DEST_PATH"
if [ $? -eq 0 ]; then
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Permissions set successfully."
else
    echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Failed to set permissions. Try running with 'sudo'."
    exit 1
fi
echo ""

# --- Step 4: Verify installation ---
echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 4 of 4: Verifying final installation..."
if [ -f "$DEST_PATH" ] && [ -x "$DEST_PATH" ]; then
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Installation complete! '$BINARY_NAME' is ready."
    echo ""
    echo "You can now run '$BINARY_NAME' from any folder."
else
    echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Verification failed. File is not executable or not found."
    exit 1
fi