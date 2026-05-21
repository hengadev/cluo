# ADR 0002 — One primary CaseSubject per Case

## Status
Accepted

## Context
The initial schema included a `cases.case_subject_cases` junction table supporting a many-to-many relationship between Cases and CaseSubjects, with a `roles[]` array per row. The `Investigation` Go struct also has a `case_subject_id` FK column pointing to a single primary subject.

Two competing models existed simultaneously.

## Decision
A Case has **at most one** CaseSubject. The `cases.cases.case_subject_id` FK is the authoritative link. The `case_subject_cases` junction table will be dropped.

The `Subject` domain entity's `CaseID []uuid.UUID` and `Roles []PersonRole` fields (which loaded from the junction table) will also be removed.

## Rationale
The PI's workflow always targets a single identifiable person per investigation. A many-to-many adds query complexity and an extra join for no real benefit. If multi-subject cases are ever needed, a migration adding the junction table back can be done with full context.

## Consequences
- `cases.case_subject_cases` table and its migration must be dropped.
- `Subject.CaseID` and `Subject.Roles` fields removed from the domain entity.
- CaseSubject CRUD is independent; attachment to a Case is done by setting `case_subject_id` on the Case.
