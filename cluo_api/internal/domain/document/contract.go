package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/validation"
)

// Contract represents a formal agreement between parties for investigative services.
type Contract struct {
	DocumentBase

	ContractNumber    string      `encx:"encrypt" json:"contract_number" db:"contract_number_encrypted"`
	StartDate         time.Time   `json:"start_date" db:"start_date"`
	EndDate           *time.Time  `json:"end_date,omitempty" db:"end_date"`
	ScopeOfServices   string      `encx:"encrypt" json:"scope_of_services" db:"scope_of_services_encrypted"`
	PaymentTerms      string      `encx:"encrypt" json:"payment_terms" db:"payment_terms_encrypted"`
	Confidentiality   string      `encx:"encrypt" json:"confidentiality" db:"confidentiality_encrypted"`
	TerminationClause string      `encx:"encrypt" json:"termination_clause" db:"termination_clause_encrypted"`
	Signatures        []Signature `encx:"encrypt" json:"signatures" db:"signatures_encrypted"`
	LinkedMandateID   *uuid.UUID  `json:"linked_mandate_id,omitempty" db:"linked_mandate_id"`
	ContractValue     *float64    `encx:"encrypt" json:"contract_value,omitempty" db:"contract_value_encrypted"`
	Currency          *string     `json:"currency,omitempty" db:"currency"`
	RenewalTerms      *string     `encx:"encrypt" json:"renewal_terms,omitempty" db:"renewal_terms_encrypted"`
	GoverningLaw      *string     `json:"governing_law,omitempty" db:"governing_law"`
}

// GetType returns the document type.
func (c *Contract) GetType() DocumentType {
	return DocumentTypeContract
}

// Validate performs validation on the contract.
// TODO: Implement comprehensive contract validation:
// - ContractNumber: required, unique within system, follows format CNT-YYYY-NNN
// - StartDate: required, cannot be in the past beyond reasonable limits, not too far in future
// - EndDate: optional, must be after StartDate if provided, reasonable duration limits
// - ScopeOfServices: required, detailed description of services, min/max length constraints
// - PaymentTerms: required, clear payment schedule and terms, reasonable business constraints
// - Confidentiality: required, standard confidentiality clauses
// - TerminationClause: required, termination conditions and notice periods
// - Signatures: at least one signature required for validity
// - LinkedMandateID: optional, if provided must reference existing mandate for same case
// - ContractValue: optional, must be positive if provided
// - Currency: optional, valid ISO currency code if provided
// - RenewalTerms: optional, reasonable length if provided
// - GoverningLaw: optional, valid jurisdiction if provided
// - Business rules: cannot be active without required signatures
// - Business rules: contract duration must be reasonable (e.g., max 5 years)
// - Business rules: all signatories must be from same case/client
// - Business rules: cannot modify contract after signatures are added
// - Business rules: LinkedMandateID must belong to same case and client
func (c *Contract) Validate() error {
	// TODO: Add comprehensive validation implementation
	if c.ContractNumber == "" {
		return fmt.Errorf("contract number is required")
	}

	if c.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}

	if c.ScopeOfServices == "" {
		return fmt.Errorf("scope of services is required")
	}

	if len(c.ScopeOfServices) < 50 {
		return fmt.Errorf("scope of services must be at least 50 characters")
	}

	if c.PaymentTerms == "" {
		return fmt.Errorf("payment terms are required")
	}

	if c.Confidentiality == "" {
		return fmt.Errorf("confidentiality clause is required")
	}

	if c.TerminationClause == "" {
		return fmt.Errorf("termination clause is required")
	}

	// Validate date ranges
	if c.EndDate != nil && c.EndDate.Before(c.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	// Validate contract value
	if c.ContractValue != nil && *c.ContractValue <= 0 {
		return fmt.Errorf("contract value must be positive")
	}

	// Validate currency code
	if c.Currency != nil && *c.Currency != "" {
		if err := validation.ValidateCurrency(*c.Currency); err != nil {
			return err
		}
	}

	// Validate signatures
	if len(c.Signatures) == 0 && c.Status == DocumentStatusActive {
		return fmt.Errorf("active contract must have at least one signature")
	}

	return nil
}

// IsExpired checks if the contract has expired.
func (c *Contract) IsExpired() bool {
	if c.EndDate == nil {
		return false // No expiry date set
	}
	return time.Now().After(*c.EndDate)
}

// IsActive checks if the contract is currently active and valid.
func (c *Contract) IsActive() bool {
	return c.Status == DocumentStatusActive &&
		!c.IsExpired() &&
		time.Now().After(c.StartDate) &&
		len(c.Signatures) > 0
}

