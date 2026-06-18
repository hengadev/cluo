package document

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateDocument updates an existing document's content. Only draft documents
// can be edited this way — once a document has been sent, its content is
// frozen until a versioning workflow exists.
func (s *Service) UpdateDocument(ctx context.Context, id string, req *document.UpdateDocumentRequest) (document.Documentable, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing (encrypted) document
	encDoc, err := s.repo.GetByID(ctx, id, req.Type)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "document")
	}

	if encDoc.GetStatus() != document.DocumentStatusDraft {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("cannot edit document in status: %s", encDoc.GetStatus()))
	}

	// Decrypt so the patch is applied to plain fields
	decrypted, err := document.DecryptDocumentable(ctx, s.crypto, encDoc)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("document", err)
	}

	// Apply the requested changes on top of the decrypted document, leaving
	// any field absent from req.Data untouched.
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid update data: %s", err.Error()))
	}
	if err := json.Unmarshal(dataBytes, decrypted); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("invalid update data: %s", err.Error()))
	}

	doc, ok := decrypted.(document.Documentable)
	if !ok {
		return nil, fmt.Errorf("document: decrypted document does not implement Documentable")
	}
	doc.UpdateTimestamp()

	// Validate updated document
	if err := doc.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("document validation failed: %s", err.Error()))
	}

	// Re-encrypt and save
	reEncrypted, err := document.EncryptDocumentable(ctx, s.crypto, decrypted)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("document", err)
	}
	if err := s.repo.Update(ctx, reEncrypted); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, doc, &authorID, req.Reason); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return doc, nil
}