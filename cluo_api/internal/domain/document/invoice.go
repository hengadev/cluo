package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Invoice represents a billing document for services rendered.
type Invoice struct {
	DocumentBase

	InvoiceNumber    string        `json:"invoice_number" db:"invoice_number"`
	IssueDate        time.Time     `json:"issue_date" db:"issue_date"`
	DueDate          time.Time     `json:"due_date" db:"due_date"`
	LineItems        []InvoiceItem `json:"line_items" db:"line_items"`
	TotalAmount      float64       `json:"total_amount" db:"total_amount"`
	TaxRate          float64       `json:"tax_rate" db:"tax_rate"`
	TaxAmount        float64       `json:"tax_amount" db:"tax_amount"`
	Notes            *string       `json:"notes,omitempty" db:"notes"`
	PaymentStatus    PaymentStatus `json:"payment_status" db:"payment_status"`
	PaidAt           *time.Time    `json:"paid_at,omitempty" db:"paid_at"`
	PaidAmount       *float64      `json:"paid_amount,omitempty" db:"paid_amount"`
	PaymentMethod    *string       `json:"payment_method,omitempty" db:"payment_method"`
	LinkedContractID *uuid.UUID    `json:"linked_contract_id,omitempty" db:"linked_contract_id"`
	Currency         *string       `json:"currency,omitempty" db:"currency"`
	PaymentTerms     *string       `json:"payment_terms,omitempty" db:"payment_terms"`
	LateFee          *float64      `json:"late_fee,omitempty" db:"late_fee"`
	LateFeeRate      *float64      `json:"late_fee_rate,omitempty" db:"late_fee_rate"`
}

// InvoiceItem represents a single line item in an invoice.
type InvoiceItem struct {
	Description string  `json:"description" db:"description"`
	Quantity    float64 `json:"quantity" db:"quantity"`
	UnitPrice   float64 `json:"unit_price" db:"unit_price"`
	Subtotal    float64 `json:"subtotal" db:"subtotal"` // Quantity * UnitPrice (persisted or computed)
}

// PaymentStatus represents the payment status of an invoice.
type PaymentStatus string

const (
	PaymentStatusUnpaid        PaymentStatus = "unpaid"
	PaymentStatusPaid          PaymentStatus = "paid"
	PaymentStatusPartiallyPaid PaymentStatus = "partially_paid"
	PaymentStatusOverdue       PaymentStatus = "overdue"
	PaymentStatusRefunded      PaymentStatus = "refunded"
	PaymentStatusVoid          PaymentStatus = "void"
)

// GetType returns the document type.
func (i *Invoice) GetType() DocumentType {
	return DocumentTypeInvoice
}

