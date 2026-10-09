#!/usr/bin/env bash

set -e

REPO="henristr/hftch"
BINARY="hftch"
INSTALL_DIR="/usr/local/bin"

ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)
        GOARCH="amd64"
        ;;
    aarch64|arm64)
        GOARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name":' \
    | sed -E 's/.*"([^"]+)".*/\1/')"

if [ -z "$VERSION" ]; then
    echo "Could not determine latest version."
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY}_${VERSION}_linux_${GOARCH}.tar.gz"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading hftch $VERSION..."

curl -fL "$URL" -o "$TMP_DIR/$BINARY.tar.gz"

tar -xzf "$TMP_DIR/$BINARY.tar.gz" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/$BINARY" ]; then
    echo "Binary not found in archive."
    exit 1
fi

echo "Installing hftch to $INSTALL_DIR..."

sudo install -m 755 "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"

echo
echo "hftch $VERSION successfully installed!"
echo

"$INSTALL_DIR/$BINARY"