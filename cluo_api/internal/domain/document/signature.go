package document

import (
	"time"

	"github.com/google/uuid"
)

// Signature represents a digital or physical signature on a document.
type Signature struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Role             string    `json:"role"`
	SignatureFileURL string    `json:"signature_file_url,omitempty"`
	Method           string    `json:"method,omitempty"`
	IPAddress        *string   `json:"ip_address,omitempty"`
	UserAgent        *string   `json:"user_agent,omitempty"`
	SignedAt         time.Time `json:"signed_at"`
}

// NewSignature creates a new signature with the given name and role.
func NewSignature(name, role string) Signature {
	return Signature{
		ID:       uuid.New(),
		Name:     name,
		Role:     role,
		SignedAt: time.Now(),
	}
}

// NewSignatureWithFile creates a new signature with a signature file URL.
func NewSignatureWithFile(name, role, signatureFileURL, method string) Signature {
	return Signature{
		ID:               uuid.New(),
		Name:             name,
		Role:             role,
		SignatureFileURL: signatureFileURL,
		Method:           method,
		SignedAt:         time.Now(),
	}
}
