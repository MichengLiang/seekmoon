# Prompt For Batch E Re-Reviewer

Read this re-review package first:

`/home/t103o/workbench/projects/seekmoon/bookshelf/books/10-seekmoon-wbs-work-packages/coordination/review-packages/020-batch-e-re-review.md`

Then execute it exactly.

Operational constraints:

- Continue as the independent Batch E reviewer.
- You are not the builder and not the principal coordinator.
- You are not alone in the repository; do not revert unrelated changes.
- The re-review object includes the original WP13 acceptance files plus the broader `cmd/` and `internal/` lint/security remediation needed for WP13 gates.
- Preserve the include-reading rule: whole-file includes require whole-file reads; `lines=` includes require exactly those ranges.
- If approved, write the re-review report and commit only approved Batch E, remediation, and coordination paths.
- If rejected, write the report and do not commit.
- Keep final chat response brief with verdict, report path, and commit hash if created.
