#!/usr/bin/env bash
set -e

# Auto-detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case $OS in
    linux) PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Download and install
BINARY="mayhem-v1.0.0-${PLATFORM}-${ARCH}"
URL="https://github.com/pgaijin66/mayhem/releases/download/v1.0.0/${BINARY}.tar.gz"

echo "🔥 Installing Mayhem for ${PLATFORM}/${ARCH}..."
echo "📥 Downloading from: $URL"

wget -qO- "$URL" | tar xz
sudo mv mayhem /usr/local/bin/

echo "✅ Mayhem installed successfully!"
mayhem -version