package document

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/hengadev/cluo_api/internal/ports"
)

func RegisterRoutes(r chi.Router, service ports.DocumentService) {
	handler := New(service)

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// Base routes for documents
	r.Route("/documents", func(r chi.Router) {
		// Generic document operations
		r.Get("/", handler.ListDocuments)
		r.Post("/", handler.CreateDocument)

		// Document workflow routes
		r.Route("/{id}/{type}", func(r chi.Router) {
			r.Get("/", handler.GetDocument)
			r.Patch("/", handler.UpdateDocument)
			r.Delete("/", handler.DeleteDocument)
			r.Post("/send", handler.SendDocument)
			r.Post("/sign", handler.SignDocument)
			r.Post("/archive", handler.ArchiveDocument)
			r.Get("/history", handler.GetDocumentHistory)
		})

		// Case-specific document workflow
		r.Get("/workflow/{caseId}", handler.GetDocumentWorkflow)
	})

	// Estimate-specific routes
	r.Route("/estimates", func(r chi.Router) {
		r.Post("/", handler.CreateEstimate)

		r.Route("/{id}", func(r chi.Router) {
			r.Patch("/", handler.UpdateEstimate)
			r.Post("/accept", handler.AcceptEstimate)
		})
	})

	// Mandate-specific routes
	r.Route("/mandates", func(r chi.Router) {
		r.Post("/", handler.CreateMandate)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/sign", handler.SignMandate)
			r.Post("/activate", handler.ActivateMandate)
			r.Post("/create-contract", handler.CreateContractFromMandate)
		})
	})

	// Contract-specific routes
	r.Route("/contracts", func(r chi.Router) {
		r.Post("/", handler.CreateContract)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/sign", handler.SignContract)
			r.Post("/activate", handler.ActivateContract)
			r.Post("/create-invoice", handler.CreateInvoiceFromContract)
		})
	})

	// Invoice-specific routes
	r.Route("/invoices", func(r chi.Router) {
		r.Post("/", handler.CreateInvoice)
		r.Get("/overdue", handler.GetOverdueInvoices)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/pay", handler.ProcessPayment)
			r.Post("/void", handler.VoidInvoice)
		})
	})

	// Legacy routes for backward compatibility (typed routes that redirect to generic)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/documents", func(r chi.Router) {
			r.Get("/", handler.ListDocuments)
			r.Post("/", handler.CreateDocument)

			r.Route("/{id}/{type}", func(r chi.Router) {
				r.Get("/", handler.GetDocument)
				r.Patch("/", handler.UpdateDocument)
				r.Delete("/", handler.DeleteDocument)
				r.Post("/send", handler.SendDocument)
				r.Post("/sign", handler.SignDocument)
				r.Post("/archive", handler.ArchiveDocument)
				r.Get("/history", handler.GetDocumentHistory)
			})
		})
	})
}

