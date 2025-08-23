#!/bin/bash

# Define the binary name and destination path
BINARY_NAME="crawlx"
SOURCE_PATH="./dist/$BINARY_NAME"
DEST_PATH="/usr/local/bin/$BINARY_NAME"

echo "🔧 Installing $BINARY_NAME to $DEST_PATH..."

# Check if the binary exists
if [ ! -f "$SOURCE_PATH" ]; then
    echo "❌ Error: Binary not found at $SOURCE_PATH. Please build the project first."
    exit 1
fi

# Copy the binary with sudo for permissions
sudo cp "$SOURCE_PATH" "$DEST_PATH"

if [ $? -eq 0 ]; then
    echo "✅ Installation complete. You can now run '$BINARY_NAME' from any folder!"
else
    echo "❌ Installation failed. Check permissions or try running with 'sudo'."
    exit 1
fi