# Claude Development Guide for Miller

## Project Overview

Miller is a command-line data processing tool for working with CSV, TSV, JSON, and other data formats. It's written in Go (v1.18+) and provides SQL-like operations on data.

## Initial Setup

### Setting Up staticcheck

The `make staticcheck` target requires the staticcheck tool. To set it up:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

This installs staticcheck to `~/go/bin/` (the default Go binaries directory). For `make staticcheck` to work, you need `~/go/bin` in your `PATH`.

**Add to your shell profile** (`.bashrc`, `.zshrc`, or equivalent):
```bash
export PATH="$PATH:$HOME/go/bin"
```

**Verify the installation:**
```bash
staticcheck -version
```

If this works, you can use `make staticcheck` without any further setup.

## Build & Test

### Building
```bash
make build          # Build the mlr executable
make quiet          # Build silently (no output messages)
```

### Testing
```bash
make check          # Run all tests (unit + regression)
make unit-test      # Run unit tests only
make regression-test # Run regression tests only
make bench          # Run benchmarks
```

### Code Quality
```bash
make fmt            # Format code with go fmt
make staticcheck    # Run static analysis (see Initial Setup section above)
```

### Full Developer Workflow
```bash
make dev            # Format, build, test, generate docs (comprehensive check before pushing)
```

## Project Structure

- `cmd/mlr` - Main Miller executable entry point
- `pkg/` - Core library code organized by functionality
  - `pkg/lib/` - Utility libraries
  - `pkg/scan/` - Input scanning
  - `pkg/mlrval/` - Miller value types
  - `pkg/bifs/` - Built-in functions
  - `pkg/input/` - Input format handlers
- `regression_test.go` - Regression test suite
- `docs/` - Documentation (Markdown with live code samples)
- `man/` - Man page generation

## Development Conventions

### Code Style
- Follow Go conventions and `go fmt` output
- Use meaningful variable and function names
- Keep functions focused and testable

### Testing
- Add unit tests for new functionality in `pkg/*/` directories
- Use `go test` for unit tests
- Run `make regression-test` for integration testing
- The `mlr regtest` command provides more control for interactive debugging

### Documentation
- Update relevant `.md.in` files in `docs/src/` when adding features
- These are processed into live documentation with actual code examples
- Run `make dev` to rebuild documentation

### Git Workflow
- Create descriptive commit messages
- Reference issue numbers when fixing bugs or implementing features
- Test locally with `make check` before committing
- For major changes, run `make dev` to ensure docs and tests pass

## Key Dependencies

Miller uses the Go standard library. Check `go.mod` for specific versions.

## Common Tasks

### Adding a New Built-in Function
1. Implement in appropriate package (likely `pkg/bifs/`)
2. Add unit tests alongside
3. Update documentation in `docs/src/`
4. Run `make check` to verify

### Fixing a Bug
1. Create a minimal test case (unit test or regression test)
2. Fix the code
3. Verify with `make check`
4. Commit with reference to issue number

### Performance Optimization
1. Run `make bench` to establish baseline
2. Make changes
3. Run benchmarks again to measure improvement
4. Consider adding permanent benchmark in test files

## Documentation Building

Documentation is built from `.md.in` template files that contain live code samples executed via Miller itself. When you make changes that affect command output or behavior, you may need to update these templates and rebuild docs with `make -C docs/src forcebuild`.

## Before Pushing

Always run:
```bash
make dev
```

This ensures code formatting, builds successfully, passes all tests, and documentation is up to date.

## Additional Resources

- [Full documentation](https://miller.readthedocs.io/)
- [Contributing guidelines](https://miller.readthedocs.io/en/latest/contributing/)
- [Issue labeling notes](https://github.com/johnkerl/miller/wiki/Issue-labeling)
