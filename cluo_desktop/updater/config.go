package updater

// Version, ManifestURL, and PublicKey are injected at build time using ldflags.
// Example: go build -ldflags "-X cluo_desktop/updater.Version=1.2.0 -X cluo_desktop/updater.ManifestURL=https://... -X cluo_desktop/updater.PublicKey=abcdef..."
var (
	Version     = "0.0.0-dev"
	ManifestURL = ""
	PublicKey   = "" // hex-encoded Ed25519 public key; empty = skip verification (dev mode)
)
