package document

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Generic document routes
	router.HandleFunc("GET /documents", RequireAdministrator(mw.EnableCORS(h.ListDocuments)))
	router.HandleFunc("POST /documents", RequireAdministrator(mw.EnableCORS(h.CreateDocument)))
	router.HandleFunc("GET /documents/{id}/{type}", RequireAdministrator(mw.EnableCORS(h.GetDocument)))
	router.HandleFunc("PATCH /documents/{id}/{type}", RequireAdministrator(mw.EnableCORS(h.UpdateDocument)))
	router.HandleFunc("DELETE /documents/{id}/{type}", RequireAdministrator(mw.EnableCORS(h.DeleteDocument)))
	router.HandleFunc("POST /documents/{id}/{type}/send", RequireAdministrator(mw.EnableCORS(h.SendDocument)))
	router.HandleFunc("POST /documents/{id}/{type}/sign", RequireAdministrator(mw.EnableCORS(h.SignDocument)))
	router.HandleFunc("POST /documents/{id}/{type}/archive", RequireAdministrator(mw.EnableCORS(h.ArchiveDocument)))
	router.HandleFunc("GET /documents/{id}/{type}/history", RequireAdministrator(mw.EnableCORS(h.GetDocumentHistory)))
	router.HandleFunc("GET /documents/workflow/{caseId}", RequireAdministrator(mw.EnableCORS(h.GetDocumentWorkflow)))

	// Generic document lifecycle action routes (type inferred from stored document)
	router.HandleFunc("POST /documents/{id}/accept", RequireAdministrator(mw.EnableCORS(h.AcceptDocument)))
	router.HandleFunc("POST /documents/{id}/activate", RequireAdministrator(mw.EnableCORS(h.ActivateDocument)))
	router.HandleFunc("POST /documents/{id}/pay", RequireAdministrator(mw.EnableCORS(h.PayDocument)))
	router.HandleFunc("POST /documents/{id}/void", RequireAdministrator(mw.EnableCORS(h.VoidDocument)))

	// Estimate-specific routes
	router.HandleFunc("POST /estimates", RequireAdministrator(mw.EnableCORS(h.CreateEstimate)))
	router.HandleFunc("PATCH /estimates/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateEstimate)))
	router.HandleFunc("POST /estimates/{id}/accept", RequireAdministrator(mw.EnableCORS(h.AcceptEstimate)))

	// Mandate-specific routes
	router.HandleFunc("POST /mandates", RequireAdministrator(mw.EnableCORS(h.CreateMandate)))
	router.HandleFunc("POST /mandates/{id}/sign", RequireAdministrator(mw.EnableCORS(h.SignMandate)))
	router.HandleFunc("POST /mandates/{id}/activate", RequireAdministrator(mw.EnableCORS(h.ActivateMandate)))
	router.HandleFunc("POST /mandates/{id}/create-contract", RequireAdministrator(mw.EnableCORS(h.CreateContractFromMandate)))

	// Contract-specific routes
	router.HandleFunc("POST /contracts", RequireAdministrator(mw.EnableCORS(h.CreateContract)))
	router.HandleFunc("POST /contracts/{id}/sign", RequireAdministrator(mw.EnableCORS(h.SignContract)))
	router.HandleFunc("POST /contracts/{id}/activate", RequireAdministrator(mw.EnableCORS(h.ActivateContract)))
	router.HandleFunc("POST /contracts/{id}/create-invoice", RequireAdministrator(mw.EnableCORS(h.CreateInvoiceFromContract)))

	// Invoice-specific routes
	router.HandleFunc("POST /invoices", RequireAdministrator(mw.EnableCORS(h.CreateInvoice)))
	router.HandleFunc("GET /invoices/overdue", RequireAdministrator(mw.EnableCORS(h.GetOverdueInvoices)))
	router.HandleFunc("POST /invoices/{id}/pay", RequireAdministrator(mw.EnableCORS(h.ProcessPayment)))
	router.HandleFunc("POST /invoices/{id}/void", RequireAdministrator(mw.EnableCORS(h.VoidInvoice)))
}
