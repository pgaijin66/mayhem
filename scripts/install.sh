#!/usr/bin/env bash
set -e

set -v

cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo "❌ Installation failed with exit code $exit_code"
        # Clean up any temporary files
        [ -f "mayhem" ] && rm -f mayhem
        [ -f "/tmp/mayhem.tar.gz" ] && rm -f /tmp/mayhem.tar.gz
    fi
    exit $exit_code
}

trap cleanup EXIT

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

download_file() {
    local url="$1"
    local output="$2"
    
    if command_exists curl; then
        echo "📥 Using curl to download..."
        if ! curl -fsSL "$url" -o "$output"; then
            echo "❌ Failed to download with curl"
            return 1
        fi
    elif command_exists wget; then
        echo "📥 Using wget to download..."
        if ! wget -qO "$output" "$url"; then
            echo "❌ Failed to download with wget"
            return 1
        fi
    else
        echo "❌ Neither curl nor wget found. Please install one of them."
        return 1
    fi
}

echo "🔍 Detecting platform..."
if ! OS=$(uname -s 2>/dev/null); then
    echo "❌ Failed to detect operating system"
    exit 1
fi

if ! ARCH=$(uname -m 2>/dev/null); then
    echo "❌ Failed to detect architecture"
    exit 1
fi

OS=$(echo "$OS" | tr '[:upper:]' '[:lower:]')

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) 
        echo "❌ Unsupported architecture: $ARCH"
        echo "ℹ️  Supported architectures: x86_64, aarch64, arm64"
        exit 1 
        ;;
esac

case $OS in
    linux) PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    *) 
        echo "❌ Unsupported OS: $OS"
        echo "ℹ️  Supported platforms: Linux, macOS (Darwin)"
        exit 1 
        ;;
esac

echo "✅ Platform detected: ${PLATFORM}/${ARCH}"

echo "🔍 Checking prerequisites..."

if ! command_exists tar; then
    echo "❌ tar command not found. Please install tar."
    exit 1
fi

if ! command_exists sudo; then
    echo "❌ sudo command not found. Please install sudo or run as root."
    exit 1
fi

# Check if /usr/local/bin exists and is writable
if [ ! -d "/usr/local/bin" ]; then
    echo "❌ Directory /usr/local/bin does not exist"
    echo "ℹ️  Creating directory with sudo..."
    if ! sudo mkdir -p /usr/local/bin; then
        echo "❌ Failed to create /usr/local/bin directory"
        exit 1
    fi
fi

# Test sudo access early
echo "🔐 Verifying sudo access..."
if ! sudo -n true 2>/dev/null; then
    echo "ℹ️  This script requires sudo access to install to /usr/local/bin"
    if ! sudo true; then
        echo "❌ Failed to obtain sudo access"
        exit 1
    fi
fi

BINARY="mayhem-v1.0.0-${PLATFORM}-${ARCH}"
URL="https://github.com/pgaijin66/mayhem/releases/download/v2.0.0/${BINARY}.tar.gz"
TEMP_FILE="/tmp/mayhem.tar.gz"

echo "🔥 Installing Mayhem for ${PLATFORM}/${ARCH}..."
echo "📥 Downloading from: $URL"

if ! download_file "$URL" "$TEMP_FILE"; then
    echo "❌ Failed to download Mayhem binary"
    exit 1
fi

if [ ! -f "$TEMP_FILE" ]; then
    echo "❌ Downloaded file not found"
    exit 1
fi

if [ ! -s "$TEMP_FILE" ]; then
    echo "❌ Downloaded file is empty"
    exit 1
fi

echo "✅ Download completed successfully"

echo "📦 Extracting archive..."
if ! tar -xzf "$TEMP_FILE" -C "$(pwd)"; then
    echo "❌ Failed to extract archive"
    exit 1
fi

if [ ! -f "mayhem" ]; then
    echo "❌ Mayhem binary not found after extraction"
    echo "ℹ️  Archive contents:"
    tar -tzf "$TEMP_FILE" 2>/dev/null || echo "Could not list archive contents"
    exit 1
fi

if [ ! -x "mayhem" ]; then
    echo "🔧 Making binary executable..."
    if ! chmod +x mayhem; then
        echo "❌ Failed to make binary executable"
        exit 1
    fi
fi

echo "🧪 Testing binary..."
if ! ./mayhem -version >/dev/null 2>&1; then
    echo "⚠️  Binary test failed, but proceeding with installation..."
    echo "ℹ️  The binary might need to be in PATH to work correctly"
fi

echo "🚀 Installing to /usr/local/bin..."
if ! sudo mv mayhem /usr/local/bin/; then
    echo "❌ Failed to install binary to /usr/local/bin"
    exit 1
fi

if [ ! -f "/usr/local/bin/mayhem" ]; then
    echo "❌ Installation verification failed"
    exit 1
fi

if [ ! -x "/usr/local/bin/mayhem" ]; then
    echo "❌ Installed binary is not executable"
    exit 1
fi

rm -f "$TEMP_FILE"

echo "✅ Mayhem installed successfully!"

echo "🧪 Testing installation..."
if mayhem -version; then
    echo "🎉 Installation completed and verified!"
else
    echo "⚠️  Installation completed but version check failed"
    echo "ℹ️  Make sure /usr/local/bin is in your PATH"
    echo "ℹ️  Current PATH: $PATH"
    exit 1
fi