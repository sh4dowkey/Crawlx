#!/bin/bash

# Configuration
BINARY_NAME="crawlx"
INSTALL_DIR="/usr/local/bin"
DEST_PATH="$INSTALL_DIR/$BINARY_NAME"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
RESET='\033[0m'

# Logging functions
log_info()    { echo -e "${BLUE}[INFO]${RESET} $1"; }
log_success() { echo -e "${GREEN}[OK]${RESET} $1"; }
log_warning() { echo -e "${YELLOW}[WARN]${RESET} $1"; }
log_error()   { echo -e "${RED}[ERROR]${RESET} $1"; }
log_step()    { echo -e "${CYAN}•${RESET} $1"; }

show_header() {
    echo ""
    echo "CrawlX Installation Script"
    echo -e "${GRAY}==========================${RESET}"
    echo ""
}

check_requirements() {
    log_step "Checking system requirements"
    
    # Check if running as root for sudo verification
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root (not recommended)"
    fi
    
    # Check sudo availability
    if ! command -v sudo &> /dev/null; then
        log_error "sudo is required but not installed"
        exit 1
    fi
    
    # Test sudo access
    if ! sudo -n true 2>/dev/null; then
        log_info "Testing sudo privileges (password may be required)"
        if ! sudo -v; then
            log_error "sudo privileges required"
            exit 1
        fi
    fi
    
    log_success "System requirements met"
}

find_binary() {
    log_step "Locating binary file"
    
    # Detect system
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
    esac
    
    # Look for binary
    candidates=(
        "./dist/${BINARY_NAME}-${OS}-${ARCH}"
        "./dist/${BINARY_NAME}"
    )
    
    for candidate in "${candidates[@]}"; do
        if [[ -f "$candidate" ]]; then
            SOURCE_PATH="$candidate"
            size=$(du -h "$candidate" | cut -f1)
            log_success "Found binary: $(basename "$candidate") (${size})"
            return 0
        fi
    done
    
    log_error "Binary not found in ./dist/"
    echo ""
    echo -e "${GRAY}Expected folder structure:${RESET}"
    echo -e "${GRAY}  ├── dist/${RESET}"
    echo -e "${GRAY}  │   └── crawlx${RESET}"
    echo -e "${GRAY}  └── scripts/${RESET}"
    echo -e "${GRAY}      └── setup.sh${RESET}"
    exit 1
}

check_existing() {
    log_step "Checking for existing installations"
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        existing_path=$(command -v "$BINARY_NAME")
        log_success "Found existing installation: $existing_path"
    else
        log_success "No existing installation found"
    fi
}

install_binary() {
    log_step "Installing to system directory"
    
    if ! sudo cp "$SOURCE_PATH" "$DEST_PATH"; then
        log_error "Failed to copy binary"
        exit 1
    fi
    
    if ! sudo chmod +x "$DEST_PATH"; then
        log_error "Failed to set permissions"
        exit 1
    fi
    
    log_success "Binary installed to $INSTALL_DIR"
}

verify_installation() {
    log_step "Verifying installation"
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        log_success "Installation completed successfully"
        echo ""
        echo "Installation Details:"
        echo -e "${GRAY}  Location: $DEST_PATH${RESET}"
        echo -e "${GRAY}  Permissions: $(ls -l "$DEST_PATH" | cut -d' ' -f1)${RESET}"
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

show_next_steps() {
    echo ""
    echo "Next Steps:"
    echo -e "${GRAY}1. Open a new terminal window${RESET}"
    echo -e "${GRAY}2. Test: crawlx -u https://example.com -d 1${RESET}"
    echo -e "${GRAY}3. Get help: crawlx --help${RESET}"
    echo ""
    echo -e "${YELLOW}To uninstall: sudo ./scripts/setup.sh uninstall${RESET}"
}

install_crawlx() {
    show_header
    
    check_requirements
    find_binary
    check_existing
    install_binary
    verify_installation
    show_next_steps
}

uninstall_crawlx() {
    echo ""
    echo "CrawlX Uninstall"
    echo -e "${GRAY}================${RESET}"
    echo ""
    
    log_step "Removing installation"
    
    if [[ -f "$DEST_PATH" ]]; then
        if sudo rm "$DEST_PATH"; then
            log_success "Binary removed from $INSTALL_DIR"
        else
            log_error "Failed to remove binary"
            exit 1
        fi
    else
        log_warning "Binary not found at $DEST_PATH"
    fi
    
    log_success "Uninstallation completed"
    echo ""
    echo -e "${YELLOW}Terminal restart recommended${RESET}"
}

show_help() {
    echo ""
    echo "CrawlX Setup Script"
    echo -e "${GRAY}===================${RESET}"
    echo ""
    echo -e "${YELLOW}USAGE:${RESET}"
    echo "  sudo ./setup.sh           Install CrawlX"
    echo "  sudo ./setup.sh uninstall Remove CrawlX"
    echo "  ./setup.sh help           Show this help"
    echo ""
    echo -e "${YELLOW}REQUIREMENTS:${RESET}"
    echo "  • sudo privileges"
    echo "  • Linux or macOS"
    echo "  • /usr/local/bin writable"
    echo ""
}

# Main execution
case "${1:-install}" in
    "uninstall")
        uninstall_crawlx
        ;;
    "help"|"--help"|"-h")
        show_help
        ;;
    *)
        install_crawlx
        ;;
esac