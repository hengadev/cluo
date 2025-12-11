package document

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeInvoice  DocumentType = "invoice"
	DocumentTypeMandate  DocumentType = "mandate"
	DocumentTypeEstimate DocumentType = "estimate"
	DocumentTypeContract DocumentType = "contract"
	DocumentTypeReport   DocumentType = "report"
	DocumentTypeOther    DocumentType = "other"
)
