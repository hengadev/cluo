# Installation & Updates

## First install

### Windows

Download the latest binary directly from the releases server:

```
https://releases.henga.dev/cluo/manifest.json   ← lists the current download URL
```

The manifest `downloads.windows_amd64` field points to the `.exe`. Download it and run it —
no installer, no UAC prompt, just a self-contained binary. Place it wherever you want it to
live (e.g. `C:\Users\<you>\AppData\Local\cluo\cluo_desktop.exe`), because the in-app updater
replaces the file in-place at that path when a new version ships.

### Linux (dev / local build)

There is no Linux binary in CI yet. Build locally:

```sh
cd cluo_desktop
wails build
# output: build/bin/cluo_desktop
```

Run `./build/bin/cluo_desktop`. The in-app updater will still check for updates and download
a new binary when one is available for `linux_amd64`.

---

## Shipping a release

The CI pipeline (`/.woodpecker.yml`) triggers on any tag matching `v*`. To publish a new version:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The pipeline will:
1. Cross-compile the Windows binary with the version and manifest URL baked in via ldflags.
2. Compute the SHA256 checksum.
3. Sign the manifest with the Ed25519 private key (`UPDATE_PRIVATE_KEY` Woodpecker secret).
4. Copy both the `.exe` and the signed `manifest.json` to `https://releases.henga.dev/cluo/`.

The manifest at `releases.henga.dev/cluo/manifest.json` always points to the latest version.
Any running instance will see the new version the next time it checks.

---

## In-app update flow

The updater (in `updater/`) runs entirely inside the app. When the user triggers a check:

1. **`CheckForUpdate()`** — fetches `manifest.json`, verifies the Ed25519 signature against
   the public key embedded at build time, and compares semver. Returns `UpdateInfo{Available: true}`
   if the remote version is newer.

2. **`DownloadAndInstall()`** — downloads the binary for the current platform to a temp file,
   verifies the SHA256 checksum, then replaces the running executable in-place:
   - **Windows**: renames the running `.exe` to `.exe.old`, moves the new file into its place.
   - **Linux/macOS**: backs up the current binary, moves the new file in, restores permissions.

3. **`RestartApp()`** — launches a new process from the (now updated) executable path, then
   calls `runtime.Quit` on the current process.

Frontend events emitted during the flow:

| Event | Payload |
|---|---|
| `updater:status` | `"downloading"` → `"installing"` → `"ready"` |
| `updater:progress` | `{ downloaded, total, percent }` |
| `updater:error` | error message string |

In dev mode (binary built without ldflags), `ManifestURL` is empty and `PublicKey` is empty,
so `CheckForUpdate()` returns an error immediately and signature verification is skipped.
See `updater/config.go`.
