#!/bin/sh
set -e

REPO="agejevasv/swk"
INSTALL_DIR="${SWK_INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "^location:" | sed 's|.*/||' | tr -d '\r')

if [ -z "$VERSION" ]; then
  echo "Failed to detect latest version" >&2
  exit 1
fi

BINARY="swk-${OS}-${ARCH}"
URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "Installing swk $VERSION ($OS/$ARCH)..."

TMPFILE=$(mktemp)
trap 'rm -f "$TMPFILE"' EXIT

curl -sL "$URL" -o "$TMPFILE"

if [ ! -s "$TMPFILE" ]; then
  echo "Download failed" >&2
  exit 1
fi

chmod +x "$TMPFILE"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMPFILE" "$INSTALL_DIR/swk"
else
  sudo mv "$TMPFILE" "$INSTALL_DIR/swk"
fi

echo "Installed swk to $INSTALL_DIR/swk"
swk --version
