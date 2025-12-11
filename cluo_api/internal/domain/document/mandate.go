package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Mandate represents a legal authorization document that grants permission to conduct investigation.
type Mandate struct {
	DocumentBase

	MandateNumber         string     `encx:"encrypt" json:"mandate_number" db:"mandate_number_encrypted"`
	IssueDate             time.Time  `json:"issue_date" db:"issue_date"`
	ScopeOfWork           string     `encx:"encrypt" json:"scope_of_work" db:"scope_of_work_encrypted"`
	ValidFrom             time.Time  `json:"valid_from" db:"valid_from"`
	ValidUntil            *time.Time `json:"valid_until,omitempty" db:"valid_until"`
	TermsConditions       string     `encx:"encrypt" json:"terms_conditions" db:"terms_conditions_encrypted"`
	ClientSignature       *Signature `encx:"encrypt" json:"client_signature,omitempty" db:"client_signature_encrypted"`
	InvestigatorSignature *Signature `encx:"encrypt" json:"investigator_signature,omitempty" db:"investigator_signature_encrypted"`
	LinkedEstimateID      *uuid.UUID `json:"linked_estimate_id,omitempty" db:"linked_estimate_id"`
	SpecialInstructions   *string    `encx:"encrypt" json:"special_instructions,omitempty" db:"special_instructions_encrypted"`
	Jurisdiction          *string    `json:"jurisdiction,omitempty" db:"jurisdiction"`
}

// GetType returns the document type.
func (m *Mandate) GetType() DocumentType {
	return DocumentTypeMandate
}

// Validate performs validation on the mandate.
// TODO: Implement comprehensive mandate validation:
// - MandateNumber: required, unique within system, follows format MND-YYYY-NNN
// - IssueDate: required, cannot be in the future, reasonable business constraints
// - ScopeOfWork: required, min length 20, max length 2000, must describe investigation scope
// - ValidFrom: required, cannot be before IssueDate, reasonable business constraints
// - ValidUntil: optional, must be after ValidFrom if provided, reasonable duration limits
// - TermsConditions: required, standard legal terms, reasonable length constraints
// - ClientSignature: required for mandate to be valid/active
// - InvestigatorSignature: optional but recommended for formal records
// - LinkedEstimateID: optional, if provided must reference existing estimate
// - SpecialInstructions: optional, reasonable length limit if provided
// - Jurisdiction: optional, valid jurisdiction format if provided
// - Business rules: cannot be active without client signature
// - Business rules: validity period (ValidUntil) must be reasonable (e.g., max 2 years)
// - Business rules: scope of work must be specific and actionable
// - Business rules: cannot modify mandate after signatures are added
// - Business rules: LinkedEstimateID must belong to same case and client
func (m *Mandate) Validate() error {
	// TODO: Add comprehensive validation implementation
	if m.MandateNumber == "" {
		return fmt.Errorf("mandate number is required")
	}

	if m.IssueDate.IsZero() {
		return fmt.Errorf("issue date is required")
	}

	if m.ScopeOfWork == "" {
		return fmt.Errorf("scope of work is required")
	}

	if len(m.ScopeOfWork) < 20 {
		return fmt.Errorf("scope of work must be at least 20 characters")
	}

	if m.ValidFrom.IsZero() {
		return fmt.Errorf("valid from date is required")
	}

	if m.ValidFrom.Before(m.IssueDate) {
		return fmt.Errorf("valid from date cannot be before issue date")
	}

	if m.TermsConditions == "" {
		return fmt.Errorf("terms and conditions are required")
	}

	// Validate date ranges
	if m.ValidUntil != nil && m.ValidUntil.Before(m.ValidFrom) {
		return fmt.Errorf("valid until date cannot be before valid from date")
	}

	// Validate that at least client signature is present for active mandates
	if m.Status == DocumentStatusActive && m.ClientSignature == nil {
		return fmt.Errorf("active mandate must have client signature")
	}

	return nil
}

// IsExpired checks if the mandate has expired.
func (m *Mandate) IsExpired() bool {
	if m.ValidUntil == nil {
		return false // No expiry date set
	}
	return time.Now().After(*m.ValidUntil)
}

