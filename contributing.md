# Contributing

Thank you for your interest in improving sql-migrator.

## Reporting issues

When you open an issue, include:

- The **exact command** you ran (including `-from`, `-to`, and paths).
- A **minimal** SQL snippet that shows the problem (or states clearly if the input cannot be shared).
- What you **expected** versus what you **got**.

## Pull requests

1. **Fork** the repository and create a **branch** for your change (`fix/postgres-timestamp`, `feat/oracle-dialect`, etc.).
2. **Keep changes focused**: one logical change per pull request is easier to review than a large refactor mixed with unrelated edits.
3. **Run tests** before submitting:

   ```bash
   go test ./...
   go build ./...
   ```

4. **Match existing style**: naming, package layout, and the level of commentary in `internal/convert` and `internal/dialect`.
5. If you add conversion rules, consider adding a **small test** in `internal/convert/convert_test.go` that locks in the behavior you care about.

## Design notes

- Conversion is **regex-based** and intentionally simple. New rules should be ordered carefully (more specific patterns before general ones) to avoid partial matches breaking valid SQL.
- Avoid replacements that are likely to corrupt **string literals** or **comments** unless the pattern is very narrow.
- Documenting limitations in the README is preferred over over-engineering a single rule.

## Code of conduct

Be respectful and constructive in issues and reviews. Assume good intent.
