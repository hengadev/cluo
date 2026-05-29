// Package main provides a CLI tool for managing update manifest signatures.
//
// Usage:
//
//	sign-manifest genkey <private-key-file> <public-key-file>
//	  Generate a new Ed25519 key pair and write hex-encoded keys to files.
//
//	sign-manifest sign <manifest.json> <private-key-file>
//	  Read an unsigned manifest, sign it, and overwrite the file with the signed version.
//
//	sign-manifest verify <manifest.json> <public-key-hex>
//	  Verify a signed manifest against the given public key.
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cluo_desktop/updater"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "genkey":
		cmdGenKey()
	case "sign":
		cmdSign()
	case "verify":
		cmdVerify()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  sign-manifest genkey <private-key-file> <public-key-file>\n")
	fmt.Fprintf(os.Stderr, "  sign-manifest sign   <manifest.json> <private-key-file>\n")
	fmt.Fprintf(os.Stderr, "  sign-manifest verify <manifest.json> <public-key-hex>\n")
}

func cmdGenKey() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: sign-manifest genkey <private-key-file> <public-key-file>\n")
		os.Exit(1)
	}

	privFile := os.Args[2]
	pubFile := os.Args[3]

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fatal("failed to generate key pair: %v", err)
	}

	if err := os.WriteFile(privFile, []byte(hex.EncodeToString(privateKey)), 0600); err != nil {
		fatal("failed to write private key: %v", err)
	}
	if err := os.WriteFile(pubFile, []byte(hex.EncodeToString(publicKey)), 0644); err != nil {
		fatal("failed to write public key: %v", err)
	}

	fmt.Printf("Generated Ed25519 key pair:\n  Private: %s\n  Public:  %s\n", privFile, pubFile)
	fmt.Printf("Public key (for ldflags): %s\n", hex.EncodeToString(publicKey))
}

func cmdSign() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: sign-manifest sign <manifest.json> <private-key-file>\n")
		os.Exit(1)
	}

	manifestPath := os.Args[2]
	privKeyPath := os.Args[3]

	// Read manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fatal("failed to read manifest: %v", err)
	}

	var m updater.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		fatal("failed to parse manifest: %v", err)
	}

	// Read private key (trim whitespace so trailing newlines don't break hex decode)
	privKeyHex, err := os.ReadFile(privKeyPath)
	if err != nil {
		fatal("failed to read private key: %v", err)
	}

	// Sign
	sig, err := updater.SignManifest(&m, strings.TrimSpace(string(privKeyHex)))
	if err != nil {
		fatal("failed to sign manifest: %v", err)
	}

	// Write signed manifest
	m.Signature = sig
	signedData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fatal("failed to marshal signed manifest: %v", err)
	}

	if err := os.WriteFile(manifestPath, append(signedData, '\n'), 0644); err != nil {
		fatal("failed to write signed manifest: %v", err)
	}

	fmt.Printf("Signed manifest written to %s\n", manifestPath)
}

func cmdVerify() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: sign-manifest verify <manifest.json> <public-key-hex>\n")
		os.Exit(1)
	}

	manifestPath := os.Args[2]
	pubKeyHex := os.Args[3]

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fatal("failed to read manifest: %v", err)
	}

	var m updater.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		fatal("failed to parse manifest: %v", err)
	}

	if err := updater.VerifyManifestSignature(&m, pubKeyHex); err != nil {
		fatal("verification failed: %v", err)
	}

	fmt.Println("✓ Manifest signature is valid")
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
