# Contributing to keyp

Thank you for your interest in contributing to keyp! We appreciate all kinds of contributions, whether it's bug reports, feature requests, documentation improvements, or code contributions.

## Code of Conduct

Please be respectful and constructive in all interactions. We're committed to providing a welcoming and inclusive environment for all contributors.

## How to Report Bugs

### Before Submitting a Bug Report

- Check the [issue list](https://github.com/TheEditor/keyp/issues) to see if the bug has already been reported
- Check the [troubleshooting guide](./TROUBLESHOOTING.md) for common issues
- Verify you're using the latest version: `npm install -g @theeditor/keyp@latest`

### Submitting a Bug Report

Create an issue with the following information:

```markdown
**Description:** Brief description of the bug

**Steps to Reproduce:**
1. First step
2. Second step
3. ...

**Expected Behavior:** What you expected to happen

**Actual Behavior:** What actually happened

**Environment:**
- keyp version: [e.g., 0.2.0]
- Node version: [e.g., 18.0.0]
- OS: [e.g., macOS 13.0, Ubuntu 22.04, Windows 11]

**Additional Context:** Any other relevant information
```

## How to Suggest Enhancements

Create an issue with the following information:

```markdown
**Title:** Brief description of enhancement

**Description:** Detailed description of what you'd like to see

**Use Case:** Why would this be useful?

**Possible Implementation:** Optional: How you might implement this

**Related Issues:** Links to related issues
```

## Development Setup

### Prerequisites

- Node.js 18.0.0 or higher
- npm 8.0.0 or higher
- Git

### Local Development

1. Fork the repository
   ```bash
   git clone https://github.com/YOUR-USERNAME/keyp.git
   cd keyp
   ```

2. Install dependencies
   ```bash
   npm install
   ```

3. Build the project
   ```bash
   npm run build
   ```

4. Run tests
   ```bash
   npm test
   ```

5. Watch mode for development
   ```bash
   npm run dev
   ```

### Project Structure

```
keyp/
├── src/                    # TypeScript source code
│   ├── crypto.ts          # Encryption/decryption
│   ├── vault-manager.ts   # Vault lifecycle
│   ├── secrets.ts         # Secret CRUD
│   ├── types.ts           # Type definitions
│   ├── config.ts          # Configuration
│   ├── git-sync.ts        # Git synchronization
│   ├── cli/               # CLI commands
│   │   ├── index.ts       # Main CLI entry
│   │   ├── utils.ts       # Shared utilities
│   │   └── commands/      # Individual commands
│   ├── *.test.ts          # Unit tests
│
├── lib/                    # Compiled JavaScript
├── bin/                    # Executable entry point
├── completions/           # Shell completion scripts
├── docs/                  # Documentation
└── package.json           # Package configuration
```

## Code Style

We use TypeScript with strict type checking. Follow these guidelines:

### TypeScript

- Use `strict: true` in tsconfig.json
- No `any` types without explicit `// @ts-ignore` comment with reason
- Use descriptive variable names
- Add JSDoc comments for public functions

### Example

```typescript
/**
 * Encrypts data using AES-256-GCM
 *
 * @param plaintext - Text to encrypt
 * @param password - Master password
 * @returns Encryption result with ciphertext, IV, salt, and auth tag
 */
export function encrypt(plaintext: string, password: string): EncryptionResult {
  // Implementation
}
```

### Formatting

- Use 2-space indentation
- Use single quotes for strings (except JSON)
- Line length: max 100 characters
- Use `const` by default, `let` when needed, avoid `var`

## Testing

### Running Tests

```bash
# Run all tests
npm test

# Run specific test file
node --test lib/crypto.test.js

# Watch mode (during development)
npm run dev
```

### Writing Tests

- Use Node's built-in test runner (no external framework)
- Follow existing test patterns
- Aim for comprehensive coverage
- Test both success and error cases
- Use descriptive test names

```typescript
describe('Function name', () => {
  it('should do something when condition is met', () => {
    // Arrange
    const input = 'test';

    // Act
    const result = functionUnderTest(input);

    // Assert
    assert.strictEqual(result, 'expected');
  });
});
```

## Committing Changes

### Before Committing

1. Run tests: `npm test`
2. Build: `npm run build`
3. Check TypeScript: `npm run tsc`
4. Fix any formatting issues

### Commit Messages

Use clear, descriptive commit messages following conventional commits:

```
type(scope): subject

body

footer
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring without feature/bug changes
- `test`: Adding or updating tests
- `chore`: Build process, dependencies, etc.

**Examples:**

```
feat(crypto): add support for AES-256-GCM encryption

Implements authenticated encryption using AES-256-GCM with PBKDF2
key derivation for secure vault encryption.

Fixes #123
```

```
fix(cli): handle special characters in secret names

Properly escape special characters when storing secret names to
prevent command injection vulnerabilities.

Fixes #456
```

## Submitting Pull Requests

### Before Submitting

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Add/update tests
5. Update documentation if needed
6. Commit with clear messages
7. Push to your fork
8. Submit a pull request

### Pull Request Template

```markdown
## Description
Brief description of what this PR does

## Related Issues
Fixes #123

## Type of Change
- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] Tested on Windows
- [ ] Tested on macOS
- [ ] Tested on Linux

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests cover new code
- [ ] CHANGELOG updated
```

### Review Process

- At least one review required
- CI tests must pass
- Code must follow project standards
- Documentation must be complete

## Documentation

### When to Update Docs

- New features need documentation
- API changes need documentation updates
- Bug fixes that change behavior should be documented
- Performance improvements should be highlighted

### Documentation Files

- **CLI.md** - Command-line interface documentation
- **API.md** - Library API reference
- **SECURITY.md** - Security information
- **GIT_SYNC.md** - Git synchronization guide
- **TROUBLESHOOTING.md** - Common issues and solutions
- **README.md** - Project overview and quick start

## Release Process

*Maintainers Only*

1. Update version in package.json
2. Update CHANGELOG.md
3. Create git tag: `git tag v1.2.3`
4. Push: `git push --tags`
5. CI/CD pipeline publishes to NPM

## Getting Help

- Check [documentation](../README.md#documentation)
- Browse [existing issues](https://github.com/TheEditor/keyp/issues)
- Read [troubleshooting guide](./TROUBLESHOOTING.md)
- Start a [discussion](https://github.com/TheEditor/keyp/discussions)

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub repository contributors list

Thank you for contributing to keyp! 🔒
