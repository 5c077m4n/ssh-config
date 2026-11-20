# Code Standards

- Run `gofmt`/`go vet ./...` before committing; formatting is tool-enforced.
- Use testify (`require`/`assert`) in tests.
- Prefer explicit `any` over `interface{}`.
- Wrap errors with context (`fmt.Errorf` with `%w`); sentinel errors live outside functions.
- No comments unless they explain something the code cannot.
