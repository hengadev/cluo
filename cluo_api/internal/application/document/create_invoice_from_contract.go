package document

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateInvoiceFromContract derives a draft invoice from an active contract.
// All invoice fields are derived from the contract — no body is required from the caller.
func (s *Service) CreateInvoiceFromContract(ctx context.Context, contractID string) (*document.Invoice, error) {
	contractEncx, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "contract")
	}

	contract, err := document.DecryptContractEncx(ctx, s.crypto, contractEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contract", err)
	}

	if !contract.IsActive() {
		return nil, errs.NewInvalidValueErr("contract must be active to create invoice")
	}

	parsedContractID, err := uuid.Parse(contractID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid contract ID format")
	}

	if contract.ContractValue == nil || *contract.ContractValue <= 0 {
		return nil, errs.NewInvalidValueErr("contract must have a positive contract value to derive an invoice")
	}

	// Derive a single line item from the contract scope and value.
	desc := contract.ScopeOfServices
	if len(desc) > 100 {
		desc = desc[:100]
	}
	lineItems := []document.InvoiceItem{{
		Description: desc,
		Quantity:    1,
		UnitPrice:   *contract.ContractValue,
		Subtotal:    *contract.ContractValue,
	}}

	// Derive invoice number from the contract number (CNT-… → FAC-…).
	invoiceNumber := "FAC-" + strings.TrimPrefix(contract.ContractNumber, "CNT-")

	now := time.Now()
	dueDate := now.AddDate(0, 0, 30) // NET 30

	invoice := document.NewInvoice(
		contract.CaseID,
		contract.ClientID,
		invoiceNumber,
		lineItems,
		0, // tax rate defaults to 0%
		dueDate,
	)
	invoice.LinkedContractID = &parsedContractID

	// Inherit currency from the contract, defaulting to EUR.
	if contract.Currency != nil && *contract.Currency != "" {
		invoice.Currency = contract.Currency
	} else {
		eur := "EUR"
		invoice.Currency = &eur
	}

	if contract.PaymentTerms != "" {
		invoice.Notes = &contract.PaymentTerms
	}

	if err := invoice.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr("invoice validation failed: " + err.Error())
	}

	invoice.Status = document.DocumentStatusDraft

	invoiceEncx, err := document.ProcessInvoiceEncx(ctx, s.crypto, invoice)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("invoice", err)
	}

	if err := s.repo.CreateInvoice(ctx, invoiceEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "invoice")
	}

	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr("Created from contract")); err != nil {
		// Non-fatal: versioning failure does not block invoice creation
	}

	return invoice, nil
}
