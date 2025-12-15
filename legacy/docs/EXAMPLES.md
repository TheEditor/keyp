# keyp Usage Examples

Real-world examples and common workflows for using keyp.

## Table of Contents

1. [Basic Secret Management](#basic-secret-management)
2. [Development Workflows](#development-workflows)
3. [Multiple Environments](#multiple-environments)
4. [Team Collaboration](#team-collaboration)
5. [CI/CD Integration](#cicd-integration)
6. [Automation Scripts](#automation-scripts)
7. [Advanced Patterns](#advanced-patterns)

## Basic Secret Management

### Store and Retrieve Secrets

```bash
# Initialize vault (first time only)
$ keyp init
Enter master password: ●●●●●●●●
Confirm master password: ●●●●●●●●
✓ Vault initialized successfully!

# Store API key
$ keyp set github-token
Enter master password: ●●●●●●●●
Enter value for "github-token": ●●●●●●●●●●●●
Enter master password to save: ●●●●●●●●
✓ Secret "github-token" saved

# Retrieve secret (copies to clipboard)
$ keyp get github-token
Enter master password: ●●●●●●●●
✓ Copied to clipboard (clears in 45 seconds)

# List all secrets
$ keyp list
Enter master password: ●●●●●●●●

  • github-token
  • db-password
  • api-key

3 secrets stored
```

### Search Secrets

```bash
# Find all GitHub-related secrets
$ keyp list --search github
Enter master password: ●●●●●●●●

  • github-token
  • github-ssh-key

2 secrets stored

# Get secret count
$ keyp list --count
Enter master password: ●●●●●●●●
5
```

### Manage Secrets

```bash
# Rename a secret
$ keyp rename old-name new-name
Enter master password: ●●●●●●●●
✓ Secret renamed: old-name → new-name

# Copy a secret
$ keyp copy prod-api-key staging-api-key
Enter master password: ●●●●●●●●
✓ Secret copied: prod-api-key → staging-api-key

# Delete a secret
$ keyp delete unused-key
Enter master password: ●●●●●●●●
Delete secret "unused-key"? (y/N): y
✓ Secret "unused-key" deleted
```

## Development Workflows

### Local Development Setup

```bash
# Initialize vault for your project
cd ~/projects/my-app
keyp init

# Store all development credentials
keyp set db-host localhost
keyp set db-user dev_user
keyp set db-password dev_password
keyp set api-key sk_test_123456

# Create a shell alias for quick access
alias getdb='keyp get db-password && echo'
alias getapi='keyp get api-key && echo'

# Use in terminal
$ getdb
# Password copied to clipboard

# Retrieve to variable
$ DB_PASS=$(keyp get db-password --stdout)
```

### Environment-Specific Secrets

```bash
# Store secrets for different environments
keyp set dev-api-key sk_test_dev
keyp set staging-api-key sk_test_staging
keyp set prod-api-key sk_live_prod

# Use in scripts
#!/bin/bash
ENV=$1
API_KEY=$(keyp get ${ENV}-api-key --stdout)
curl -H "Authorization: Bearer $API_KEY" https://api.example.com
```

### Shell Aliases

```bash
# Add to ~/.bashrc or ~/.zshrc

# Get secret and show briefly
alias keyp-show='keyp get'

# Get and copy
alias keyp-copy='keyp get'

# Get to stdout
alias keyp-stdout='keyp get --stdout'

# Get and pipe to another command
alias keyp-pipe='keyp get --stdout |'

# Usage examples
$ keyp-stdout | jq .  # Parse JSON secret
$ keyp-pipe xclip -i  # Copy to clipboard explicitly
```

## Multiple Environments

### Separate Vaults Per Environment

```bash
# For development
KEYP_DIR=~/.keyp/dev keyp init
KEYP_DIR=~/.keyp/dev keyp set api-key sk_test_dev

# For staging
KEYP_DIR=~/.keyp/staging keyp init
KEYP_DIR=~/.keyp/staging keyp set api-key sk_test_staging

# For production
KEYP_DIR=~/.keyp/prod keyp init
KEYP_DIR=~/.keyp/prod keyp set api-key sk_live_prod

# Use in scripts
#!/bin/bash
case "$ENVIRONMENT" in
  dev)
    KEYP_DIR=~/.keyp/dev keyp get api-key --stdout
    ;;
  staging)
    KEYP_DIR=~/.keyp/staging keyp get api-key --stdout
    ;;
  prod)
    KEYP_DIR=~/.keyp/prod keyp get api-key --stdout
    ;;
esac
```

### Synchronized Across Machines

```bash
# Machine A: Set up Git sync
keyp sync init https://github.com/myusername/keyp-backup.git
keyp sync push

# Machine B: Clone and pull
keyp sync init https://github.com/myusername/keyp-backup.git
keyp sync pull

# Both machines now have same secrets
keyp list
# Shows all secrets from Machine A

# After adding new secret on Machine B
keyp set new-secret new-value
keyp sync push

# Back on Machine A
keyp sync pull
keyp list
# Shows the new secret
```

## Team Collaboration

### Multiple Team Members

```bash
# Each developer has their own vault (unique password)
# Shared backup repository syncs their encrypted vaults

# Developer 1
keyp init  # Creates with their password
keyp sync init git@github.com:team/secrets-backup.git
keyp sync push

# Developer 2
keyp init  # Creates with their own password
keyp sync init git@github.com:team/secrets-backup.git
keyp sync pull  # Gets Developer 1's secrets

# Note: Developer 2 must have already set up git sync
# to pull Developer 1's vault
```

### Shared Vault on Single Machine

```bash
# Multiple users sharing one machine
# Each can use separate KEYP_DIR

# User 1
KEYP_DIR=~/.keyp-alice keyp init
KEYP_DIR=~/.keyp-alice keyp set personal-key value1

# User 2
KEYP_DIR=~/.keyp-bob keyp init
KEYP_DIR=~/.keyp-bob keyp set personal-key value2

# Keep separate shell aliases
alias alice-keyp='KEYP_DIR=~/.keyp-alice keyp'
alias bob-keyp='KEYP_DIR=~/.keyp-bob keyp'

# Use
alice-keyp list
bob-keyp list  # Different secrets!
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy

on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install keyp
        run: npm install -g @theeditor/keyp

      - name: Create vault
        run: |
          echo "${{ secrets.VAULT_PASSWORD }}" | keyp init
          echo "${{ secrets.VAULT_DATA }}" > ~/.keyp/vault.json.encrypted

      - name: Deploy application
        run: |
          API_KEY=$(keyp get api-key --stdout)
          DB_URL=$(keyp get db-url --stdout)
          npm run deploy
        env:
          API_KEY: ${{ env.API_KEY }}
          DB_URL: ${{ env.DB_URL }}
```

### GitLab CI

```yaml
# .gitlab-ci.yml
deploy:
  image: node:18
  script:
    - npm install -g @theeditor/keyp
    - echo "$VAULT_PASSWORD" | keyp init
    - echo "$VAULT_DATA" > ~/.keyp/vault.json.encrypted
    - export API_KEY=$(keyp get api-key --stdout)
    - export DB_URL=$(keyp get db-url --stdout)
    - npm run deploy
  only:
    - main
```

### Docker

```dockerfile
FROM node:18

# Install keyp
RUN npm install -g @theeditor/keyp

# Create vault from build args
ARG VAULT_PASSWORD
ARG VAULT_DATA
RUN mkdir -p ~/.keyp && \
    echo "$VAULT_DATA" > ~/.keyp/vault.json.encrypted && \
    echo "Init with password" || true

# Deploy application
COPY . /app
WORKDIR /app
RUN npm install
CMD ["node", "deploy.js"]
```

## Automation Scripts

### Backup Script

```bash
#!/bin/bash
# backup-secrets.sh - Back up all secrets

BACKUP_DIR="$HOME/keyp-backups"
mkdir -p "$BACKUP_DIR"

DATE=$(date +%Y-%m-%d)
BACKUP_FILE="$BACKUP_DIR/secrets-$DATE.keyp"

# Export encrypted secrets
keyp export "$BACKUP_FILE" || exit 1

echo "✓ Secrets backed up to: $BACKUP_FILE"

# Keep only last 30 days of backups
find "$BACKUP_DIR" -name "secrets-*.keyp" -mtime +30 -delete
```

### Secret Rotation Script

```bash
#!/bin/bash
# rotate-secrets.sh - Rotate API keys

old_key=$(keyp get api-key --stdout)

# Call API to generate new key
new_key=$(curl -H "Authorization: Bearer $old_key" \
  https://api.example.com/rotate | jq -r '.new_key')

# Update keyp
keyp delete api-key -f
keyp set api-key "$new_key"

# Archive old key
keyp set api-key-archive-$(date +%s) "$old_key"

keyp sync push
echo "✓ API key rotated"
```

### Local Development Setup

```bash
#!/bin/bash
# setup-dev.sh - Initialize development environment

set -e

echo "Setting up development environment..."

# Initialize vault if needed
if [ ! -f ~/.keyp/vault.json ]; then
  echo "Creating vault..."
  keyp init
fi

# Store development secrets
keyp set db-host localhost
keyp set db-user dev_user
keyp set db-password changeme
keyp set api-key sk_test_123456
keyp set redis-host localhost

echo "✓ Development environment ready"
echo ""
echo "Next steps:"
echo "  keyp list              # See all secrets"
echo "  keyp get db-password   # Retrieve a secret"
echo ""
```

## Advanced Patterns

### Dynamic Configuration

```bash
#!/bin/bash
# Load config from keyp

load_config() {
  local env=${1:-dev}

  export DB_HOST=$(keyp get ${env}-db-host --stdout)
  export DB_USER=$(keyp get ${env}-db-user --stdout)
  export DB_PASS=$(keyp get ${env}-db-password --stdout)
  export API_KEY=$(keyp get ${env}-api-key --stdout)
}

# Usage
load_config prod
npm run deploy
```

### Secret Generation

```bash
#!/bin/bash
# Generate and store random secrets

generate_and_store() {
  local name=$1
  local length=${2:-32}

  secret=$(openssl rand -base64 $length | tr -d '\n')
  keyp set "$name" "$secret"

  echo "✓ Generated secret: $name"
  keyp get "$name" --stdout
}

# Usage
generate_and_store jwt-secret 64
generate_and_store api-token 32
```

### Vault Migration

```bash
#!/bin/bash
# Migrate secrets from one vault to another

migrate_vault() {
  local source_dir=$1
  local dest_dir=$2

  # Export from source
  KEYP_DIR=$source_dir keyp export backup.keyp --plain

  # Import to destination
  KEYP_DIR=$dest_dir keyp import backup.keyp --replace

  rm backup.keyp
  echo "✓ Vault migrated"
}

# Usage
migrate_vault ~/.keyp-old ~/.keyp-new
```

### Conditional Secret Retrieval

```bash
#!/bin/bash
# Get secret with fallback

get_secret_or_default() {
  local secret_name=$1
  local default_value=${2:-""}

  if keyp list | grep -q "$secret_name"; then
    keyp get "$secret_name" --stdout
  else
    echo "$default_value"
  fi
}

# Usage
DB_URL=$(get_secret_or_default "db-url" "localhost:5432")
API_KEY=$(get_secret_or_default "api-key" "dev_key_12345")
```

### Secret Expiration Alert

```bash
#!/bin/bash
# Check for secrets nearing expiration

check_expiring_secrets() {
  local vault_list=$(keyp list --stdout 2>/dev/null | grep -E ".*-expires-" || true)

  while IFS= read -r secret; do
    if [ -n "$secret" ]; then
      expiry=$(keyp get "$secret" --stdout)
      expiry_date=$(date -d "$expiry" +%s)
      today=$(date +%s)
      days_left=$(( ($expiry_date - $today) / 86400 ))

      if [ $days_left -lt 7 ]; then
        echo "⚠ Secret '$secret' expires in $days_left days"
      fi
    fi
  done <<< "$vault_list"
}

# Schedule with cron
# 0 9 * * * /path/to/check_expiring_secrets.sh
```

---

**More Examples?** Check out the [Git Sync Guide](./GIT_SYNC.md) and [CLI Reference](./CLI.md) for additional workflows!
