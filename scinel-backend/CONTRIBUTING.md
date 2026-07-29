@"
# Contributing to ScIntel Backend

First off, thank you for considering contributing to ScIntel! 🎉

## Code of Conduct
By participating in this project, you agree to abide by our code of conduct: Be respectful, professional, and constructive.

## Getting Started
1. Fork the repository
2. Create a new branch: \`git checkout -b feature/your-feature-name\`
3. Make your changes
4. Run tests: \`make test\`
5. Run linter: \`make lint\`
6. Commit your changes (see commit conventions below)
7. Push to your fork and open a Pull Request

## Development Standards

### Go Code Style
- Run \`gofmt -s -w .\` before committing (included in \`make format\`)
- Use \`golangci-lint\` to catch common issues
- Follow Go naming conventions (CamelCase for exported, camelCase for unexported)

### Commit Conventions
We follow [Conventional Commits](https://www.conventionalcommits.org/):
- \`feat:\` New feature
- \`fix:\` Bug fix
- \`docs:\` Documentation updates
- \`test:\` Testing changes
- \`chore:\` Maintenance tasks

**Example:** \`feat(auth): implement JWT refresh token endpoint\`

### Testing Requirements
- **Coverage**: All new code must maintain >80% coverage
- **Unit Tests**: Must pass locally before pushing
- **Integration Tests**: Must pass in CI pipeline

## Pull Request Process
1. Ensure your PR description clearly describes the changes and why they are needed.
2. Update the README.md or documentation if necessary.
3. Wait for CI to pass. All checks must be green.
4. Request a review (or self-review for solo projects).

## Reporting Bugs
Please use GitHub Issues with the \`bug\` label. Include:
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, etc.)

Thanks for contributing! 🚀
"@ | Out-File -FilePath CONTRIBUTING.md -Encoding utf8