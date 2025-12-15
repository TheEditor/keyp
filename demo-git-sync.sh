#!/bin/bash
# keyp Git Sync Demo Script
# Demonstrates: sync init → sync status → sync push
# Creates a local git repository for safe demo purposes

set -e

DEMO_PASSWORD="SecureDemo123!"
DEMO_DIR=$(mktemp -d)
DEMO_GIT_REPO=$(mktemp -d)
export KEYP_HOME="$DEMO_DIR"

echo "════════════════════════════════════════"
echo "   keyp Git Sync Demo"
echo "════════════════════════════════════════"
echo ""

# ============================================
# Step 1: Initialize Vault
# ============================================
echo "📋 Step 1: Initialize Vault"
echo "Command: keyp init"
echo ""
echo -e "$DEMO_PASSWORD\n$DEMO_PASSWORD" | keyp init
echo ""
sleep 2

# ============================================
# Step 2: Add Some Secrets
# ============================================
echo "📋 Step 2: Add Secrets to Vault"
echo ""
echo -e "$DEMO_PASSWORD\ngh_demo_token_123\n$DEMO_PASSWORD" | keyp set github-token
echo -e "$DEMO_PASSWORD\ndb_pass_secure\n$DEMO_PASSWORD" | keyp set db-password
echo ""
sleep 2

# ============================================
# Step 3: Initialize Git Sync
# ============================================
echo "📋 Step 3: Initialize Git Sync"
echo "Command: keyp sync init <repo-url>"
echo ""
# Create a bare git repo for demo
cd "$DEMO_GIT_REPO"
git init --bare keyp-backup.git
cd - > /dev/null

GIT_REPO_PATH="$DEMO_GIT_REPO/keyp-backup.git"
echo -e "$DEMO_PASSWORD" | keyp sync init "file://$GIT_REPO_PATH"
echo ""
sleep 2

# ============================================
# Step 4: Check Sync Status
# ============================================
echo "📋 Step 4: Check Sync Status"
echo "Command: keyp sync status"
echo ""
echo -e "$DEMO_PASSWORD" | keyp sync status
echo ""
sleep 2

# ============================================
# Step 5: Push Vault to Git
# ============================================
echo "📋 Step 5: Push Vault to Git"
echo "Command: keyp sync push"
echo ""
echo -e "$DEMO_PASSWORD" | keyp sync push
echo ""
sleep 2

# ============================================
# Step 6: Verify Status
# ============================================
echo "📋 Step 6: Verify Sync Status"
echo "Command: keyp sync status"
echo ""
echo -e "$DEMO_PASSWORD" | keyp sync status
echo ""

# ============================================
# Cleanup
# ============================================
echo "════════════════════════════════════════"
echo "Git Sync Demo Complete!"
echo "════════════════════════════════════════"
echo ""
echo "Demo repositories:"
echo "  Vault: $DEMO_DIR"
echo "  Git backup: $GIT_REPO_PATH"
echo ""
echo "To clean up:"
echo "  rm -rf $DEMO_DIR $DEMO_GIT_REPO"
echo ""
