package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Signature represents a digital or wet signature on a document.
type Signature struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Name             string     `encx:"encrypt" json:"name" db:"name_encrypted"`
	Role             string     `json:"role" db:"role"` // e.g., "client", "investigator", "witness"
	SignedAt         time.Time  `json:"signed_at" db:"signed_at"`
	SignatureFileURL string     `encx:"encrypt" json:"signature_file_url" db:"signature_file_url_encrypted"` // URL to stored signature (image/PDF) or digital signature envelope ID
	Method           string     `json:"method" db:"method"`                                                  // "e-sign", "wet", "pdf-stamp", "third-party"
	SignerID         *uuid.UUID `json:"signer_id" db:"signer_id,omitempty"`                                  // Reference to the user who signed
	IPAddress        *string    `encx:"encrypt" json:"ip_address" db:"ip_address_encrypted,omitempty"`       // IP address when signature was captured
	UserAgent        *string    `encx:"encrypt" json:"user_agent" db:"user_agent_encrypted,omitempty"`       // Browser/user agent when signature was captured
}

// Validate performs validation on the signature.
// TODO: Implement comprehensive signature validation:
// - Name: required, min length 2, max length 100, no special characters
// - Role: required, must be one of predefined roles (client, investigator, witness, notary, etc.)
// - SignedAt: must not be in the future, must be after document creation
// - SignatureFileURL: required if Method != "wet", must be valid URL format
// - Method: required, must be one of ["e-sign", "wet", "pdf-stamp", "third-party"]
// - SignerID: required for non-wet signatures, must exist in users table
// - IPAddress: optional, must be valid IP format if provided
// - UserAgent: optional, reasonable length limit if provided
func (s Signature) Validate() error {
	if len(s.Name) < 2 || len(s.Name) > 100 {
		return fmt.Errorf("name must be between 2 and 100 characters")
	}

	validRoles := map[string]bool{"client": true, "investigator": true, "witness": true, "notary": true}
	if !validRoles[s.Role] {
		return fmt.Errorf("role must be one of: client, investigator, witness, notary")
	}

	validMethods := map[string]bool{"e-sign": true, "wet": true, "pdf-stamp": true, "third-party": true}
	if !validMethods[s.Method] {
		return fmt.Errorf("method must be one of: e-sign, wet, pdf-stamp, third-party")
	}

	if s.Method != "wet" && s.SignatureFileURL == "" {
		return fmt.Errorf("signature file URL is required for non-wet signatures")
	}

	if s.SignedAt.IsZero() {
		return fmt.Errorf("signed at timestamp is required")
	}

	return nil
}

// NewSignature creates a new signature with default values.
func NewSignature(name, role, method, signatureFileURL string, signerID *uuid.UUID) Signature {
	return Signature{
		ID:               uuid.New(),
		Name:             name,
		Role:             role,
		SignedAt:         time.Now(),
		SignatureFileURL: signatureFileURL,
		Method:           method,
		SignerID:         signerID,
	}
}

// IsElectronic returns true if this is an electronic signature.
func (s Signature) IsElectronic() bool {
	return s.Method == "e-sign" || s.Method == "pdf-stamp" || s.Method == "third-party"
}

// IsWetSignature returns true if this is a physical/wet signature.
func (s Signature) IsWetSignature() bool {
	return s.Method == "wet"
}

