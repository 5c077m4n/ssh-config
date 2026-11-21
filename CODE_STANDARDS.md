# Code Standards

- Use full descriptive variable names; no single-letter names (except
  conventional receivers/indices like `i`, `j`).
- Run `gofmt`/`go vet ./...` before committing; formatting is tool-enforced.
- Use testify (`require`/`assert`) in tests.
- Prefer explicit `any` over `interface{}`.
- Wrap errors with context (`fmt.Errorf` with `%w`); sentinel errors live
  outside functions.
- No comments unless they explain something the code cannot.
- Declare functions before using them (callers appear above callees where
  practical).