// IsWithinValidityPeriod checks if a given date is within the contract's validity period.
func (c *Contract) IsWithinValidityPeriod(date time.Time) bool {
	if date.Before(c.StartDate) {
		return false
	}
	if c.EndDate != nil && date.After(*c.EndDate) {
		return false
	}
	return true
}

// CanBeSigned checks if the contract can be signed by the specified person.
func (c *Contract) CanBeSigned(signerName, role string) bool {
	// Cannot sign if already cancelled or archived
	if c.Status == DocumentStatusCancelled || c.Status == DocumentStatusArchived {
		return false
	}

	// Check if this person has already signed
	for _, sig := range c.Signatures {
		if sig.Name == signerName && sig.Role == role {
			return false // Already signed
		}
	}

	return true
}

// AddSignature adds a signature to the contract.
func (c *Contract) AddSignature(signature Signature) error {
	if !c.CanBeSigned(signature.Name, signature.Role) {
		return fmt.Errorf("contract cannot be signed by %s as %s", signature.Name, signature.Role)
	}

	c.Signatures = append(c.Signatures, signature)

	// If this is the first signature, move to signed status
	if c.Status == DocumentStatusDraft || c.Status == DocumentStatusSent {
		c.SetStatus(DocumentStatusSigned)
	}

	c.UpdateTimestamp()
	return nil
}

// Activate marks the contract as active if it has required signatures and valid dates.
func (c *Contract) Activate() error {
	if len(c.Signatures) == 0 {
		return fmt.Errorf("contract must have at least one signature to be activated")
	}

	if time.Now().Before(c.StartDate) {
		return fmt.Errorf("cannot activate contract before start date")
	}

	if c.IsExpired() {
		return fmt.Errorf("cannot activate expired contract")
	}

	c.SetStatus(DocumentStatusActive)
	return nil
}

// LinkToMandate links this contract to a mandate.
func (c *Contract) LinkToMandate(mandateID uuid.UUID) error {
	if c.LinkedMandateID != nil {
		return fmt.Errorf("contract is already linked to a mandate")
	}

	c.LinkedMandateID = &mandateID
	c.UpdateTimestamp()
	return nil
}

// NewContract creates a new contract with the given details.
func NewContract(caseID, clientID uuid.UUID, contractNumber, scopeOfServices, paymentTerms, confidentiality, terminationClause string, startDate time.Time) *Contract {
	contract := &Contract{
		DocumentBase:      NewDocumentBase(caseID, clientID),
		ContractNumber:    contractNumber,
		StartDate:         startDate,
		ScopeOfServices:   scopeOfServices,
		PaymentTerms:      paymentTerms,
		Confidentiality:   confidentiality,
		TerminationClause: terminationClause,
		Signatures:        []Signature{},
	}

	return contract
}

// SetEndDate sets the end date for the contract.
func (c *Contract) SetEndDate(endDate time.Time) error {
	if endDate.Before(c.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	c.EndDate = &endDate
	c.UpdateTimestamp()
	return nil
}

// ExtendContract extends the contract's end date.
func (c *Contract) ExtendContract(newEndDate time.Time) error {
	if c.Status == DocumentStatusCancelled || c.Status == DocumentStatusArchived {
		return fmt.Errorf("cannot extend cancelled or archived contract")
	}

	if newEndDate.Before(c.StartDate) {
		return fmt.Errorf("new end date cannot be before start date")
	}

	if c.EndDate != nil && newEndDate.Before(*c.EndDate) {
		return fmt.Errorf("new end date cannot be before current end date")
	}

	c.EndDate = &newEndDate
	c.UpdateTimestamp()
	return nil
}

// SetContractValue sets the contract value and currency.
func (c *Contract) SetContractValue(value float64, currency string) error {
	if value <= 0 {
		return fmt.Errorf("contract value must be positive")
	}

	if currency == "" {
		currency = "USD" // Default currency
	}

	// TODO: Validate currency code format (ISO 4217)

	c.ContractValue = &value
	c.Currency = &currency
	c.UpdateTimestamp()
	return nil
}

// GetRemainingDays returns the number of days remaining until contract expiration.
func (c *Contract) GetRemainingDays() int {
	if c.EndDate == nil {
		return -1 // No end date
	}

	now := time.Now()
	if now.After(*c.EndDate) {
		return 0 // Already expired
	}

	duration := c.EndDate.Sub(now)
	return int(duration.Hours() / 24)
}

// HasSignature checks if the contract has a signature from a specific person/role.
func (c *Contract) HasSignature(signerName, role string) bool {
	for _, sig := range c.Signatures {
		if sig.Name == signerName && sig.Role == role {
			return true
		}
	}
	return false
}
