package document

import (
	"time"

	"github.com/google/uuid"
)

// EstimateEncx Documentable implementation.
// CaseID/ClientID are encrypted and not available without decryption; GetCaseID returns uuid.Nil.
func (e *EstimateEncx) GetID() uuid.UUID              { return e.ID }
func (e *EstimateEncx) GetType() DocumentType         { return DocumentTypeEstimate }
func (e *EstimateEncx) GetCaseID() uuid.UUID          { return uuid.Nil }
func (e *EstimateEncx) GetStatus() DocumentStatus     { return e.Status }
func (e *EstimateEncx) Validate() error               { return nil }
func (e *EstimateEncx) SetStatus(s DocumentStatus)    { e.Status = s; e.UpdatedAt = time.Now() }
func (e *EstimateEncx) UpdateTimestamp()               { e.UpdatedAt = time.Now() }
func (e *EstimateEncx) SetCaseID(_ uuid.UUID)          {}
func (e *EstimateEncx) SetClientID(_ uuid.UUID)        {}

// MandateEncx Documentable implementation.
func (m *MandateEncx) GetID() uuid.UUID              { return m.ID }
func (m *MandateEncx) GetType() DocumentType         { return DocumentTypeMandate }
func (m *MandateEncx) GetCaseID() uuid.UUID          { return uuid.Nil }
func (m *MandateEncx) GetStatus() DocumentStatus     { return m.Status }
func (m *MandateEncx) Validate() error               { return nil }
func (m *MandateEncx) SetStatus(s DocumentStatus)    { m.Status = s; m.UpdatedAt = time.Now() }
func (m *MandateEncx) UpdateTimestamp()               { m.UpdatedAt = time.Now() }
func (m *MandateEncx) SetCaseID(_ uuid.UUID)          {}
func (m *MandateEncx) SetClientID(_ uuid.UUID)        {}

// ContractEncx Documentable implementation.
func (c *ContractEncx) GetID() uuid.UUID              { return c.ID }
func (c *ContractEncx) GetType() DocumentType         { return DocumentTypeContract }
func (c *ContractEncx) GetCaseID() uuid.UUID          { return uuid.Nil }
func (c *ContractEncx) GetStatus() DocumentStatus     { return c.Status }
func (c *ContractEncx) Validate() error               { return nil }
func (c *ContractEncx) SetStatus(s DocumentStatus)    { c.Status = s; c.UpdatedAt = time.Now() }
func (c *ContractEncx) UpdateTimestamp()               { c.UpdatedAt = time.Now() }
func (c *ContractEncx) SetCaseID(_ uuid.UUID)          {}
func (c *ContractEncx) SetClientID(_ uuid.UUID)        {}

// InvoiceEncx Documentable implementation.
func (i *InvoiceEncx) GetID() uuid.UUID              { return i.ID }
func (i *InvoiceEncx) GetType() DocumentType         { return DocumentTypeInvoice }
func (i *InvoiceEncx) GetCaseID() uuid.UUID          { return uuid.Nil }
func (i *InvoiceEncx) GetStatus() DocumentStatus     { return i.Status }
func (i *InvoiceEncx) Validate() error               { return nil }
func (i *InvoiceEncx) SetStatus(s DocumentStatus)    { i.Status = s; i.UpdatedAt = time.Now() }
func (i *InvoiceEncx) UpdateTimestamp()               { i.UpdatedAt = time.Now() }
func (i *InvoiceEncx) SetCaseID(_ uuid.UUID)          {}
func (i *InvoiceEncx) SetClientID(_ uuid.UUID)        {}
