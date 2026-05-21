package rapport

import (
	"context"
	"fmt"

	"github.com/hengadev/encx"
)

// ProcessRapportEncx encrypts the rapport content and returns a RapportEncx ready for storage.
func ProcessRapportEncx(ctx context.Context, crypto encx.CryptoService, r *Rapport) (*RapportEncx, error) {
	// Generate DEK
	dek, err := crypto.GenerateDEK()
	if err != nil {
		return nil, fmt.Errorf("rapport DEK generation: %w", err)
	}

	// Encrypt content directly (raw bytes — no serialization needed)
	contentEncrypted, err := crypto.EncryptData(ctx, r.Content, dek)
	if err != nil {
		return nil, fmt.Errorf("rapport content encryption: %w", err)
	}

	// Encrypt DEK with KEK
	dekEncrypted, err := crypto.EncryptDEK(ctx, dek)
	if err != nil {
		return nil, fmt.Errorf("rapport DEK encryption: %w", err)
	}

	// Get current KEK version
	keyVersion, err := crypto.GetCurrentKEKVersion(ctx, crypto.GetAlias())
	if err != nil {
		return nil, fmt.Errorf("rapport KEK version retrieval: %w", err)
	}

	return &RapportEncx{
		ID:               r.ID,
		CaseID:           r.CaseID,
		ContentEncrypted: contentEncrypted,
		DEKEncrypted:     dekEncrypted,
		KeyVersion:       keyVersion,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}, nil
}

// DecryptRapportEncx decrypts a RapportEncx back to a Rapport.
func DecryptRapportEncx(ctx context.Context, crypto encx.CryptoService, r *RapportEncx) (*Rapport, error) {
	// Decrypt DEK
	dek, err := crypto.DecryptDEKWithVersion(ctx, r.DEKEncrypted, r.KeyVersion)
	if err != nil {
		return nil, fmt.Errorf("rapport DEK decryption: %w", err)
	}

	// Decrypt content directly (raw bytes — no deserialization needed)
	content, err := crypto.DecryptData(ctx, r.ContentEncrypted, dek)
	if err != nil {
		return nil, fmt.Errorf("rapport content decryption: %w", err)
	}

	return &Rapport{
		ID:        r.ID,
		CaseID:    r.CaseID,
		Content:   content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}
