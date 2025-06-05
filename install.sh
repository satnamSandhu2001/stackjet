#!/bin/bash

OS=$(uname -s)

# Detect OS
if [ "$OS" != "Linux" ]; then
  echo "Unsupported OS: $OS. Only Linux is supported."
  exit 1
fi

# Detect architecture
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64 | arm64)
    ARCH="arm64"
    ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Construct binary name
BINARY="stackjet-linux-${ARCH}"

echo "Downloading stackjet for $OS..."
curl -L -o stackjet "https://github.com/satnamSandhu2001/stackjet/releases/latest/download/$BINARY"

if [ $? -ne 0 ]; then
    echo "Error: Failed to download stackjet. Please check your internet connection or the GitHub repository."
    exit 1
fi

echo "Making stackjet executable..."
chmod +x stackjet

echo "Moving stackjet to /usr/local/bin/ (requires sudo)..."
sudo mv stackjet /usr/local/bin/

if [ $? -ne 0 ]; then
    echo "Error: Failed to move stackjet. Do you have sudo permissions?"
    exit 1
fi

echo "stackjet installed successfully!"
echo "You can now run 'stackjet help' to see available options."
