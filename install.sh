#!/bin/sh
set -e

# koto - Interactive ToDo management CLI tool installer
# This script installs the latest version of koto from GitHub releases

# Configuration
REPO="syeeel/koto-cli-go"
BINARY_NAME="koto"
INSTALL_DIR="$HOME/.local/bin"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo "${GREEN}✓${NC} $1"
}

print_warning() {
    echo "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo "${RED}✗${NC} $1"
}

echo ""
echo "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "${GREEN}  koto - ToDo Management CLI Installer${NC}"
echo "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Detect OS and architecture
print_info "Detecting system information..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        print_error "Unsupported architecture: $ARCH"
        echo "Supported architectures: x86_64 (amd64), aarch64/arm64"
        exit 1
        ;;
esac

print_success "Detected: ${OS}/${ARCH}"

# Fetch latest version from GitHub API
print_info "Fetching latest version from GitHub..."
VERSION=$(curl -sL https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    print_error "Failed to get latest version from GitHub"
    echo ""
    echo "Please check:"
    echo "  1. Your internet connection"
    echo "  2. GitHub API availability"
    echo "  3. Repository: https://github.com/$REPO"
    exit 1
fi

print_success "Latest version: ${VERSION}"

# Construct download URL
# Format: koto_v1.0.0_darwin_amd64.tar.gz
ARCHIVE_NAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/${ARCHIVE_NAME}"

echo ""
print_info "Downloading koto ${VERSION}..."
echo "URL: ${DOWNLOAD_URL}"
echo ""

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download archive
if ! curl -L --progress-bar "$DOWNLOAD_URL" -o "${ARCHIVE_NAME}"; then
    print_error "Download failed"
    echo ""
    echo "Attempted to download: ${DOWNLOAD_URL}"
    echo ""
    echo "Please check:"
    echo "  1. The release exists for your platform (${OS}/${ARCH})"
    echo "  2. Visit: https://github.com/$REPO/releases/latest"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Extract archive
print_info "Extracting archive..."
if ! tar -xzf "${ARCHIVE_NAME}"; then
    print_error "Failed to extract archive"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Verify binary exists
if [ ! -f "$BINARY_NAME" ]; then
    print_error "Binary not found in archive"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Create installation directory
print_info "Installing to ${INSTALL_DIR}..."
mkdir -p "$INSTALL_DIR"

# Install binary
mv "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
print_success "koto ${VERSION} installed successfully!"
echo "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Installation path: ${INSTALL_DIR}/${BINARY_NAME}"
echo ""

# Check if install directory is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    print_warning "${INSTALL_DIR} is not in your PATH"
    echo ""
    echo "Add the following line to your shell configuration file:"
    echo ""

    # Detect shell and suggest appropriate config file
    if [ -n "$BASH_VERSION" ]; then
        echo "  ${YELLOW}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc${NC}"
        echo "  ${YELLOW}source ~/.bashrc${NC}"
    elif [ -n "$ZSH_VERSION" ]; then
        echo "  ${YELLOW}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc${NC}"
        echo "  ${YELLOW}source ~/.zshrc${NC}"
    else
        echo "  ${YELLOW}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
        echo ""
        echo "  Add this to your shell config file (~/.bashrc, ~/.zshrc, etc.)"
    fi
    echo ""
else
    print_success "Ready to use! Try running: ${BINARY_NAME}"
    echo ""
    echo "Quick start:"
    echo "  ${BLUE}${BINARY_NAME}${NC}              # Start the TUI"
    echo ""
fi

echo "For more information, visit:"
echo "  https://github.com/$REPO"
echo ""
