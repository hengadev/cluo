# ADR-0001: CaseType as a Prefilled Lookup Table

## Status
Accepted

## Context
A Case needs to be categorised by the type of investigation it represents (e.g. surveillance, insurance fraud, arson investigation, background check). The range of investigation types handled by the PI is too varied and open-ended to hardcode as a fixed enum in the domain model.

## Decision
`CaseType` is stored as a reference in a prefilled database lookup table rather than as a hardcoded enum. The table is maintained by the Investigator and can be extended over time without a code change.

The primary purpose of `CaseType` is to drive a checklist of expected **Pièces** (exhibits) for a Case — certain investigation types require specific supporting documents that others do not.

## Consequences
- New case types can be added without a code deployment.
- The application can validate that expected Pièces are present for a given CaseType, surfacing a checklist to the Investigator.
- The domain layer holds `CaseType` as a `string` rather than a typed enum; validation of acceptable values happens at the application/persistence layer.
