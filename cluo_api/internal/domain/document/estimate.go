package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Estimate represents a price quotation for services to be provided.
type Estimate struct {
	DocumentBase

	EstimateNumber string         `encx:"encrypt" json:"estimate_number" db:"estimate_number_encrypted"`
	IssueDate      time.Time      `json:"issue_date" db:"issue_date"`
	ValidUntil     *time.Time     `json:"valid_until,omitempty" db:"valid_until"`
	LineItems      []EstimateItem `encx:"encrypt" json:"line_items" db:"line_items_encrypted"`
	EstimatedTotal float64        `encx:"encrypt" json:"estimated_total" db:"estimated_total_encrypted"`
	Notes          *string        `encx:"encrypt" json:"notes,omitempty" db:"notes_encrypted"`
	Accepted       bool           `json:"accepted" db:"accepted"`
	AcceptedAt     *time.Time     `json:"accepted_at,omitempty" db:"accepted_at"`
	AcceptedBy     *uuid.UUID     `json:"accepted_by,omitempty" db:"accepted_by"`
}

// EstimateItem represents a single line item in an estimate.
type EstimateItem struct {
	Description string  `json:"description" db:"description"`
	Quantity    float64 `json:"quantity" db:"quantity"`
	UnitPrice   float64 `json:"unit_price" db:"unit_price"`
	Subtotal    float64 `json:"subtotal" db:"subtotal"` // Quantity * UnitPrice (persisted or computed)
}

// GetType returns the document type.
func (e *Estimate) GetType() DocumentType {
	return DocumentTypeEstimate
}

// Validate performs validation on the estimate.
// TODO: Implement comprehensive estimate validation:
// - EstimateNumber: required, unique within system, follows format EST-YYYY-NNN
// - IssueDate: required, cannot be in the future, cannot be too far in past
// - ValidUntil: optional, must be after IssueDate if provided
// - LineItems: required, at least one item, no empty descriptions
// - LineItems[].Description: required, min length 5, max length 500
// - LineItems[].Quantity: required, positive number, reasonable upper limit (e.g., 9999)
// - LineItems[].UnitPrice: required, positive number, reasonable upper limit
// - LineItems[].Subtotal: should equal Quantity * UnitPrice (allow small rounding differences)
// - EstimatedTotal: must equal sum of all line item subtotals
// - Notes: optional, reasonable length limit if provided
// - Accepted: false initially, only settable through AcceptEstimate operation
// - AcceptedAt/ AcceptedBy: must be set together, AcceptedAt must be after IssueDate
// - Business rules: cannot accept if expired (ValidUntil < now)
// - Business rules: cannot modify estimate after acceptance
// - Business rules: cannot delete estimate if linked to mandate
func (e *Estimate) Validate() error {
	// TODO: Add comprehensive validation implementation
	if e.EstimateNumber == "" {
		return fmt.Errorf("estimate number is required")
	}

	if e.IssueDate.IsZero() {
		return fmt.Errorf("issue date is required")
	}

	if len(e.LineItems) == 0 {
		return fmt.Errorf("at least one line item is required")
	}

	// Calculate expected total and validate line items
	var calculatedTotal float64
	for i, item := range e.LineItems {
		if item.Description == "" {
			return fmt.Errorf("line item %d description is required", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("line item %d quantity must be positive", i+1)
		}
		if item.UnitPrice <= 0 {
			return fmt.Errorf("line item %d unit price must be positive", i+1)
		}

		expectedSubtotal := item.Quantity * item.UnitPrice
		// Allow small rounding differences
		if abs(item.Subtotal-expectedSubtotal) > 0.01 {
			return fmt.Errorf("line item %d subtotal mismatch: expected %.2f, got %.2f",
				i+1, expectedSubtotal, item.Subtotal)
		}

		calculatedTotal += item.Subtotal
	}

	// Allow small rounding differences for total
	if abs(e.EstimatedTotal-calculatedTotal) > 0.01 {
		return fmt.Errorf("estimated total mismatch: expected %.2f, got %.2f",
			calculatedTotal, e.EstimatedTotal)
	}

	// Validate dates
	if e.ValidUntil != nil && e.ValidUntil.Before(e.IssueDate) {
		return fmt.Errorf("valid until date cannot be before issue date")
	}

	return nil
}

// CalculateTotal recalculates the total from line items.
func (e *Estimate) CalculateTotal() {
	var total float64
	for i := range e.LineItems {
		e.LineItems[i].Subtotal = e.LineItems[i].Quantity * e.LineItems[i].UnitPrice
		total += e.LineItems[i].Subtotal
	}
	e.EstimatedTotal = total
	e.UpdateTimestamp()
}

// IsExpired checks if the estimate has expired.
func (e *Estimate) IsExpired() bool {
	if e.ValidUntil == nil {
		return false // No expiry date set
	}
	return time.Now().After(*e.ValidUntil)
}

// CanBeAccepted checks if the estimate can be accepted.
func (e *Estimate) CanBeAccepted() bool {
	return !e.Accepted && !e.IsExpired() && e.Status == DocumentStatusSent
}

// Accept marks the estimate as accepted.
func (e *Estimate) Accept(acceptedBy uuid.UUID) error {
	if !e.CanBeAccepted() {
		return fmt.Errorf("estimate cannot be accepted")
	}

	now := time.Now()
	e.Accepted = true
	e.AcceptedAt = &now
	e.AcceptedBy = &acceptedBy
	e.SetStatus(DocumentStatusSigned)
	return nil
}

// NewEstimate creates a new estimate with the given details.
func NewEstimate(caseID, clientID uuid.UUID, estimateNumber string, lineItems []EstimateItem) *Estimate {
	estimate := &Estimate{
		DocumentBase:   NewDocumentBase(caseID, clientID),
		EstimateNumber: estimateNumber,
		IssueDate:      time.Now(),
		LineItems:      lineItems,
		Accepted:       false,
	}

	estimate.CalculateTotal()
	return estimate
}

// AddLineItem adds a new line item to the estimate.
func (e *Estimate) AddLineItem(description string, quantity, unitPrice float64) {
	item := EstimateItem{
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    quantity * unitPrice,
	}
	e.LineItems = append(e.LineItems, item)
	e.CalculateTotal()
}

// UpdateLineItem updates an existing line item.
func (e *Estimate) UpdateLineItem(index int, description string, quantity, unitPrice float64) error {
	if index < 0 || index >= len(e.LineItems) {
		return fmt.Errorf("line item index out of range")
	}

	if e.Accepted {
		return fmt.Errorf("cannot modify accepted estimate")
	}

	e.LineItems[index].Description = description
	e.LineItems[index].Quantity = quantity
	e.LineItems[index].UnitPrice = unitPrice
	e.LineItems[index].Subtotal = quantity * unitPrice
	e.CalculateTotal()
	return nil
}

// RemoveLineItem removes a line item from the estimate.
func (e *Estimate) RemoveLineItem(index int) error {
	if index < 0 || index >= len(e.LineItems) {
		return fmt.Errorf("line item index out of range")
	}

	if e.Accepted {
		return fmt.Errorf("cannot modify accepted estimate")
	}

	if len(e.LineItems) <= 1 {
		return fmt.Errorf("estimate must have at least one line item")
	}

	e.LineItems = append(e.LineItems[:index], e.LineItems[index+1:]...)
	e.CalculateTotal()
	return nil
}

// Helper function for floating point comparison
func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