// Validate performs validation on the invoice.
// TODO: Implement comprehensive invoice validation:
// - InvoiceNumber: required, unique within system, follows format INV-YYYY-NNN
// - IssueDate: required, cannot be in future, reasonable business constraints
// - DueDate: required, must be after IssueDate, reasonable payment terms
// - LineItems: required, at least one item, no empty descriptions
// - LineItems[].Description: required, min length 5, max length 500
// - LineItems[].Quantity: required, positive number, reasonable upper limit
// - LineItems[].UnitPrice: required, positive number, reasonable upper limit
// - LineItems[].Subtotal: should equal Quantity * UnitPrice (allow small rounding differences)
// - TotalAmount: must equal sum of all line item subtotals before tax
// - TaxRate: must be between 0 and 1 (or 0-100 depending on representation)
// - TaxAmount: must equal TotalAmount * TaxRate
// - Notes: optional, reasonable length limit if provided
// - PaymentStatus: required, must be valid enum value
// - PaidAt/ PaidAmount: must be set together when status is paid/partially_paid
// - PaymentMethod: required if paid/partially_paid, must be valid payment method
// - LinkedContractID: optional, if provided must reference existing contract for same case
// - Currency: optional, valid ISO currency code if provided
// - PaymentTerms: optional, standard payment terms if provided
// - LateFee/ LateFeeRate: optional, must be positive if provided
// - Business rules: cannot pay more than total amount + tax + fees
// - Business rules: cannot modify invoice after it's been paid
// - Business rules: payment status transitions must be validated
// - Business rules: LinkedContractID must belong to same case and client
func (i *Invoice) Validate() error {
	// TODO: Add comprehensive validation implementation
	if i.InvoiceNumber == "" {
		return fmt.Errorf("invoice number is required")
	}

	if i.IssueDate.IsZero() {
		return fmt.Errorf("issue date is required")
	}

	if i.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}

	if i.DueDate.Before(i.IssueDate) {
		return fmt.Errorf("due date cannot be before issue date")
	}

	if len(i.LineItems) == 0 {
		return fmt.Errorf("at least one line item is required")
	}

	// Calculate expected totals and validate line items
	var subtotalTotal float64
	for idx, item := range i.LineItems {
		if item.Description == "" {
			return fmt.Errorf("line item %d description is required", idx+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("line item %d quantity must be positive", idx+1)
		}
		if item.UnitPrice <= 0 {
			return fmt.Errorf("line item %d unit price must be positive", idx+1)
		}

		expectedSubtotal := item.Quantity * item.UnitPrice
		// Allow small rounding differences
		if abs(item.Subtotal-expectedSubtotal) > 0.01 {
			return fmt.Errorf("line item %d subtotal mismatch: expected %.2f, got %.2f",
				idx+1, expectedSubtotal, item.Subtotal)
		}

		subtotalTotal += item.Subtotal
	}

	// Validate totals
	expectedTaxAmount := subtotalTotal * i.TaxRate
	expectedTotalAmount := subtotalTotal + expectedTaxAmount

	// Allow small rounding differences
	if abs(i.TotalAmount-subtotalTotal) > 0.01 {
		return fmt.Errorf("total amount mismatch: expected %.2f, got %.2f",
			subtotalTotal, i.TotalAmount)
	}

	if abs(i.TaxAmount-expectedTaxAmount) > 0.01 {
		return fmt.Errorf("tax amount mismatch: expected %.2f, got %.2f",
			expectedTaxAmount, i.TaxAmount)
	}

	// Validate payment status consistency
	if !i.PaymentStatus.IsValid() {
		return fmt.Errorf("invalid payment status: %s", i.PaymentStatus)
	}

	// Validate payment fields consistency
	if (i.PaidAt != nil) != (i.PaidAmount != nil) {
		return fmt.Errorf("paid at and paid amount must be set together")
	}

	return nil
}

// IsOverdue checks if the invoice is overdue.
func (i *Invoice) IsOverdue() bool {
	return i.PaymentStatus == PaymentStatusUnpaid && time.Now().After(i.DueDate)
}

// IsPaid checks if the invoice is fully paid.
func (i *Invoice) IsPaid() bool {
	return i.PaymentStatus == PaymentStatusPaid
}

// GetOutstandingBalance returns the remaining amount to be paid.
func (i *Invoice) GetOutstandingBalance() float64 {
	totalDue := i.TotalAmount + i.TaxAmount

	if i.LateFee != nil {
		totalDue += *i.LateFee
	}

	if i.PaidAmount != nil {
		return totalDue - *i.PaidAmount
	}

	return totalDue
}

// AddPayment records a payment for the invoice.
func (i *Invoice) AddPayment(amount float64, method string) error {
	if amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}

	if i.PaymentStatus == PaymentStatusPaid {
		return fmt.Errorf("invoice is already fully paid")
	}

	if i.PaymentStatus == PaymentStatusVoid {
		return fmt.Errorf("cannot make payment on void invoice")
	}

	outstanding := i.GetOutstandingBalance()
	if amount > outstanding {
		return fmt.Errorf("payment amount %.2f exceeds outstanding balance %.2f", amount, outstanding)
	}

	now := time.Now()
	if i.PaidAmount == nil {
		i.PaidAmount = &amount
	} else {
		*i.PaidAmount += amount
	}

	i.PaidAt = &now
	i.PaymentMethod = &method

	// Update payment status
	if *i.PaidAmount >= outstanding {
		i.PaymentStatus = PaymentStatusPaid
	} else {
		i.PaymentStatus = PaymentStatusPartiallyPaid
	}

	i.UpdateTimestamp()
	return nil
}

