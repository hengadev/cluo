package piece

import (
	"context"

	"github.com/hengadev/encx"
	"github.com/hengadev/errsx"
)

// ProcessPieceEncx encrypts the sensitive fields of a Piece and returns a PieceEncx ready for persistence.
func ProcessPieceEncx(ctx context.Context, crypto encx.CryptoService, p *Piece) (*PieceEncx, error) {
	var errs errsx.Map

	result := &PieceEncx{
		ID:         p.ID,
		CaseID:     p.CaseID,
		StorageKey: p.StorageKey,
		MimeType:   p.MimeType,
		SizeBytes:  p.SizeBytes,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}

	// Generate DEK
	dek, err := crypto.GenerateDEK()
	if err != nil {
		errs.Set("DEK generation", err)
		return result, errs.AsError()
	}

	// Encrypt Filename
	filenameBytes, err := encx.SerializeValue(p.Filename)
	if err != nil {
		errs.Set("Filename serialization", err)
	} else {
		result.FilenameEncrypted, err = crypto.EncryptData(ctx, filenameBytes, dek)
		if err != nil {
			errs.Set("Filename encryption", err)
		}
	}

	// Encrypt Notes (optional)
	if p.Notes != "" {
		notesBytes, err := encx.SerializeValue(p.Notes)
		if err != nil {
			errs.Set("Notes serialization", err)
		} else {
			result.NotesEncrypted, err = crypto.EncryptData(ctx, notesBytes, dek)
			if err != nil {
				errs.Set("Notes encryption", err)
			}
		}
	}

	// Encrypt DEK with KEK
	result.DEKEncrypted, err = crypto.EncryptDEK(ctx, dek)
	if err != nil {
		errs.Set("DEK encryption", err)
	}

	// Get current KEK version
	result.KeyVersion, err = crypto.GetCurrentKEKVersion(ctx, crypto.GetAlias())
	if err != nil {
		errs.Set("KEK version retrieval", err)
	}

	return result, errs.AsError()
}

// DecryptPieceEncx decrypts a PieceEncx back into a plain Piece.
func DecryptPieceEncx(ctx context.Context, crypto encx.CryptoService, p *PieceEncx) (*Piece, error) {
	var errs errsx.Map

	result := &Piece{
		ID:         p.ID,
		CaseID:     p.CaseID,
		StorageKey: p.StorageKey,
		MimeType:   p.MimeType,
		SizeBytes:  p.SizeBytes,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}

	// Decrypt DEK
	dek, err := crypto.DecryptDEKWithVersion(ctx, p.DEKEncrypted, p.KeyVersion)
	if err != nil {
		errs.Set("DEK decryption", err)
		return result, errs.AsError()
	}

	// Decrypt Filename
	if len(p.FilenameEncrypted) > 0 {
		filenameBytes, err := crypto.DecryptData(ctx, p.FilenameEncrypted, dek)
		if err != nil {
			errs.Set("Filename decryption", err)
		} else {
			if err = encx.DeserializeValue(filenameBytes, &result.Filename); err != nil {
				errs.Set("Filename deserialization", err)
			}
		}
	}

	// Decrypt Notes (optional)
	if len(p.NotesEncrypted) > 0 {
		notesBytes, err := crypto.DecryptData(ctx, p.NotesEncrypted, dek)
		if err != nil {
			errs.Set("Notes decryption", err)
		} else {
			if err = encx.DeserializeValue(notesBytes, &result.Notes); err != nil {
				errs.Set("Notes deserialization", err)
			}
		}
	}

	return result, errs.AsError()
}
