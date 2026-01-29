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
#   --bin-dir DIR    Binary directory (default: bin)
#   --version VER    Install specific version (default: latest)
#   --prerelease     Install latest prerelease version
#   --help           Show this help message
#
# Environment:
#   MANDOR_VERSION   Version to install (overrides --version)
#   MANDOR_PREFIX    Install prefix (overrides --prefix)
#   MANDOR_PRERELEASE Set to "true" to install prerelease
#

set -e

REPO_OWNER="sanxzy"
REPO_NAME="mandor"
BIN_NAME="mandor"
GITHUB_API="https://api.github.com"
GITHUB_RAW="https://raw.githubusercontent.com"

# Default values
VERSION="latest"
PRERELEASE=""
PREFIX="${HOME}/.local"
BIN_DIR="bin"
DOWNLOAD_DIR="${TMPDIR:-/tmp}/mandor-install-$$"

# Parse arguments
while [ $# -gt 0 ]; do
    case "$1" in
        --prefix)
            PREFIX="$2"
            shift 2
            ;;
        --bin-dir)
            BIN_DIR="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --prerelease|-p)
            PRERELEASE="true"
            ;;
        --help|-h)
            sed -n '/^#/!d; s/^#//; /^Usage/,$p' "$0" | head -${LINES:-99}
            exit 0
            ;;
        *)
            echo "Unknown option: $1" >&2
            echo "Run with --help for usage" >&2
            exit 1
            ;;
    esac
done

# Use environment variables if set
[ -n "$MANDOR_VERSION" ] && VERSION="$MANDOR_VERSION"
[ -n "$MANDOR_PREFIX" ] && PREFIX="$MANDOR_PREFIX"
[ -n "$MANDOR_PRERELEASE" ] && PRERELEASE="$MANDOR_PRERELEASE"

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Darwin|*[Dd]arwin*)
            echo "darwin"
            ;;
        Linux|*[Ll]inux*)
            echo "linux"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "win32"
            ;;
        *)
            echo "Unsupported platform: $(uname -s)" >&2
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64*|amd64*)
            echo "x64"
            ;;
        arm64*|aarch64*)
            echo "arm64"
            ;;
        i*86)
            echo "x64"
            ;;
        *)
            echo "Unsupported architecture: $(uname -m)" >&2
            exit 1
            ;;
    esac
}

# Get latest version number
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        local api_url="${GITHUB_API}/repos/${REPO_OWNER}/${REPO_NAME}/releases"
        if [ -n "$PRERELEASE" ]; then
            # Get most recent release (including prereleases)
            VERSION=$(curl -fsSL "${api_url}" | \
                sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)
        else
            # Get latest stable release only
            VERSION=$(curl -fsSL "${api_url}/latest" | \
                sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
        fi
    fi
    echo "$VERSION"
}

# Download file with progress
download_file() {
    local url="$1"
    local dest="$2"
    local label="$3"

    echo "Downloading ${label}..."
    if ! curl -fsSL --progress-bar -o "$dest" "$url"; then
        echo "Failed to download: $url" >&2
        echo "Please check your internet connection or try a different version." >&2
        rm -f "$dest"
        exit 1
    fi
}

# Install binary
install_binary() {
    local platform="$1"
    local arch="$2"
    local version="$3"
    local install_dir="${PREFIX}/${BIN_DIR}"

    # Create install directory
    if ! mkdir -p "$install_dir"; then
        echo "Failed to create directory: $install_dir" >&2
        exit 1
    fi

    # Download URL
    local asset_name="${BIN_NAME}-${platform}-${arch}.tar.gz"
    local download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/${asset_name}"
    local tarball="${DOWNLOAD_DIR}/${asset_name}"

    # Download
    download_file "$download_url" "$tarball" "Mandor ${version} (${platform}-${arch})"

    # Verify download
    if [ ! -s "$tarball" ]; then
        echo "Download failed: file is empty or missing" >&2
        echo "" >&2
        echo "Available assets for ${version}:" >&2
        curl -fsSL "${GITHUB_API}/repos/${REPO_OWNER}/${REPO_NAME}/releases/tags/${version}" | \
            sed -n 's/.*"name": *"\([^"]*\)".*/\1/p' >&2
        exit 1
    fi

    # Extract
    echo "Extracting..."
    mkdir -p "${DOWNLOAD_DIR}/${BIN_NAME}-${platform}-${arch}"
    tar -xzf "${tarball}" -C "${DOWNLOAD_DIR}/${BIN_NAME}-${platform}-${arch}"

    # Install binary
    local binary="${DOWNLOAD_DIR}/${BIN_NAME}-${platform}-${arch}/${BIN_NAME}"
    if [ "$platform" = "win32" ]; then
        binary="${binary}.exe"
    fi

    if [ ! -f "$binary" ]; then
        echo "Binary not found after extraction" >&2
        exit 1
    fi

    cp "$binary" "${install_dir}/${BIN_NAME}"
    chmod 755 "${install_dir}/${BIN_NAME}"

    # Cleanup
    rm -rf "$DOWNLOAD_DIR"

    echo ""
    echo "Mandor installed successfully!"
    echo ""
    echo "Binary location: ${install_dir}/${BIN_NAME}"
    echo ""
    if [ "$install_dir" = "${HOME}/.local/bin" ] || [ "$install_dir" = "/usr/local/bin" ]; then
        echo "Make sure ${install_dir} is in your PATH."
        echo ""
        echo "Add to PATH:"
        echo "  For bash: echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
        echo "  For zsh:  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.zshrc && source ~/.zshrc"
    fi
}

# Main
main() {
    local platform
    local arch
    local version

    echo "Mandor Installer"
    echo "================"
    echo ""

    platform=$(detect_platform)
    arch=$(detect_arch)
    version=$(get_latest_version)

    echo "Platform:   ${platform}-${arch}"
    echo "Version:    ${version}"
    if [ -n "$PRERELEASE" ]; then
        echo "Type:       prerelease"
    else
        echo "Type:       stable"
    fi
    echo "Prefix:     ${PREFIX}"
    echo ""

    # Create temp directory
    mkdir -p "$DOWNLOAD_DIR"

    # Install
    install_binary "$platform" "$arch" "$version"
}

main "$@"
