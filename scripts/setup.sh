#!/bin/bash

# --- Configuration ---
BINARY_NAME="crawlx"
INSTALL_DIR="/usr/local/bin"
DEST_PATH="$INSTALL_DIR/$BINARY_NAME"

# --- ANSI Color Codes ---
COLOR_CYAN='\033[0;36m'
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[0;33m'
COLOR_RESET='\033[0m'

# --- Banners ---
get_banner() {
    # (Your awesome ASCII art banner remains the same)
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

# --- Helper Functions ---
detect_binary() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
    esac

    # Construct the expected binary name (e.g., crawlx-linux-amd64)
    # This matches the output of our cross-compilation commands.
    local detected_binary_name="${BINARY_NAME}-${OS}-${ARCH}"
    if [ -f "./dist/$detected_binary_name" ]; then
        SOURCE_PATH="./dist/$detected_binary_name"
    elif [ -f "./dist/$BINARY_NAME" ]; then
        # Fallback for a generic binary name
        SOURCE_PATH="./dist/$BINARY_NAME"
    else
        echo ""
        return 1
    fi
    return 0
}

# --- Main Logic ---
main() {
    echo -e "$(get_banner)"
    echo -e "${COLOR_CYAN}         CRAWLX CLI INSTALLER${COLOR_RESET}"
    echo ""

    # --- Step 1: Check for binary ---
    echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 1 of 4: Detecting binary..."
    if ! detect_binary; then
        echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Could not find a suitable binary in './dist/'."
        echo "       The installation cannot proceed."
        # (Your helpful build guide remains the same)
        exit 1
    fi
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Binary found at '$SOURCE_PATH'."
    echo ""

    # --- Step 2: Copy binary ---
    echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 2 of 4: Copying file to system directory..."
    sudo cp "$SOURCE_PATH" "$DEST_PATH"
    if [ $? -ne 0 ]; then
        echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Installation failed. Check permissions or try running with 'sudo'."
        exit 1
    fi
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} File copied successfully."
    echo ""

    # --- Step 3: Grant executable permissions ---
    echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 3 of 4: Granting executable permissions..."
    sudo chmod +x "$DEST_PATH"
    if [ $? -ne 0 ]; then
        echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Failed to set permissions. Try running with 'sudo'."
        exit 1
    fi
    echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Permissions set successfully."
    echo ""

    # --- Step 4: Verify installation ---
    echo -e " ${COLOR_CYAN}>>${COLOR_RESET} Step 4 of 4: Verifying final installation..."
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Installation complete! '$BINARY_NAME' is ready."
        echo ""
        echo "You can now run '$BINARY_NAME' from any folder."
        echo -e "To uninstall, run: ${COLOR_YELLOW}bash setup.sh uninstall${COLOR_RESET}"
    else
        echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Verification failed. Please check your system's PATH."
        exit 1
    fi
}

uninstall() {
    echo -e "${COLOR_YELLOW}Starting uninstallation process...${COLOR_RESET}"
    if [ -f "$DEST_PATH" ]; then
        echo "Removing '$DEST_PATH'..."
        sudo rm "$DEST_PATH"
        if [ $? -eq 0 ]; then
            echo -e "${COLOR_GREEN}[OK]${COLOR_RESET} Uninstallation successful."
        else
            echo -e "${COLOR_RED}[FAIL]${COLOR_RESET} Could not remove the file. Please check permissions."
        fi
    else
        echo -e "${COLOR_YELLOW}[INFO]${COLOR_RESET} '$BINARY_NAME' is not installed in '$INSTALL_DIR'."
    fi
}

# --- Script Entry Point ---
if [ "$1" == "uninstall" ]; then
    uninstall
else
    main
fi