# Documents Domain — Implementation Specification

This file describes everything that needs to be implemented for the **documents domain** (Estimate, Mandate, Contract, Invoice).  
It assumes other entities (Client, Investigator, Case) exist but intentionally **excludes their definitions and implementations** — focus is only on documents, their storage, lifecycle and API/operations.

---

## Purpose

Capture the complete technical requirements for implementing document generation, storage, lifecycle, relationships, signing and versioning so an LLM or engineer can implement the backend and supporting APIs.

---

## Table of Contents

1. Base types
2. Document type definitions
3. Document relationships / workflow
4. Storage options (typed tables and unified JSON table)
5. Versioning & audit
6. Operations (business logic)
7. Suggested SQL schema (Postgres)
8. Suggested REST API endpoints
9. Example event flows
10. Notes & implementation hints

---

## 1) Base Types

Use a common base struct embedded in all document-specific structs.

~~~go
// DocumentBase: metadata shared by all documents.
type DocumentBase struct {
    ID              string    // uuid
    CaseID          string    // foreign key to Case
    ClientID        string    // foreign key to Client
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Status          string    // draft, sent, signed, active, archived, cancelled
}
~~~

Signature helper used by Mandate / Contract:

~~~go
type Signature struct {
    Name             string
    Role             string    // e.g., "client", "investigator", "witness"
    SignedAt         time.Time
    SignatureFileURL string    // URL to stored signature (image/PDF) or a digital-signature envelope id
    Method           string    // "e-sign", "wet", "pdf-stamp", "third-party"
}
~~~

---

## 2) Document Type Definitions

Implement typed structs for each document type. Keep fields explicit and typed.

### Estimate

~~~go
type Estimate struct {
    DocumentBase

    EstimateNumber string
    IssueDate      time.Time
    ValidUntil     *time.Time
    LineItems      []EstimateItem
    EstimatedTotal float64
    Notes          *string
    Accepted       bool
    AcceptedAt     *time.Time
}

type EstimateItem struct {
    Description string
    Quantity    float64
    UnitPrice   float64
    Subtotal    float64 // Quantity * UnitPrice (persisted or computed)
}
~~~

### Mandate

~~~go
type Mandate struct {
    DocumentBase

    MandateNumber         string
    IssueDate             time.Time
    ScopeOfWork           string    // clear textual description of permitted actions
    ValidFrom             time.Time
    ValidUntil            *time.Time
    TermsConditions       string
    ClientSignature       *Signature
    InvestigatorSignature *Signature
    LinkedEstimateID      *string
}
~~~

### Contract

~~~go
type Contract struct {
    DocumentBase

    ContractNumber     string
    StartDate          time.Time
    EndDate            *time.Time
    ScopeOfServices    string
    PaymentTerms       string
    Confidentiality    string
    TerminationClause  string
    Signatures         []Signature
    LinkedMandateID    *string
}
~~~

### Invoice

~~~go
type Invoice struct {
    DocumentBase

    InvoiceNumber    string
    IssueDate        time.Time
    DueDate          time.Time
    LineItems        []InvoiceItem
    TotalAmount      float64
    TaxRate          float64
    TaxAmount        float64
    Notes            *string
    PaymentStatus    string // unpaid, paid, partially_paid, overdue, refunded
    PaidAt           *time.Time
    LinkedContractID *string
}

type InvoiceItem struct {
    Description string
    Quantity    float64
    UnitPrice   float64
    Subtotal    float64
}
~~~

---

## 3) Document Relationships / Workflow

Documents typically progress in this order:

Estimate → Mandate → Contract → Invoice


Linking fields:
- `Mandate.LinkedEstimateID` → `Estimate.ID`
- `Contract.LinkedMandateID` → `Mandate.ID`
- `Invoice.LinkedContractID` → `Contract.ID`

These links are optional but recommended for traceability and reconstruction of the negotiation lifecycle.

---

## 4) Storage Options

You can implement either or both patterns:

### Option A — Typed relational tables (recommended for clarity + typed queries)

- `estimates`, `mandates`, `contracts`, `invoices` — each table matches the typed struct.
- Pros: strong typing, easier SQL queries, safer constraints, smaller records.
- Cons: schema migrations when fields change.

### Option B — Unified JSON documents table (flexible)

Single table `documents`:

