#!/bin/sh
#
# install.sh - Install Mandor CLI
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/sanxzy/mandor/main/scripts/install.sh | sh -s -- --help
#
# Options:
#   --prefix DIR     Install prefix (default: $HOME/.local)
#   --version VER    Install specific version (default: latest)
#   --prerelease     Install latest prerelease
#   --help           Show this help
#

set -e

REPO="sanxzy/mandor"
INSTALL_DIR="${HOME}/.local/bin"
VERSION="latest"
PRERELEASE=""
TAG="latest"

while [ $# -gt 0 ]; do
    case "$1" in
        --prefix)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --prerelease)
            PRERELEASE="1"
            TAG="latest"
            shift
            ;;
        --help|-h)
            head -20 "$0"
            exit 0
            ;;
    esac
done

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="x64" ;;
    arm64|aarch64) ARCH="arm64" ;;
esac

case "$OS" in
    darwin) ;;
    linux) ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

echo "Mandor Installer"
echo "================"
echo "OS: $OS-$ARCH"

if [ "$VERSION" = "latest" ]; then
    if [ -n "$PRERELEASE" ]; then
        echo "Fetching latest prerelease..."
        VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)
    else
        echo "Fetching latest release..."
        VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
    fi
fi

echo "Version: $VERSION"
echo "Install dir: $INSTALL_DIR"
echo ""

ASSET_NAME="${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"
TEMP_DIR=$(mktemp -d)
TARFILE="${TEMP_DIR}/${ASSET_NAME}"

echo "Downloading ${ASSET_NAME}..."
if ! curl -fsSL -o "$TARFILE" "$DOWNLOAD_URL"; then
    echo "Download failed: $DOWNLOAD_URL"
    rm -rf "$TEMP_DIR"
    exit 1
fi

echo "Extracting..."
mkdir -p "$INSTALL_DIR"
tar -xzf "$TARFILE" -C "$INSTALL_DIR"
chmod 755 "${INSTALL_DIR}/mandor"

rm -rf "$TEMP_DIR"

echo ""
echo "Installed: ${INSTALL_DIR}/mandor"
echo ""

if [ -d "$HOME/.local/bin" ]; then
    echo "Add to PATH: export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
