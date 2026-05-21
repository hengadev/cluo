# Cluo — Project Instructions

## Skill output conventions

All planning artefacts are stored locally and gitignored under `.local/` at the repo root.

| Skill | Output directory | File naming |
|---|---|---|
| `/to-prd` | `.local/prd/` | `<short-kebab-description>.md` |
| `/to-issues` | `.local/issues/` | `<NNN>-<short-kebab-description>.md` (one file per issue, zero-padded) |

When running `/to-prd`: write the PRD to `.local/prd/`. When running `/to-issues`: read the PRD from `.local/prd/` and write individual issue files to `.local/issues/`. Do not publish to any external issue tracker — everything stays local.
