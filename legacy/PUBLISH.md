# Publishing keyp to npm

## Quick Start

### 1. Update package.json

Replace placeholders:
- `"author": "Your Name <your.email@example.com>"`
- `"url": "https://github.com/TheEditor/keyp.git"`

### 2. Create npm Account (if needed)

```bash
npm adduser
```

Or login:
```bash
npm login
```

### 3. Test Locally

```bash
# Make the CLI executable
chmod +x bin/keyp.js

# Test it works
node bin/keyp.js

# Should output the "coming soon" message
```

### 4. Publish to npm

```bash
# From the keyp/ directory
npm publish
```

**That's it!** The name is now reserved.

### 5. Verify It Worked

```bash
# Check on npm
npm view keyp

# Or visit
# https://www.npmjs.com/package/keyp
```

### 6. Test Installation

```bash
# Install globally
npm install -g keyp

# Run it
keyp

# Should see the "coming soon" message
```

## Important Notes

- **Must be inside the `keyp/` folder** when running `npm publish`
- **Name will be claimed** as soon as you publish
- **Can update anytime** with `npm version patch` then `npm publish`

## Next Steps After Publishing

1. **Create GitHub repo** at github.com/TheEditor/keyp
2. **Push code to GitHub**
   ```bash
   git init
   git add .
   git commit -m "Initial commit - placeholder v0.0.1"
   git branch -M main
   git remote add origin https://github.com/TheEditor/keyp.git
   git push -u origin main
   ```
3. **Update every few weeks** to show active development
4. **Build the actual CLI** following the specification

## Version Updates

As you develop, publish updates:

```bash
# Increment version
npm version patch  # 0.0.1 -> 0.0.2

# Publish
npm publish
```

Good luck! 🚀
