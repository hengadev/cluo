# Auto-Updater Setup Guide

This guide explains how to configure and deploy the auto-updater for cluo_desktop.

## Overview

The auto-updater works by:
1. Fetching a JSON manifest from a URL you control
2. Comparing the manifest version with the current app version
3. Downloading the new binary if an update is available
4. Replacing the current binary and prompting for restart

## Prerequisites

- A static file hosting service (any web server, S3, GitHub Releases, etc.)
- Ability to build the app with custom ldflags

---

## Step 1: Generate TypeScript Bindings

After modifying Go code, regenerate the frontend bindings:

```bash
cd cluo_desktop
wails generate module
```

This creates `frontend/src/lib/wailsjs/go/updater/Updater.js` and `.d.ts` files.

---

## Step 2: Create the Manifest File

Create a `manifest.json` file with this structure:

```json
{
  "version": "1.0.0",
  "release_notes": "Initial release with new features:\n- Feature A\n- Bug fix B",
  "downloads": {
    "linux_amd64": "https://your-server.com/releases/v1.0.0/cluo_desktop_linux_amd64",
    "linux_arm64": "https://your-server.com/releases/v1.0.0/cluo_desktop_linux_arm64",
    "darwin_amd64": "https://your-server.com/releases/v1.0.0/cluo_desktop_darwin_amd64",
    "darwin_arm64": "https://your-server.com/releases/v1.0.0/cluo_desktop_darwin_arm64",
    "windows_amd64": "https://your-server.com/releases/v1.0.0/cluo_desktop_windows_amd64.exe"
  },
  "checksums": {
    "linux_amd64": "sha256:abc123...",
    "linux_arm64": "sha256:def456...",
    "darwin_amd64": "sha256:789ghi...",
    "darwin_arm64": "sha256:jkl012...",
    "windows_amd64": "sha256:mno345..."
  }
}
```

### Manifest Fields

| Field | Required | Description |
|-------|----------|-------------|
| `version` | Yes | Semantic version (e.g., "1.2.3" or "1.2.3-beta") |
| `release_notes` | No | Text shown to users (supports newlines with `\n`) |
| `downloads` | Yes | Map of platform to download URL |
| `checksums` | No | Map of platform to SHA256 checksum (recommended) |

### Platform Identifiers

| Platform | Identifier |
|----------|------------|
| Linux (Intel/AMD 64-bit) | `linux_amd64` |
| Linux (ARM 64-bit) | `linux_arm64` |
| macOS (Intel) | `darwin_amd64` |
| macOS (Apple Silicon) | `darwin_arm64` |
| Windows (64-bit) | `windows_amd64` |

---

## Step 3: Build for Release

Build the application with version and manifest URL injected:

```bash
# Set your variables
VERSION="1.0.0"
MANIFEST_URL="https://your-server.com/releases/manifest.json"

# Build for current platform
wails build -ldflags "-X cluo_desktop/updater.Version=$VERSION -X cluo_desktop/updater.ManifestURL=$MANIFEST_URL"

# Build for specific platform (cross-compile)
wails build -platform linux/amd64 -ldflags "-X cluo_desktop/updater.Version=$VERSION -X cluo_desktop/updater.ManifestURL=$MANIFEST_URL"
```

### Build Script Example

Create a `build-release.sh` script:

```bash
#!/bin/bash
set -e

VERSION="${1:?Usage: $0 <version>}"
MANIFEST_URL="https://your-server.com/releases/manifest.json"
LDFLAGS="-X cluo_desktop/updater.Version=$VERSION -X cluo_desktop/updater.ManifestURL=$MANIFEST_URL"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

mkdir -p dist/v$VERSION

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT="cluo_desktop_${OS}_${ARCH}"

    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for $PLATFORM..."
    wails build -platform "$PLATFORM" -ldflags "$LDFLAGS" -o "dist/v$VERSION/$OUTPUT"
done

# Generate checksums
cd dist/v$VERSION
sha256sum * > checksums.txt
echo "Checksums:"
cat checksums.txt
```

---

## Step 4: Upload Release Files

Upload to your static file host:

```
your-server.com/releases/
├── manifest.json          # Always points to latest version
├── v1.0.0/
│   ├── cluo_desktop_linux_amd64
│   ├── cluo_desktop_linux_arm64
│   ├── cluo_desktop_darwin_amd64
│   ├── cluo_desktop_darwin_arm64
│   └── cluo_desktop_windows_amd64.exe
├── v1.1.0/
│   └── ...
```

---

## Step 5: Update the Manifest

When releasing a new version:

1. Build new binaries with the new version number
2. Upload binaries to a new version folder
3. Generate SHA256 checksums: `sha256sum cluo_desktop_* > checksums.txt`
4. Update `manifest.json` with:
   - New version number
   - New download URLs
   - New checksums
   - Release notes

---

## Testing Locally

### 1. Start a Local HTTP Server

```bash
# Create test directory
mkdir -p /tmp/update-test/v1.1.0

# Copy a test binary (or the same binary for testing)
cp build/bin/cluo_desktop /tmp/update-test/v1.1.0/cluo_desktop_linux_amd64

# Generate checksum
cd /tmp/update-test/v1.1.0
sha256sum cluo_desktop_linux_amd64

# Create manifest.json in /tmp/update-test/
cat > /tmp/update-test/manifest.json << 'EOF'
{
  "version": "1.1.0",
  "release_notes": "Test update",
  "downloads": {
    "linux_amd64": "http://localhost:8080/v1.1.0/cluo_desktop_linux_amd64"
  },
  "checksums": {
    "linux_amd64": "sha256:<paste-checksum-here>"
  }
}
EOF

# Start HTTP server
cd /tmp/update-test
python3 -m http.server 8080
```

### 2. Build Test Version

```bash
# Build v1.0.0 pointing to local manifest
wails build -ldflags "-X cluo_desktop/updater.Version=1.0.0 -X cluo_desktop/updater.ManifestURL=http://localhost:8080/manifest.json"
```

### 3. Test the Flow

1. Run the v1.0.0 build
2. Click your profile icon → "Check for Updates"
3. Should show v1.1.0 available
4. Click "Download & Install"
5. Verify download progress works
6. Click "Restart Now"

---

## Hosting Options

### Simple Options (Free)

| Service | Notes |
|---------|-------|
| GitHub Releases | Host manifest as a Gist, binaries as release assets |
| Cloudflare R2 | Free tier with custom domain support |
| Backblaze B2 | Free tier, works with Cloudflare CDN |

### Self-Hosted

Any web server that can serve static files (nginx, Apache, Caddy, etc.)

### Example: GitHub Releases

1. Create releases with binaries attached
2. Host `manifest.json` as a GitHub Gist (raw URL)
3. Use raw GitHub URLs for downloads:
   ```
   https://github.com/user/repo/releases/download/v1.0.0/cluo_desktop_linux_amd64
   ```

---

## Troubleshooting

### "Update checking is not configured"

The `ManifestURL` was not set at build time. Rebuild with `-ldflags`.

### Checksum mismatch

Regenerate checksums after building. Ensure you're using SHA256:
```bash
sha256sum filename
```

### Download fails on macOS/Linux

Ensure the binary has execute permissions after download. The updater sets permissions automatically, but verify your server isn't stripping them.

### Windows: "Access denied" during install

The app may be running from a protected directory. Users should run as administrator or install to a user-writable location.
