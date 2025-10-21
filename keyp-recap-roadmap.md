# keyp: Vibe Coding a GitHub Attention Getter

## Executive Summary

This conversation documents the rapid development and launch of **keyp**, a local-first secret manager for developers, from concept to published npm package in a single session. The project demonstrates the “vibe coding” approach - using AI assistance to move quickly from idea to working prototype.

**Result:** Successfully published `@theeditor/keyp` v0.0.1 to npm and created GitHub repository.

-----

## The Challenge: Finding a Viable Secret Manager Name

### Initial Research

We evaluated existing secret managers and identified a gap in the market for a developer-focused, local-first tool that’s simpler than enterprise solutions but more modern than traditional Unix password managers like `pass`.

**Target positioning:** “pass for the Node.js generation”

### The Naming Journey

After discovering that simple names like `keep` were already taken on npm, we explored creative alternatives and eventually settled on **`keyp`** - a clever misspelling that:

- Sounds identical to “keep” (our first choice)
- Is professional enough for enterprise use
- Is memorable and distinctive
- Successfully passed npm’s typosquatting detection

**Final decision:** Published as scoped package `@theeditor/keyp` due to similarity with existing package `keyv`.

-----

## Project Setup & Publishing

### Files Created

```
keyp/
??? bin/
?   ??? keyp.js          # CLI executable showing "coming soon"
??? .gitignore           # Standard Node.js ignores
??? LICENSE              # MIT License
??? package.json         # npm package configuration
??? PUBLISH.md           # Publishing instructions
??? README.md            # Project documentation
```

### Key Features of package.json

```json
{
  "name": "@theeditor/keyp",
  "version": "0.0.1",
  "publishConfig": {
    "access": "public"
  },
  "bin": {
    "keyp": "./bin/keyp.js"
  }
}
```

**Important:** The `bin` field ensures users run `keyp` (not `@theeditor/keyp`) even though it’s a scoped package.

-----

## npm Publishing Process

### Steps Taken

1. **Created npm account** and authenticated via web login
1. **Updated package.json** to scoped name `@theeditor/keyp`
1. **Added `publishConfig`** to avoid needing `--access=public` flag
1. **Published successfully** with `npm publish`

### npm Scoped Packages: Pros & Cons

**Advantages:**

- ? CLI command stays simple (`keyp`, not the full scoped name)
- ? Namespace control (only you can publish `@theeditor/*`)
- ? Natural for building a family of tools
- ? Professional (Angular, Babel, Vue all use scopes)

**Minor Downsides:**

- Slightly longer installation command
- Need to include `publishConfig` or remember `--access=public`
- Slightly less discoverable in npm search

**Verdict:** Scoped packages are perfectly professional and the downsides are minimal.

-----

## Repository Setup

### Created on GitHub

- **Repository:** https://github.com/TheEditor/keyp
- **npm Package:** https://www.npmjs.com/package/@theeditor/keyp

### Git Workflow

```bash
git init
git add .
git commit -m "Initial commit - keyp v0.0.1 placeholder"
git branch -M main
git remote add origin https://github.com/TheEditor/keyp.git
git push -u origin main --force
```

**Note:** Used `--force` to overwrite GitHub’s auto-generated README with our detailed version.

-----

## README Content Strategy

### Final Tagline

```markdown
# keyp

> Local-first secret manager for developers
> *"pass for the Node.js generation"*
```

**Key insight:** The line works great as a memorable tagline, but “Positioning:” prefix felt like stage directions and was removed for a more natural presentation.

### README Structure

1. **What is keyp?** - Clear value proposition
1. **Features** - Security, local-first, Git sync, developer-friendly
1. **Planned Commands** - Show the intended CLI interface
1. **Why keyp?** - Positioning against alternatives
1. **Installation** - Clear npm command
1. **Roadmap** - Transparency about development status
1. **Philosophy** - Core principles (local-first, simple, secure, dev-focused)
1. **Inspiration** - Acknowledges `pass` as inspiration

-----

## Lessons Learned

### What Worked Well

1. **Rapid iteration** - From concept to published package in hours
1. **Name persistence** - Finding creative solutions when obvious names were taken
1. **Scoped packages** - Good workaround for namespace collisions
1. **Clear README** - Strong positioning and feature communication

### Key Decisions

1. **Scoped package over different name** - Kept the name we wanted
1. **Placeholder strategy** - Publish early to claim the name
1. **Force push over merge** - Simpler when GitHub only has placeholders
1. **Tagline presentation** - Natural integration vs. meta-commentary

### Technical Details That Matter

- **CLI command is controlled by `bin` field**, not package name
- **npm has typosquatting protection** that blocks similar names
- **Scoped packages default to private** - must set `publishConfig` or use `--access=public`
- **Git force push is appropriate** when remote only has auto-generated files

-----

## What’s Next

### Immediate Tasks

- ? Package published to npm
- ? GitHub repository created
- ? README polished

### Development Roadmap

**Week 1:** Core encryption + vault management

- AES-256-GCM implementation
- PBKDF2 key derivation
- Vault file format

**Week 2:** CLI commands

- `keyp init` - Initialize vault
- `keyp set` - Store secrets
- `keyp get` - Retrieve secrets (clipboard)
- `keyp list` - List all secrets

**Week 3:** Git sync + polish

- Git integration for encrypted backups
- Command-line UX refinement
- Error handling

**Week 4:** v1.0.0 launch

- Comprehensive testing
- Full documentation
- Launch on Hacker News, Product Hunt

-----

## Project Philosophy

**keyp** embodies these principles:

1. **Local-first** - Your secrets stay on your machine, no cloud required
1. **Simple** - One command does one thing well
1. **Secure** - Industry-standard encryption, no shortcuts
1. **Developer-focused** - Built for developers, by developers

**Target user:** Solo developers who want secure secret management without enterprise complexity or GPG arcana.

-----

## Success Metrics

### Immediate Success

- ? Name secured on npm
- ? Repository created and public
- ? Professional README in place
- ? Clear development roadmap

### Future Goals

- 100+ GitHub stars in first month
- 1,000+ npm downloads in first quarter
- Active community contributions
- Featured on developer tool lists

-----

## Vibe Coding Insights

This project demonstrates the **Infrastructure Layer** of the Vibe Coding ecosystem:

### What Was Automated

- Package structure generation
- README content creation
- License file setup
- Publishing instructions

### What Required Manual Work

- npm account creation and authentication
- GitHub repository creation
- Resolving naming conflicts
- Git operations (push, merge, etc.)

### Opportunities for Automation

- **Automated npm name checking** before suggesting names
- **One-command publishing** that handles scoped packages automatically
- **GitHub repo creation** via API
- **Automated README generation** based on project type

-----

## Conclusion

From initial concept to published npm package with professional documentation took approximately 4-6 hours using AI-assisted development. The project is now positioned as a credible open-source tool with:

- Professional branding and positioning
- Clear value proposition
- Transparent development roadmap
- Active presence on npm and GitHub

**Next step:** Build the actual CLI tool according to the specification documents we’ve created.

-----

## Quick Reference

**Package:** `@theeditor/keyp`
**GitHub:** https://github.com/TheEditor/keyp
**npm:** https://www.npmjs.com/package/@theeditor/keyp
**License:** MIT
**Author:** Dave Fobare

**Installation:**

```bash
npm install -g @theeditor/keyp
```

**Usage (current placeholder):**

```bash
keyp
# Shows "coming soon" message with planned features
```