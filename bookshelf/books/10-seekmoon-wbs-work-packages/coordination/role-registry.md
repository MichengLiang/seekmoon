# SeekMoon WBS Role Registry

## Active Roles

| Agent | Runtime role | Work package | Boundary | Reuse intent | Status |
|---|---|---|---|---|---|
| principal coordinator | coordinator | full WBS | process artifacts, handoffs, review routing, promotion decisions | continuous | active |
| `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`) | builder | Batch D revision 1 | service behavior files listed in Batch D handoff | persistent specialist reuse | revision returned |
| `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`) | reviewer | Batch D re-review | independent review, evidence check, commit on approval | reused for Batch D re-review | re-review ready |

## Rules For Specialized Executors

- Do not directly address the user.
- Do not promote your own work.
- Do not revert changes made by others.
- Read the handoff and required sources before editing or reviewing.
- Put long reports in the requested file path, then summarize briefly in the final response.
- Use `git commit --only -m "..." -- <paths>` for your own committed files when committing is requested and necessary. New files must be tracked first.

## Current Agent IDs

- Batch A builder: `019ef59d-058d-78c1-bcb6-9cf1417d8b8c` (`Dewey`)
- Batch A reviewer: `019ef5b4-df23-7151-85dd-41239d63c743` (`Curie`)
