# Cluo

Case management system for a private investigator (PI). Covers the full investigation lifecycle — intake, fieldwork, report writing, document signing, and client delivery.

## Structure

| Package | Description |
|---|---|
| `cluo_api` | Go backend (hexagonal architecture). Owns all business logic, PostgreSQL persistence, S3 storage, and AI integrations. |
| `cluo_desktop` | Wails + SvelteKit desktop app for the investigator. Full case management interface including report editor, media management, and legal documents. |
| `cluo_mobile` | SvelteKit PWA for field use. Captures and uploads audio recordings linked to active cases. |
| `cluo_web` | SvelteKit client portal. Token-gated, read-only access for clients to view and download case documents, the final report, and published media. |
| `infrastructure` | Terraform + Ansible for VPS provisioning and deployment. |

## Key Concepts

- **Case** — the core entity, progressing through `in_progress → ready → released`.
- **Rapport** — the final report, written in a rich-text editor and delivered as a PDF to the client.
- **Portal access** — the PI releases a case by emailing the client a magic link (token-gated, 30-day expiry). No client account is required.
- **AI features** — speech-to-text transcription of field recordings, transcript analysis, text transformation in the report editor, and an AI chat assistant with full case context.

## Stack

- **Backend:** Go, PostgreSQL, Redis, S3
- **Frontend:** SvelteKit (TypeScript)
- **Desktop shell:** Wails (Go + WebView)
- **Infrastructure:** Terraform, Ansible, Docker