CREATE TABLE documents (
  id uuid PRIMARY KEY,
  type text NOT NULL, -- 'estimate' | 'mandate' | 'contract' | 'invoice'
  case_id uuid NOT NULL,
  client_id uuid NOT NULL,
  data jsonb NOT NULL, -- serialized typed struct
  status text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX ON documents ((data->>'EstimateNumber'));

- Pros: very flexible, easier to store unknown fields.
- Cons: harder to enforce invariants, more app-side validation.

**Hybrid approach**: keep typed tables for invoices/payments (financial data) and store other docs in `documents` JSON table.

---

## 5) Versioning & Audit

Implement `document_versions` table (or version subcollection) to keep every change.

type DocumentVersion struct {
    DocumentID string
    Version    int
    AuthorID   *string   // who made the change
    Data       json.RawMessage
    CreatedAt  time.Time
    Reason     *string   // optional note for the edit
}

Rules:
- Increment version on any state-changing update (including minor edits).
- Never delete old versions; mark as superseded.
- Provide a `GetDocumentHistory(documentID)` API.

---

## 6) Business Operations (functions to implement)

Implement the following high-level operations with input validation and state transitions:

- `CreateDocument(type, payload) -> Document`  
  - Create a new document in `draft` status.

- `UpdateDocument(documentID, patch) -> Document`  
  - Allowed only in `draft` or `sent` (depending on your policy). Create a new version on update.

- `SendDocument(documentID, sendOptions) -> Document`  
  - Set `Status = sent`, record `sent_at`, optionally create email/notification event and a PDF generation job.

- `SignDocument(documentID, signerInfo) -> Document`  
  - Attach `Signature`, mark document as `signed` or `active`. Validate required signatures (e.g., Mandate requires client signature).

- `AcceptEstimate(estimateID) -> Mandate`  
  - Optionally create a Mandate pre-populated from Estimate fields and link `LinkedEstimateID`.

- `CreateContractFromMandate(mandateID, contractPayload) -> Contract`  
  - Link back `LinkedMandateID`.

- `CreateInvoiceFromContract(contractID, invoicePayload) -> Invoice`  
  - Link back `LinkedContractID`.

- `PayInvoice(invoiceID, paymentInfo) -> Invoice`  
  - Update `PaymentStatus` and record `PaidAt`, payment receipts.

- `ArchiveDocument(documentID) -> Document`  
  - Move to `archived`, restrict further edits.

- `GetDocumentsByCase(caseID) -> []DocumentSummary`  
  - Return all documents sorted by createdAt with type and status.

- `GetDocumentHistory(documentID) -> []DocumentVersion`

Add access control checks in each operation (who can create/update/send/sign/archive).

---

## 7) Suggested SQL Schema (Postgres) — typed tables (minimal)

> Note: these are simplified DDLs to get started.

-- Estimates
CREATE TABLE estimates (
  id uuid PRIMARY KEY,
  case_id uuid NOT NULL,
  client_id uuid NOT NULL,
  estimate_number text NOT NULL UNIQUE,
  issue_date timestamptz NOT NULL,
  valid_until timestamptz,
  line_items jsonb NOT NULL, -- or separate table estimate_items
  estimated_total numeric(12,2) NOT NULL,
  notes text,
  accepted boolean DEFAULT false,
  accepted_at timestamptz,
  status text NOT NULL DEFAULT 'draft',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Mandates
CREATE TABLE mandates (
  id uuid PRIMARY KEY,
  case_id uuid NOT NULL,
  client_id uuid NOT NULL,
  mandate_number text NOT NULL UNIQUE,
  issue_date timestamptz NOT NULL,
  scope_of_work text NOT NULL,
  valid_from timestamptz NOT NULL,
  valid_until timestamptz,
  terms_conditions text,
  client_signature jsonb, -- friendly to store Signature struct
  investigator_signature jsonb,
  linked_estimate_id uuid,
  status text NOT NULL DEFAULT 'draft',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Contracts
CREATE TABLE contracts (
  id uuid PRIMARY KEY,
  case_id uuid NOT NULL,
  client_id uuid NOT NULL,
  contract_number text NOT NULL UNIQUE,
  start_date timestamptz NOT NULL,
  end_date timestamptz,
  scope_of_services text,
  payment_terms text,
  confidentiality text,
  termination_clause text,
  signatures jsonb,
  linked_mandate_id uuid,
  status text NOT NULL DEFAULT 'draft',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Invoices
CREATE TABLE invoices (
  id uuid PRIMARY KEY,
  case_id uuid NOT NULL,
  client_id uuid NOT NULL,
  invoice_number text NOT NULL UNIQUE,
  issue_date timestamptz NOT NULL,
  due_date timestamptz NOT NULL,
  line_items jsonb NOT NULL,
  total_amount numeric(12,2) NOT NULL,
  tax_rate numeric(5,2) DEFAULT 0,
  tax_amount numeric(12,2) DEFAULT 0,
  notes text,
  payment_status text NOT NULL DEFAULT 'unpaid',
  paid_at timestamptz,
  linked_contract_id uuid,
  status text NOT NULL DEFAULT 'draft',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Document versions (generic)
CREATE TABLE document_versions (
  id bigserial PRIMARY KEY,
  document_id uuid NOT NULL,
  doc_type text NOT NULL,
  version int NOT NULL,
  author_id uuid,
  data jsonb NOT NULL,
  reason text,
  created_at timestamptz NOT NULL DEFAULT now()
);

---

## 8) Suggested REST API Endpoints

Design the API using RESTful resources. Replace `:id` with UUID.

### Documents (generic)
- `POST /api/documents`  
  - Body: `{ type: "estimate"|"mandate"|"contract"|"invoice", data: {...} }`  
  - Creates a draft document.

- `GET /api/documents/:id`  
  - Returns typed document data with metadata.

- `PATCH /api/documents/:id`  
  - Update draft; create new version.

- `POST /api/documents/:id/send`  
  - Mark sent and queue email/pdf.

- `POST /api/documents/:id/sign`  
  - Attach signature (body contains signer info & signature URL or e-sign envelope id).

- `GET /api/cases/:caseId/documents`  
  - List documents for a case.

- `GET /api/documents/:id/history`  
  - Return versions.

### Typed convenience endpoints
- `POST /api/estimates/:id/accept` → creates a Mandate draft (linked).
- `POST /api/mandates/:id/create-contract` → create Contract (linked).
- `POST /api/contracts/:id/create-invoice` → create Invoice (linked).
- `POST /api/invoices/:id/pay` → register a payment.

### Admin / utilities
- `POST /api/documents/:id/archive`
- `GET /api/documents?status=sent&type=invoice&from=2025-01-01&to=2025-12-31` (filtering)

---

## 9) Example Event Flows

### Flow: Estimate → Mandate → Contract → Invoice → Payment

1. User creates `Estimate` (status `draft`).
2. User sends Estimate (`status: sent`). Email + PDF generated.
3. Client accepts Estimate (via link or manual). Call `POST /api/estimates/:id/accept`.
4. Server creates `Mandate` prepopulated from `Estimate`. `Mandate.LinkedEstimateID` = estimate.id.
5. Client signs Mandate (`POST /api/mandates/:id/sign`). Mandate `status = signed`.
6. Investigator creates `Contract` referencing Mandate (`Contract.LinkedMandateID`).
7. Parties sign Contract (`POST /api/contracts/:id/sign`).
8. Investigator issues `Invoice` referencing Contract (`Invoice.LinkedContractID`).
9. Client pays invoice via `POST /api/invoices/:id/pay` → `payment_status = paid`, `paid_at = now()`.

### Flow: Quick Invoice (no contract)
- Create `Invoice` directly, send, receive payment.

---

## 10) Notes & Implementation Hints

- **Validation rules**:
  - Mandate must have client signature to be valid (business rule).
  - Invoice total must equal sum(line_items) + tax.
  - Invoice `due_date` must be >= `issue_date`.
- **Authorization**:
  - Only the investigator (or assigned operator) can create/send invoices.
  - Clients can accept Estimates and sign Mandates/Contracts.
- **Signature storage**:
  - Store e-sign envelope ids or signed PDF blobs externally (S3) and keep a URL/reference in the signature value.
- **PDF generation**:
  - Implement a job queue to render PDF on `send` (store PDF URL in document metadata).
- **Notifications**:
  - When sending a document, push email + optionally SMS + webhook to client systems.
- **Idempotency**:
  - `send` and `pay` operations should be idempotent (handle retries).
- **Data retention / GDPR**:
  - When archiving or deleting, ensure compliance (soft-delete with retention policy).
- **Testing**:
  - Unit-test state transitions (draft → sent → signed → archived).
  - Integration tests for linking flows (estimate → mandate → contract → invoice).
- **Audit / Compliance**:
  - Ensure `document_versions` capture the full content and who made changes.

---

## Appendix: Minimal JSON Schema Examples

### Estimate (example)
{
  "type": "estimate",
  "data": {
    "estimate_number": "EST-2025-001",
    "issue_date": "2025-10-12T09:00:00Z",
    "valid_until": "2025-11-12T09:00:00Z",
    "line_items": [
      { "description": "Surveillance (8h)", "quantity": 8, "unit_price": 50, "subtotal": 400 }
    ],
    "estimated_total": 400,
    "notes": "Travel expenses not included",
    "status": "draft"
  }
}

### Mandate (example)
{
  "type": "mandate",
  "data": {
    "mandate_number": "MND-2025-001",
    "issue_date": "2025-10-13T09:00:00Z",
    "scope_of_work": "Foot surveillance and photo evidence in city X",
    "valid_from": "2025-10-14T00:00:00Z",
    "valid_until": "2025-11-14T00:00:00Z",
    "terms_conditions": "Standard terms...",
    "linked_estimate_id": "uuid-of-estimate",
    "status": "signed"
  }
}

---

## Final checklist to implement

- [ ] Typed structs for Estimate, Mandate, Contract, Invoice
- [ ] `DocumentBase` and `Signature` types
- [ ] Persistence layer (typed tables or unified JSON)
- [ ] `document_versions` for audit/versioning
- [ ] Business operations: create/update/send/sign/archive/link/pay
- [ ] REST API routes (or GraphQL equivalents)
- [ ] PDF rendering job + storage
- [ ] Notifications (email/SMS/webhooks)
- [ ] Access control and validation rules
- [ ] Tests for flows and state transitions
- [ ] GDPR / retention policy handling

---
