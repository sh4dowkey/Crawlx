#!/bin/bash

# CrawlX Installation Script for Linux/macOS
# Version: 2.0
# Description: Installs CrawlX web crawler system-wide with proper validation and error handling
# Requirements: sudo privileges for system installation

set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Configuration
readonly BINARY_NAME="crawlx"
readonly INSTALL_DIR="/usr/local/bin"
readonly DEST_PATH="$INSTALL_DIR/$BINARY_NAME"
readonly SCRIPT_VERSION="2.0"
readonly MIN_FREE_SPACE_KB=50000  # 50MB minimum free space

# Color codes for output formatting
readonly COLOR_RED='\033[0;31m'
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_YELLOW='\033[0;33m'
readonly COLOR_BLUE='\033[0;34m'
readonly COLOR_CYAN='\033[0;36m'
readonly COLOR_RESET='\033[0m'

# Logging functions
log_info() {
    echo -e "${COLOR_BLUE}[INFO]${COLOR_RESET} $1"
}

log_success() {
    echo -e "${COLOR_GREEN}[SUCCESS]${COLOR_RESET} $1"
}

log_warning() {
    echo -e "${COLOR_YELLOW}[WARNING]${COLOR_RESET} $1"
}

log_error() {
    echo -e "${COLOR_RED}[ERROR]${COLOR_RESET} $1" >&2
}

log_step() {
    echo -e "${COLOR_CYAN}>>> $1${COLOR_RESET}"
}

# Display banner with version information
display_banner() {
    echo ""
    echo "================================================================================"
    echo "  CrawlX Installation Script v${SCRIPT_VERSION}"
    echo "  Fast, concurrent web crawler built in Go"
    echo "  https://github.com/sh4dowkey/crawlx"
    echo "================================================================================"
    echo ""
}

# System information detection
detect_system() {
    local os_name arch_name
    
    os_name=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch_name=$(uname -m)
    
    case "$arch_name" in
        x86_64) arch_name="amd64" ;;
        aarch64|arm64) arch_name="arm64" ;;
        armv7l) arch_name="arm" ;;
        *) 
            log_warning "Unknown architecture: $arch_name, trying amd64"
            arch_name="amd64"
            ;;
    esac
    
    log_info "Detected system: $os_name-$arch_name"
    echo "$os_name-$arch_name"
}

# Check system requirements
check_requirements() {
    log_step "Checking system requirements"
    
    # Check if running as root (for sudo check)
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root. This is not recommended for security reasons."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Installation cancelled by user"
            exit 0
        fi
    fi
    
    # Check sudo availability
    if ! command -v sudo &> /dev/null; then
        log_error "sudo is required but not installed"
        log_error "Please install sudo or run this script as root"
        exit 1
    fi
    
    # Test sudo access
    if ! sudo -n true 2>/dev/null; then
        log_info "Checking sudo privileges (you may be prompted for your password)"
        if ! sudo -v; then
            log_error "sudo privileges required for system installation"
            exit 1
        fi
    fi
    
    # Check available disk space
    local available_space
    available_space=$(df "$INSTALL_DIR" 2>/dev/null | awk 'NR==2 {print $4}' || echo "0")
    
    if [[ $available_space -lt $MIN_FREE_SPACE_KB ]]; then
        log_error "Insufficient disk space. Need at least 50MB free in $INSTALL_DIR"
        exit 1
    fi
    
    log_success "System requirements verified"
}