// Void marks the invoice as void (cancelled).
func (i *Invoice) Void() error {
	if i.PaymentStatus == PaymentStatusPaid {
		return fmt.Errorf("cannot void paid invoice")
	}

	if i.PaymentStatus == PaymentStatusRefunded {
		return fmt.Errorf("cannot void refunded invoice")
	}

	i.PaymentStatus = PaymentStatusVoid
	i.SetStatus(DocumentStatusCancelled)
	return nil
}

// LinkToContract links this invoice to a contract.
func (i *Invoice) LinkToContract(contractID uuid.UUID) error {
	if i.LinkedContractID != nil {
		return fmt.Errorf("invoice is already linked to a contract")
	}

	i.LinkedContractID = &contractID
	i.UpdateTimestamp()
	return nil
}

// CalculateTotals recalculates all monetary amounts.
func (i *Invoice) CalculateTotals() {
	var subtotalTotal float64
	for idx := range i.LineItems {
		i.LineItems[idx].Subtotal = i.LineItems[idx].Quantity * i.LineItems[idx].UnitPrice
		subtotalTotal += i.LineItems[idx].Subtotal
	}

	i.TotalAmount = subtotalTotal
	i.TaxAmount = subtotalTotal * i.TaxRate

	// Apply late fee if applicable
	if i.IsOverdue() && i.LateFeeRate != nil {
		lateFee := i.GetOutstandingBalance() * *i.LateFeeRate
		i.LateFee = &lateFee
	}

	i.UpdateTimestamp()
}

// NewInvoice creates a new invoice with the given details.
func NewInvoice(caseID, clientID uuid.UUID, invoiceNumber string, lineItems []InvoiceItem, taxRate float64, dueDate time.Time) *Invoice {
	invoice := &Invoice{
		DocumentBase:  NewDocumentBase(caseID, clientID),
		InvoiceNumber: invoiceNumber,
		IssueDate:     time.Now(),
		DueDate:       dueDate,
		LineItems:     lineItems,
		TaxRate:       taxRate,
		PaymentStatus: PaymentStatusUnpaid,
		Currency:      stringPtr("USD"), // Default currency
	}

	invoice.CalculateTotals()
	return invoice
}

// AddLineItem adds a new line item to the invoice.
func (i *Invoice) AddLineItem(description string, quantity, unitPrice float64) error {
	if i.PaymentStatus != PaymentStatusUnpaid {
		return fmt.Errorf("cannot modify invoice that has payments")
	}

	item := InvoiceItem{
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    quantity * unitPrice,
	}
	i.LineItems = append(i.LineItems, item)
	i.CalculateTotals()
	return nil
}

// UpdateLineItem updates an existing line item.
func (i *Invoice) UpdateLineItem(index int, description string, quantity, unitPrice float64) error {
	if index < 0 || index >= len(i.LineItems) {
		return fmt.Errorf("line item index out of range")
	}

	if i.PaymentStatus != PaymentStatusUnpaid {
		return fmt.Errorf("cannot modify invoice that has payments")
	}

	i.LineItems[index].Description = description
	i.LineItems[index].Quantity = quantity
	i.LineItems[index].UnitPrice = unitPrice
	i.LineItems[index].Subtotal = quantity * unitPrice
	i.CalculateTotals()
	return nil
}

// RemoveLineItem removes a line item from the invoice.
func (i *Invoice) RemoveLineItem(index int) error {
	if index < 0 || index >= len(i.LineItems) {
		return fmt.Errorf("line item index out of range")
	}

	if i.PaymentStatus != PaymentStatusUnpaid {
		return fmt.Errorf("cannot modify invoice that has payments")
	}

	if len(i.LineItems) <= 1 {
		return fmt.Errorf("invoice must have at least one line item")
	}

	i.LineItems = append(i.LineItems[:index], i.LineItems[index+1:]...)
	i.CalculateTotals()
	return nil
}

// IsValidPaymentStatus checks if the payment status is valid.
func (ps PaymentStatus) IsValid() bool {
	switch ps {
	case PaymentStatusUnpaid, PaymentStatusPaid, PaymentStatusPartiallyPaid,
		PaymentStatusOverdue, PaymentStatusRefunded, PaymentStatusVoid:
		return true
	default:
		return false
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
