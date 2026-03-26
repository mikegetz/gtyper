#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

REPO="mikegetz/gtyper"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="gtyper"

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)   OS="linux" ;;
    Darwin)  OS="darwin" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect Architecture
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64)  ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Fetch latest release version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$LATEST_VERSION" ]]; then
    echo "Failed to fetch the latest release version."
    exit 1
fi

BINARY_NAME="${BIN_NAME}-${LATEST_VERSION}-${OS}-${ARCH}"
BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/$BINARY_NAME"

echo "Downloading $BIN_NAME $LATEST_VERSION for $OS-$ARCH from $BINARY_URL..."
if command -v curl >/dev/null 2>&1; then
    curl -L -o "$BIN_NAME" "$BINARY_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$BIN_NAME" "$BINARY_URL"
else
    echo "Error: curl or wget is required to download files."
    exit 1
fi

# Make the binary executable
chmod +x "$BIN_NAME"

# Move to install directory (requires sudo for /usr/local/bin)
echo "Installing $BIN_NAME to $INSTALL_DIR..."
sudo mv "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

# Verify installation
if command -v "$BIN_NAME" >/dev/null 2>&1; then
    echo "✅ Installation complete! Run '$BIN_NAME' to start."
else
    echo "⚠️ Installation failed. Please check permissions."
fi