# Detect and validate binary
detect_binary() {
    log_step "Detecting binary file"
    
    local system_info binary_candidates
    system_info=$(detect_system)
    
    # List of possible binary names in order of preference
    binary_candidates=(
        "./dist/${BINARY_NAME}-${system_info}"
        "./dist/${BINARY_NAME}"
        "./${BINARY_NAME}"
        "./${BINARY_NAME}-${system_info}"
    )
    
    for candidate in "${binary_candidates[@]}"; do
        if [[ -f "$candidate" ]]; then
            # Verify it's actually executable
            if [[ -x "$candidate" ]] || chmod +x "$candidate" 2>/dev/null; then
                SOURCE_PATH="$candidate"
                log_success "Found binary: $SOURCE_PATH"
                
                # Display binary information
                local file_size
                file_size=$(du -h "$SOURCE_PATH" | cut -f1)
                log_info "Binary size: $file_size"
                
                # Check if it's a valid Go binary
                if file "$SOURCE_PATH" 2>/dev/null | grep -q "executable"; then
                    log_info "Binary validation passed"
                    return 0
                else
                    log_warning "Binary may not be a valid executable"
                fi
                return 0
            else
                log_warning "Found $candidate but cannot make it executable"
            fi
        fi
    done
    
    log_error "No suitable binary found in expected locations:"
    for candidate in "${binary_candidates[@]}"; do
        echo "  - $candidate"
    done
    echo ""
    log_error "Please ensure you have:"
    log_error "1. Downloaded the correct release for your system ($system_info)"
    log_error "2. Extracted the archive completely"
    log_error "3. Run this script from the extracted directory"
    echo ""
    log_info "Download from: https://github.com/sh4dowkey/crawlx/releases/latest"
    exit 1
}

# Check for existing installation
check_existing_installation() {
    log_step "Checking for existing installation"
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        local existing_path current_version
        existing_path=$(command -v "$BINARY_NAME")
        
        log_warning "CrawlX is already installed at: $existing_path"
        
        # Try to get version information
        if current_version=$("$existing_path" --version 2>/dev/null || echo "unknown"); then
            log_info "Currently installed version: $current_version"
        fi
        
        echo ""
        echo "Options:"
        echo "  1) Update/reinstall (recommended)"
        echo "  2) Cancel installation"
        echo ""
        read -p "Choose option (1/2): " -n 1 -r
        echo ""
        
        case $REPLY in
            1) log_info "Proceeding with update/reinstall" ;;
            2) log_info "Installation cancelled by user"; exit 0 ;;
            *) log_info "Invalid option, cancelling installation"; exit 0 ;;
        esac
    else
        log_info "No existing installation found"
    fi
}

# Backup existing installation
backup_existing() {
    if [[ -f "$DEST_PATH" ]]; then
        log_step "Backing up existing installation"
        local backup_path="${DEST_PATH}.backup.$(date +%Y%m%d_%H%M%S)"
        
        if sudo cp "$DEST_PATH" "$backup_path"; then
            log_success "Backup created: $backup_path"
        else
            log_warning "Could not create backup, continuing anyway"
        fi
    fi
}

# Install binary
install_binary() {
    log_step "Installing binary to system directory"
    
    # Create install directory if it doesn't exist
    if [[ ! -d "$INSTALL_DIR" ]]; then
        log_info "Creating install directory: $INSTALL_DIR"
        if ! sudo mkdir -p "$INSTALL_DIR"; then
            log_error "Failed to create install directory"
            exit 1
        fi
    fi
    
    # Copy binary with verification
    log_info "Copying $SOURCE_PATH to $DEST_PATH"
    if ! sudo cp "$SOURCE_PATH" "$DEST_PATH"; then
        log_error "Failed to copy binary to install directory"
        log_error "Check permissions and disk space"
        exit 1
    fi
    
    # Set proper permissions
    log_info "Setting executable permissions"
    if ! sudo chmod 755 "$DEST_PATH"; then
        log_error "Failed to set executable permissions"
        exit 1
    fi
    
    # Verify ownership
    sudo chown root:root "$DEST_PATH" 2>/dev/null || {
        log_warning "Could not set root ownership (this may be normal on some systems)"
    }
    
    log_success "Binary installed successfully"
}

