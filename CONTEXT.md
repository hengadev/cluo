# Cluo — Domain Context

Cluo is a case management application for a private investigator (PI). It covers the full lifecycle of an investigation mission: from intake through fieldwork, report writing, document signing, and client delivery.

---

## Actors

| Actor | Description | Primary interface |
|---|---|---|
| **Investigator** | The PI. Single user who manages all cases, writes reports, and handles all fieldwork. | `cluo_desktop` (desktop app) + `cluo_mobile` (field PWA) |
| **Client** | The person or organisation that commissioned the investigation. Receives and signs documents. No full account — access to `cluo_web` is granted via a credential sent by email or SMS. | `cluo_web` (client portal) |

---

## Core Entities

### Case
A single investigation mission. Created when a Client commissions the Investigator to investigate something. Closed when the report is delivered and the case is settled.

**Lifecycle:**

```
draft → in_progress → ready → released
```

| Status | Meaning |
|---|---|
| `draft` | Case created, work has not started yet |
| `in_progress` | Investigation is actively ongoing |
| `ready` | Investigation complete, report ready, not yet delivered to client |
| `released` | Case delivered to client and formally closed |

A Case has:
- A **Client** (who commissioned it)
- An optional **CaseSubject** (the person being investigated)
- A **CaseType** (the nature of the investigation — see below)
- A location (where the investigation takes place)

**CaseType** is not a fixed enum. It is stored in a prefilled lookup table because the range of investigation types is too varied to hardcode (e.g. surveillance, insurance fraud, arson investigation, background check, missing person). Its primary purpose is to define a checklist of expected **Pièces** for a case — certain case types require specific exhibits that others do not.

### Client
The person or organisation that hired the Investigator. Can be a private individual, insurance company, law firm, corporate entity, or government body.

A Client is distinct from a CaseSubject — the Client commissions the work; the CaseSubject is the target of the investigation.

A Client can have one or more **Contacts**: named individuals within the client organisation (e.g. a claims adjuster at an insurance firm, a partner at a law firm). A Case can be assigned to a specific Contact, who acts as the day-to-day point of contact for that mission.

### CaseSubject
The person being investigated. Holds personal details (name, contact, address, occupation) and one or more **roles** within a case (e.g. suspect, witness). A CaseSubject can appear in multiple Cases.

### Rapport
The final written investigation report produced by the Investigator. A formal narrative document describing observations, timelines, and findings from the fieldwork, supported by published media. The Rapport is written in a rich-text editor (with AI writing assistance) in the desktop app and delivered to the Client as a downloadable PDF via the web portal (`cluo_web`). The Case moves to `ready` when the Rapport is complete, and to `released` when it has been delivered to the Client.

### Réseaux (OSINT Research)
A section of a Case dedicated to open-source intelligence (OSINT) findings about the CaseSubject — social media profiles, public online presence, usernames, notable posts. The Investigator is not trained in OSINT; this section is planned to be AI-assisted, surfacing relevant findings automatically from the subject's known details.

### Pièces (Exhibits)
Supporting documents attached to a Case that were received externally or collected outside of direct fieldwork. Includes scanned documents, client-provided files, court exhibits, receipts, and any file that doesn't fit the photo/video/audio categories and isn't one of the four formal legal documents. Acts as a general-purpose file attachment bucket for the case.

### Media
A `MediaFile` is any image, video, or audio file collected during fieldwork and attached to a Case. The investigator captures a large volume of raw media; he selects a subset to publish to the client.

- **Unpublished** (`IsPublished = false`): internal only, visible to the investigator in the desktop app.
- **Published** (`IsPublished = true`): included in the deliverable; visible to the client in the web portal.

### Documents
Four document types cover the legal and financial lifecycle of a Case. All documents are expected to be present — the profession is regulated and having everything in order is the norm.

| Document | Phase | Purpose |
|---|---|---|
| **Estimate** | Pre-investigation | Quote detailing the expected scope and cost of the mission |
| **Mandate** | Pre-investigation | Legal authorisation for the Investigator to conduct the investigation (required by French PI regulation) |
| **Contract** | Pre-investigation | Commercial service agreement covering payment terms, confidentiality, and liability |
| **Invoice** | Post-investigation | Bill issued to the Client once the mission is complete |

