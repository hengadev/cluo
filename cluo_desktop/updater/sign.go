package updater

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// SignManifest signs the manifest content with the given Ed25519 private key.
// The private key must be 64 bytes, hex-encoded.
// Returns the signature as a hex-encoded string.
func SignManifest(m *Manifest, privateKeyHex string) (string, error) {
	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key size: expected %d bytes, got %d", ed25519.PrivateKeySize, len(privateKey))
	}

	payload, err := canonicalPayload(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest for signing: %w", err)
	}

	sig := ed25519.Sign(privateKey, payload)
	return hex.EncodeToString(sig), nil
}

// VerifyManifestSignature verifies the manifest's Ed25519 signature against the
// given hex-encoded public key. Returns an error if the signature is missing,
// malformed, or does not match.
func VerifyManifestSignature(m *Manifest, publicKeyHex string) error {
	if publicKeyHex == "" {
		// Dev mode: no public key configured, skip verification
		return nil
	}

	if m.Signature == "" {
		return fmt.Errorf("manifest is not signed — rejecting unsigned update")
	}

	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid public key size: expected %d bytes, got %d", ed25519.PublicKeySize, len(publicKey))
	}

	sig, err := hex.DecodeString(m.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	if len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature size: expected %d bytes, got %d", ed25519.SignatureSize, len(sig))
	}

	payload, err := canonicalPayload(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest for verification: %w", err)
	}

	if !ed25519.Verify(publicKey, payload, sig) {
		return fmt.Errorf("manifest signature verification failed — the manifest may have been tampered with")
	}

	return nil
}

// canonicalPayload produces the deterministic bytes that are signed/verified.
// It marshals the manifest fields excluding the Signature, so the signature
// itself is not part of the signed data.
func canonicalPayload(m *Manifest) ([]byte, error) {
	unsigned := &Manifest{
		Version:      m.Version,
		ReleaseNotes: m.ReleaseNotes,
		Downloads:    m.Downloads,
		Checksums:    m.Checksums,
	}
	// json.Marshal sorts map keys alphabetically, giving deterministic output.
	return json.Marshal(unsigned)
}
