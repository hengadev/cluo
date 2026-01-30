package updater

// Version and ManifestURL are injected at build time using ldflags.
// Example: go build -ldflags "-X cluo_desktop/updater.Version=1.2.0 -X cluo_desktop/updater.ManifestURL=https://..."
var (
	Version     = "0.0.0-dev"
	ManifestURL = ""
)