// IsValid checks if the mandate is currently valid and active.
func (m *Mandate) IsValid() bool {
	return m.Status == DocumentStatusActive &&
		!m.IsExpired() &&
		m.ClientSignature != nil &&
		time.Now().After(m.ValidFrom)
}

// CanBeSigned checks if the mandate can be signed by the specified role.
func (m *Mandate) CanBeSigned(role string) bool {
	// Cannot sign if already cancelled or archived
	if m.Status == DocumentStatusCancelled || m.Status == DocumentStatusArchived {
		return false
	}

	// Client can always sign (required for validity)
	if role == "client" {
		return m.ClientSignature == nil
	}

	// Investigator can sign (optional but recommended)
	if role == "investigator" {
		return m.InvestigatorSignature == nil
	}

	return false
}

// AddClientSignature adds the client signature to the mandate.
func (m *Mandate) AddClientSignature(signature Signature) error {
	if !m.CanBeSigned("client") {
		return fmt.Errorf("mandate cannot be signed by client")
	}

	m.ClientSignature = &signature

	// If this is the first signature, move to signed status
	if m.Status == DocumentStatusDraft || m.Status == DocumentStatusSent {
		m.SetStatus(DocumentStatusSigned)
	}

	m.UpdateTimestamp()
	return nil
}

// AddInvestigatorSignature adds the investigator signature to the mandate.
func (m *Mandate) AddInvestigatorSignature(signature Signature) error {
	if !m.CanBeSigned("investigator") {
		return fmt.Errorf("mandate cannot be signed by investigator")
	}

	m.InvestigatorSignature = &signature
	m.UpdateTimestamp()
	return nil
}

// Activate marks the mandate as active if it has required signatures.
func (m *Mandate) Activate() error {
	if m.ClientSignature == nil {
		return fmt.Errorf("mandate must have client signature to be activated")
	}

	if m.IsExpired() {
		return fmt.Errorf("cannot activate expired mandate")
	}

	if time.Now().Before(m.ValidFrom) {
		return fmt.Errorf("cannot activate mandate before valid from date")
	}

	m.SetStatus(DocumentStatusActive)
	return nil
}

// LinkToEstimate links this mandate to an estimate.
func (m *Mandate) LinkToEstimate(estimateID uuid.UUID) error {
	if m.LinkedEstimateID != nil {
		return fmt.Errorf("mandate is already linked to an estimate")
	}

	m.LinkedEstimateID = &estimateID
	m.UpdateTimestamp()
	return nil
}

// NewMandate creates a new mandate with the given details.
func NewMandate(caseID, clientID uuid.UUID, mandateNumber, scopeOfWork, termsConditions string, validFrom time.Time) *Mandate {
	mandate := &Mandate{
		DocumentBase:    NewDocumentBase(caseID, clientID),
		MandateNumber:   mandateNumber,
		IssueDate:       time.Now(),
		ScopeOfWork:     scopeOfWork,
		ValidFrom:       validFrom,
		TermsConditions: termsConditions,
	}

	return mandate
}

// SetValidityPeriod sets the validity period for the mandate.
func (m *Mandate) SetValidityPeriod(validFrom time.Time, validUntil *time.Time) error {
	if validFrom.Before(m.IssueDate) {
		return fmt.Errorf("valid from date cannot be before issue date")
	}

	if validUntil != nil && validUntil.Before(validFrom) {
		return fmt.Errorf("valid until date cannot be before valid from date")
	}

	m.ValidFrom = validFrom
	m.ValidUntil = validUntil
	m.UpdateTimestamp()
	return nil
}

// ExtendValidity extends the mandate's validity period.
func (m *Mandate) ExtendValidity(newValidUntil time.Time) error {
	if m.Status == DocumentStatusCancelled || m.Status == DocumentStatusArchived {
		return fmt.Errorf("cannot extend cancelled or archived mandate")
	}

	if newValidUntil.Before(m.ValidFrom) {
		return fmt.Errorf("new valid until date cannot be before valid from date")
	}

	m.ValidUntil = &newValidUntil
	m.UpdateTimestamp()
	return nil
}

