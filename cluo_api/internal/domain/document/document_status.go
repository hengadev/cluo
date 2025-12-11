package document

import (
	"fmt"
)

// DocumentStatus represents the current state of a document in its lifecycle.
type DocumentStatus string

const (
	// DocumentStatusDraft is the initial state when a document is being created.
	DocumentStatusDraft DocumentStatus = "draft"

	// DocumentStatusSent indicates the document has been sent to the client/parties.
	DocumentStatusSent DocumentStatus = "sent"

	// DocumentStatusSigned indicates the document has been signed by required parties.
	DocumentStatusSigned DocumentStatus = "signed"

	// DocumentStatusActive indicates the document is currently in effect.
	DocumentStatusActive DocumentStatus = "active"

	// DocumentStatusArchived indicates the document is no longer active but preserved.
	DocumentStatusArchived DocumentStatus = "archived"

	// DocumentStatusCancelled indicates the document has been voided or cancelled.
	DocumentStatusCancelled DocumentStatus = "cancelled"

	// DocumentStatusRejected indicates the document was rejected by a party.
	DocumentStatusRejected DocumentStatus = "rejected"

	// DocumentStatusExpired indicates the document has passed its validity period.
	DocumentStatusExpired DocumentStatus = "expired"
)

// IsValid checks if the DocumentStatus is a valid enum value.
func (s DocumentStatus) IsValid() bool {
	switch s {
	case DocumentStatusDraft, DocumentStatusSent, DocumentStatusSigned,
		DocumentStatusActive, DocumentStatusArchived, DocumentStatusCancelled,
		DocumentStatusRejected, DocumentStatusExpired:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if a status transition is allowed.
// TODO: Implement comprehensive status transition validation:
// - Draft can transition to: Sent, Cancelled
// - Sent can transition to: Signed, Rejected, Cancelled, Expired
// - Signed can transition to: Active, Cancelled
// - Active can transition to: Archived, Cancelled, Expired
// - Archived: final state, no transitions allowed
// - Cancelled: final state, no transitions allowed
// - Rejected can transition to: Draft (for revision), Cancelled
// - Expired can transition to: Archived
func (s DocumentStatus) CanTransitionTo(newStatus DocumentStatus) bool {
	// TODO: Add comprehensive transition rules with business logic validation
	// For now, allow most transitions as a basic safety check
	if !newStatus.IsValid() {
		return false
	}

	// Basic transition rules
	switch s {
	case DocumentStatusDraft:
		return newStatus == DocumentStatusSent || newStatus == DocumentStatusCancelled
	case DocumentStatusSent:
		return newStatus == DocumentStatusSigned || newStatus == DocumentStatusRejected ||
			newStatus == DocumentStatusCancelled || newStatus == DocumentStatusExpired
	case DocumentStatusSigned:
		return newStatus == DocumentStatusActive || newStatus == DocumentStatusCancelled
	case DocumentStatusActive:
		return newStatus == DocumentStatusArchived || newStatus == DocumentStatusCancelled ||
			newStatus == DocumentStatusExpired
	case DocumentStatusArchived, DocumentStatusCancelled:
		return false // Final states
	case DocumentStatusRejected:
		return newStatus == DocumentStatusDraft || newStatus == DocumentStatusCancelled
	case DocumentStatusExpired:
		return newStatus == DocumentStatusArchived
	default:
		return false
	}
}

// String returns the string representation of the status.
func (s DocumentStatus) String() string {
	return string(s)
}

// IsFinal returns true if the status is a final state that cannot be changed.
func (s DocumentStatus) IsFinal() bool {
	return s == DocumentStatusArchived || s == DocumentStatusCancelled
}

// IsActive returns true if the document is currently active or in use.
func (s DocumentStatus) IsActive() bool {
	return s == DocumentStatusActive || s == DocumentStatusSigned
}

// MarshalJSON implements json.Marshaler interface.
func (s DocumentStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *DocumentStatus) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1] // Remove quotes
	}

	status := DocumentStatus(str)
	if !status.IsValid() {
		return fmt.Errorf("invalid document status: %s", str)
	}

	*s = status
	return nil
}
