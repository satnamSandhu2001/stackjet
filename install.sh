#!/bin/bash

OS=$(uname -s)

if [ "$OS" = "Linux" ]; then
  BINARY="go-stack-cli-linux"
elif [ "$OS" = "Darwin" ]; then # macOS
  BINARY="go-stack-cli-mac"
else
  echo "Error: Unsupported operating system ($OS). This script supports Linux and macOS."
  exit 1
fi

echo "Downloading go-stack-cli for $OS..."
curl -L -o go-stack-cli "https://github.com/satnamSandhu2001/go-stack-cli/releases/latest/download/$BINARY"

if [ $? -ne 0 ]; then
    echo "Error: Failed to download go-stack-cli. Please check your internet connection or the GitHub repository."
    exit 1
fi

echo "Making go-stack-cli executable..."
chmod +x go-stack-cli

echo "Moving go-stack-cli to /usr/local/bin/ (requires sudo)..."
sudo mv go-stack-cli /usr/local/bin/

if [ $? -ne 0 ]; then
    echo "Error: Failed to move go-stack-cli. Do you have sudo permissions?"
    exit 1
fi

echo "go-stack-cli installed successfully!"
echo "You can now run 'go-stack-cli --help' to see available options."
