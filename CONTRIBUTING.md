# Contributing to Crux

Thank you for contributing to Crux! This document outlines the processes and standards for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Install pre-commit hooks: `make hooks`
4. Create a feature branch: `git checkout -b feature/your-feature-name`
5. Make your changes
6. Run the development workflow: `make dev`
7. Commit your changes following our commit message conventions
8. Push to your fork and create a pull request

## Development Workflow

### Prerequisites

- Go 1.26 or later
- `golangci-lint` installed and available in PATH
- Git

### Common Commands

```bash
make hooks      # Install pre-commit hooks
make fmt        # Format code
make lint       # Run linter
make vet        # Run go vet
make test       # Run tests with race detector
make build      # Build the binary
make dev        # Run full development workflow (fmt + lint + test)
```

## Dependency Management

### Adding New Dependencies

All new dependencies must be reviewed and approved before being added to the project. Follow this process:

1. **Propose the dependency** in your PR description:
   - Package name and version
   - Purpose and justification
   - Why existing dependencies cannot solve the problem
   - Licence information

2. **Verify the licence**:
   - Check the package licence is compatible with our project
   - Acceptable licences: MIT, Apache 2.0, BSD-3-Clause, ISC
   - Unacceptable licences: GPL, AGPL, LGPL (without explicit approval)

3. **Add the dependency**:
   ```bash
   go get package@version
   go mod tidy
   ```

4. **Commit both files**:
   ```bash
   git add go.mod go.sum
   git commit -m "deps: add package-name for purpose"
   ```

5. **PR review**:
   - At least one reviewer must confirm the licence is acceptable
   - Reviewer must verify the dependency is necessary
   - `go mod tidy` must produce no changes

### Updating Dependencies

1. **Propose the update** with reasoning in the PR description
2. **Update the dependency**:
   ```bash
   go get package@version
   go mod tidy
   ```
3. **Test thoroughly** - dependency updates can introduce breaking changes
4. **Document** any API changes in the PR description

### Current Dependencies

| Package | Version | Licence | Purpose |
|---|---|---|---|
| `github.com/spf13/cobra` | v1.10.2 | Apache-2.0 | CLI framework |
| `github.com/stretchr/testify` | v1.11.1 | MIT | Test assertions |

### Dependency Rules

- **Always run `go mod tidy`** after any dependency change
- **Always commit `go.sum`** - it ensures reproducible builds
- **Never vendor** unless explicitly decided by the team (requires ADR)
- **Pin versions** - avoid `@latest` in production code
- **Minimize dependencies** - prefer standard library when possible

## Code Style

Follow the standards documented in `docs/development-standards.md`:

- Format code with `gofmt`
- Follow Go naming conventions
- Keep functions under 50 lines
- Maximum 4 function parameters
- Write table-driven tests
- Maintain 80% test coverage minimum

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, no code change
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `deps`: Dependency updates

### Examples

```
feat(cli): add new command for service generation

Implements the core CLI command structure using cobra.
Adds basic flag parsing and validation.

Closes #123
```

```
deps: add cobra v1.10.2 for CLI framework

Cobra provides a robust CLI framework with command structure,
flag parsing, and help generation. Licence: Apache-2.0.
```

## Pull Request Process

1. **Create a feature branch** from `main`
2. **Implement changes** with tests
3. **Run full test suite**: `make dev`
4. **Update documentation** if needed
5. **Create PR** with clear description
6. **Address review comments**
7. **Squash commits** before merge (if requested)

### PR Description Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] Dependency update

## Testing
How has this been tested?

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] No new warnings
```

## Code Review Guidelines

### For Authors

- Keep PRs focused and reasonably sized
- Write clear commit messages
- Add tests for new functionality
- Update documentation
- Respond to feedback constructively

### For Reviewers

- Be constructive and specific
- Explain reasoning behind suggestions
- Use prefixes: `nit:`, `question:`, `suggestion:`, `blocker:`
- Approve when satisfied, even if minor nits remain

### Review Checklist

- [ ] Correctness: Does it work as intended?
- [ ] Tests: Are there adequate tests?
- [ ] Design: Does it follow architecture principles?
- [ ] Readability: Is it easy to understand?
- [ ] Performance: Are there obvious bottlenecks?
- [ ] Security: Are there security concerns?
- [ ] Dependencies: Are new dependencies justified and licensed correctly?

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas
- Reach out to the platform team for guidance

Thank you for contributing to Crux!