Documents are not a rigid sequential workflow. Estimate, Mandate, and Contract are typically resolved before fieldwork starts; the Invoice is issued at the end. Each document carries a status (`draft → sent → signed → active → archived`) and can be linked to the others for traceability. The Mandate and Contract are signed out-of-band (wet signature or third-party e-sign tool); the client portal is read-only.

### AI Features
Four AI-powered capabilities assist the Investigator throughout the case lifecycle.

| Feature | When used | What it does |
|---|---|---|
| **Speech-to-text** | After fieldwork | Transcribes audio recordings captured in the field into text |
| **Transcript analysis** | After transcription | One-shot structured extraction over a completed transcription — identifies people, key observations, timeline facts, and produces a structured summary for use in the Rapport |
| **Text transformation** | During Rapport writing | Rewrites or improves selected passages in the Rapport (grammar, tone, clarity) |
| **AI chat** | Anytime | Interactive assistant with full case context — the PI can ask open-ended questions about the case, get suggestions, or think through findings |

---

## Subfolders

### `cluo_api`
Go backend. Sole owner of all business logic and data persistence. Follows **hexagonal architecture** (ports & adapters):

| Internal package | Responsibility |
|---|---|
| `internal/domain/` | Pure business entities with no external dependencies. One subdirectory per domain: `case`, `client`, `case_subject`, `document`, `media`, `ai`, `user`. All sensitive fields encrypted at rest via `encx`. |
| `internal/ports/` | Interface definitions (repository contracts, service contracts). Decouples the application layer from infrastructure. |
| `internal/application/` | Use cases and business workflows. Orchestrates domain entities through ports. |
| `internal/adapters/http/` | HTTP handlers and route registration. Maps HTTP requests to application use cases; maps errors to status codes. |
| `internal/adapters/postgres/` | PostgreSQL repository implementations. One subdirectory per domain. |
| `internal/adapters/redis/` | Redis cache implementations. |
| `internal/adapters/external/` | Third-party service adapters (e.g. AI providers, storage, email/SMS). |
| `cmd/` | Service entry points and wiring (dependency injection, server startup). |
| `migrations/` | SQL migration files. Naming convention: `{timestamp}_{domain}_{action}.sql`. |

### `cluo_desktop`
Tauri + SvelteKit desktop app for the **Investigator**. Used at the office (or anywhere with the desktop). Provides the full case management interface:

| Route | Responsibility |
|---|---|
| `/cases` | List of all cases with search and filtering |
| `/cases/[id]/informations` | Case details, client info, subject info, metadata |
| `/cases/[id]/photos` | Photo management — organise, caption, publish/unpublish |
| `/cases/[id]/recordings` | Audio recordings — playback, trigger transcription |
| `/cases/[id]/rapport` | Rich-text editor for the final report, with AI writing assistance and chat sidebar |
| `/cases/[id]/documents/[type]` | Legal and financial documents (Estimate, Mandate, Contract, Invoice) |
| `/cases/[id]/pieces` | Exhibit file attachments — upload and manage supporting documents |
| `/cases/[id]/reseaux` | OSINT research section — AI-assisted social media and public record findings on the CaseSubject |

### `cluo_mobile`
SvelteKit **PWA** for the **Investigator in the field**. Optimised for mobile use during surveillance. Primary function: capture audio recordings and upload them to the backend, linked to the active Case. Recordings are then transcribed by AI on the backend. Transcriptions feed into the Rapport writing workflow on the desktop.

### `cluo_web`
SvelteKit web app for the **Client**. Credential-gated — no persistent account. Access is granted by a one-time code sent by email or SMS for a specific Case. The client can:
- View all Case documents (Estimate, Mandate, Contract, Invoice)
- Download the Rapport PDF once the Case is `released`

Signing of Mandate and Contract is handled out-of-band; by the time the case is `in_progress` both documents are already signed.

### `infrastructure`
Terraform (provisioning) and Ansible (configuration) for the VPS deployment. Manages server setup, environment configuration, and service deployment for `cluo_api`, `cluo_web`, and `cluo_mobile`.
