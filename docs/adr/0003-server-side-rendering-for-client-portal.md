# ADR-0003: Server-side rendering for client portal documents

## Status
Accepted

## Context
The client portal (`cluo_web`) needs to:
1. Display the Rapport inline — stored on the backend as an encrypted TipTap JSON blob.
2. Generate downloadable PDFs of the Rapport and the four case documents (Estimate, Mandate, Contract, Invoice).

Two approaches were considered for each requirement.

**For Rapport display:**
- Option A: Ship TipTap in read-only mode to `cluo_web`; render the JSON client-side.
- Option B: Convert TipTap JSON to HTML in the Go API; send the rendered HTML to the portal.

**For PDF generation:**
- Option A: Browser print-to-PDF (`window.print()` with a print stylesheet).
- Option B: Generate PDFs server-side in the Go API; stream them as binary downloads.

## Decision
Both rendering concerns are handled server-side in the Go API:
- The Rapport is served as rendered HTML via `GET /token/{token}/report/html`.
- PDFs (per document, Rapport, and the full case archive) are generated and streamed by the API.

The portal remains a thin display layer — it renders HTML returned by the API and triggers download links.

## Rationale
- TipTap is a heavyweight editor library not suited to a minimal read-only portal; adding it as a dependency increases bundle size for zero gain.
- Browser print-to-PDF produces inconsistent output across browsers and cannot be automated for the "download everything" archive flow, which must assemble PDFs server-side regardless.
- Centralising rendering in the API means the portal has no knowledge of document structure; it simply displays what the API returns.

## Consequences
- The Go API needs a TipTap JSON → HTML converter.
- The Go API needs a PDF generation library.
- `cluo_web` stays lightweight: no document parsing, no editor dependencies.
- A change to document templates requires only a backend deployment, not a frontend release.
