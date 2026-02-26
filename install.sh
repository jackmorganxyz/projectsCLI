#!/bin/sh
set -e

REPO="jackmorganxyz/projectsCLI"
BINARY="projects"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)             echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag
echo "Fetching latest release..."
TAG=$(curl -sI "https://github.com/${REPO}/releases/latest" \
  | grep -i '^location:' \
  | sed 's/.*tag\///' \
  | tr -d '\r\n')

if [ -z "$TAG" ]; then
  echo "Error: could not determine latest release"
  exit 1
fi

VERSION="${TAG#v}"
ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE}"

echo "Downloading ${BINARY} ${TAG} for ${OS}/${ARCH}..."

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${TAG}/checksums.txt"

curl -sL "$URL" -o "${TMPDIR}/${ARCHIVE}"
curl -sL "$CHECKSUMS_URL" -o "${TMPDIR}/checksums.txt"

# Verify integrity of downloaded archive
echo "Verifying checksum..."
EXPECTED=$(grep "${ARCHIVE}" "${TMPDIR}/checksums.txt" | awk '{print $1}')
if [ -z "$EXPECTED" ]; then
  echo "Error: no checksum found for ${ARCHIVE} in checksums.txt"
  exit 1
fi

ACTUAL=$(sha256sum "${TMPDIR}/${ARCHIVE}" | awk '{print $1}')
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "Error: checksum mismatch"
  echo "  expected: $EXPECTED"
  echo "  actual:   $ACTUAL"
  exit 1
fi
echo "Checksum OK."

tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"

# Install binary
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "${BINARY} ${TAG} installed to ${INSTALL_DIR}/${BINARY}"
echo "Run 'projects --help' to get started."
