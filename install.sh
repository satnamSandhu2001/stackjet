#!/bin/bash
ARCH=$(uname -m)
OS=$(uname -s)

if [ "$OS" = "Linux" ]; then
  BINARY="go-stack-cli-linux"
elif [ "$OS" = "Darwin" ]; then
  BINARY="go-stack-cli-mac"
else
  echo "Unsupported OS"
  exit 1
fi

curl -L -o go-stack-cli https://github.com/satnamsandhu2001/go-stack-cli/releases/latest/download/$BINARY
chmod +x go-stack-cli
sudo mv go-stack-cli /usr/local/bin/
