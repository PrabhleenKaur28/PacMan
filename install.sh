#!/usr/bin/env bash
set -euo pipefail

REPO="PrabhleenKaur28/PacMan"
BIN_NAME="t-pac"               # updated binary name
DEST="/usr/local/bin/t-pac"     # updated install destination

echo "Fetching latest release info for $REPO…"

# find latest release asset matching BIN_NAME
ASSET_URL=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep "browser_download_url" \
  | cut -d '"' -f 4 \
  | grep "${BIN_NAME}" || true)

if [ -z "$ASSET_URL" ]; then
  echo "❌ No compatible binary found for ${BIN_NAME} in the latest release."
  exit 1
fi

echo "➡️  Found release asset:"
echo "$ASSET_URL"
echo

echo "⬇️  Downloading…"
tmpfile=$(mktemp)
curl -L "$ASSET_URL" -o "$tmpfile"
chmod +x "$tmpfile"

echo "📦 Installing to $DEST (requires sudo)…"
sudo mv "$tmpfile" "$DEST"
sudo chmod +x "$DEST"

echo "✅ Installed T-Pac!"
echo "➡️ You can now run: t-pac"