# Verify installation
verify_installation() {
    log_step "Verifying installation"
    
    # Check if binary exists and is executable
    if [[ ! -f "$DEST_PATH" ]]; then
        log_error "Binary not found at $DEST_PATH"
        exit 1
    fi
    
    if [[ ! -x "$DEST_PATH" ]]; then
        log_error "Binary is not executable"
        exit 1
    fi
    
    # Check if binary is in PATH
    if ! command -v "$BINARY_NAME" &> /dev/null; then
        log_error "Binary not found in PATH"
        log_error "You may need to restart your terminal or add $INSTALL_DIR to your PATH"
        exit 1
    fi
    
    # Test binary execution
    log_info "Testing binary execution"
    if timeout 5s "$BINARY_NAME" --help &> /dev/null; then
        log_success "Binary execution test passed"
    else
        log_warning "Binary execution test failed (may be normal if --help not implemented)"
    fi
    
    # Display installation details
    local binary_path file_size
    binary_path=$(command -v "$BINARY_NAME")
    file_size=$(du -h "$binary_path" | cut -f1)
    
    echo ""
    log_success "Installation completed successfully!"
    echo "  Binary location: $binary_path"
    echo "  Binary size: $file_size"
    echo "  Installation directory: $INSTALL_DIR"
    echo ""
}

# Display usage information
display_usage_info() {
    echo "=== QUICK START ==="
    echo ""
    echo "Basic usage:"
    echo "  $BINARY_NAME -u https://example.com"
    echo ""
    echo "Deep crawl with details:"
    echo "  $BINARY_NAME -u https://example.com -d 3 --verbose"
    echo ""
    echo "Get help:"
    echo "  $BINARY_NAME --help"
    echo ""
    echo "=== NEXT STEPS ==="
    echo ""
    echo "1. Open a new terminal window (to refresh PATH)"
    echo "2. Try: $BINARY_NAME -u https://example.com -d 1"
    echo "3. Read the documentation: https://github.com/sh4dowkey/crawlx"
    echo ""
    echo "To uninstall: sudo rm $DEST_PATH"
    echo ""
}

# Uninstallation function
uninstall() {
    log_step "Starting uninstallation process"
    
    if [[ ! -f "$DEST_PATH" ]]; then
        log_warning "$BINARY_NAME is not installed at $DEST_PATH"
        return 0
    fi
    
    log_info "Removing $DEST_PATH"
    if sudo rm -f "$DEST_PATH"; then
        log_success "CrawlX uninstalled successfully"
    else
        log_error "Failed to remove binary"
        exit 1
    fi
    
    # Clean up any backup files
    local backup_count
    backup_count=$(sudo find "$INSTALL_DIR" -name "${BINARY_NAME}.backup.*" 2>/dev/null | wc -l)
    if [[ $backup_count -gt 0 ]]; then
        echo ""
        read -p "Remove $backup_count backup file(s)? (y/N): " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo find "$INSTALL_DIR" -name "${BINARY_NAME}.backup.*" -delete 2>/dev/null
            log_info "Backup files removed"
        fi
    fi
    
    echo ""
    log_success "Uninstallation completed"
}

# Main installation function
main() {
    display_banner
    
    check_requirements
    detect_binary
    check_existing_installation
    backup_existing
    install_binary
    verify_installation
    display_usage_info
    
    echo "================================================================================"
    echo "  Installation completed successfully!"
    echo "  Thank you for using CrawlX!"
    echo "================================================================================"
}

# Script entry point
if [[ "${1:-}" == "uninstall" ]]; then
    display_banner
    uninstall
elif [[ "${1:-}" == "--help" ]] || [[ "${1:-}" == "-h" ]]; then
    echo "CrawlX Installation Script"
    echo ""
    echo "Usage:"
    echo "  $0              Install CrawlX"
    echo "  $0 uninstall    Uninstall CrawlX"
    echo "  $0 --help       Show this help"
    echo ""
    echo "Requirements:"
    echo "  - sudo privileges"
    echo "  - 50MB free disk space"
    echo "  - Linux or macOS system"
    echo ""
else
    main
fi