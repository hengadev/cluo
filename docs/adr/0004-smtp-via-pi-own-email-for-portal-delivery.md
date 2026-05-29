# ADR-0004: SMTP via PI's own custom-domain email for portal delivery

## Status
Accepted

## Context
When a Case is `released`, the PI sends the Client a magic link to access `cluo_web`. The app must deliver this link by email automatically rather than requiring the PI to copy-paste it manually.

Two categories of delivery were considered:

- **Third-party transactional email SaaS** (Resend, Sendgrid, Mailgun, etc.): simple to integrate, no infrastructure to maintain, but all log sender/recipient metadata and process email content through their servers.
- **SMTP via the PI's own email account**: the app sends through a custom-domain address the PI already owns and controls (e.g. `contact@agence-xxx.fr`) using standard SMTP credentials.

## Decision
Email is sent via SMTP using the PI's own custom-domain address. No third-party transactional email SaaS is introduced.

## Rationale
The PI operates under French professional secrecy obligations (secret professionnel). Even a magic-link email contains sensitive metadata: the sender is a PI, the recipient is a client, and the timestamp reveals that an investigation is being delivered. Routing this through a US-based SaaS creates GDPR exposure and conflicts with deontological obligations, regardless of the email body content.

The PI already has a professional custom-domain email address. Using it as the SMTP sender keeps the delivery within infrastructure the PI already trusts, adds no new data processor, and remains professional in appearance.

## Consequences
- A real SMTP adapter must replace the current stub in `internal/adapters/external/`.
- SMTP credentials (host, port, from address, username, password) are stored as env vars on the VPS — never in code.
- Deliverability depends on the PI's email provider and correct SPF/DKIM configuration on their domain. This is a one-time setup, not an ongoing operational concern.
- No per-email cost and no third-party dependency to maintain or monitor.
