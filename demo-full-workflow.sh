#!/bin/bash
# Complete keyp demo workflow script
# Demonstrates: init → set → list → get → stats
# All with fake demo credentials (no real passwords exposed)

set -e  # Exit on error

DEMO_PASSWORD="SecureDemo123!"
DEMO_DIR=$(mktemp -d)
export KEYP_HOME="$DEMO_DIR"

echo "════════════════════════════════════════"
echo "   keyp Demo Workflow"
echo "════════════════════════════════════════"
echo ""

# ============================================
# Demo 1: Initialize Vault
# ============================================
echo "📋 Step 1: Initialize Vault"
echo "Command: keyp init"
echo ""
echo -e "$DEMO_PASSWORD\n$DEMO_PASSWORD" | keyp init
echo ""
sleep 2

# ============================================
# Demo 2: Store First Secret
# ============================================
echo "📋 Step 2: Store a Secret (github-token)"
echo "Command: keyp set github-token"
echo ""
echo -e "$DEMO_PASSWORD\ngh_1234567890abcdefghij\n$DEMO_PASSWORD" | keyp set github-token
echo ""
sleep 2

# ============================================
# Demo 3: Store More Secrets
# ============================================
echo "📋 Step 3: Store More Secrets"
echo ""
echo "Command: keyp set api-key"
echo -e "$DEMO_PASSWORD\nsk_live_abc123xyz\n$DEMO_PASSWORD" | keyp set api-key
echo ""
echo "Command: keyp set db-password"
echo -e "$DEMO_PASSWORD\ndb_secure_pass_123\n$DEMO_PASSWORD" | keyp set db-password
echo ""
sleep 2

# ============================================
# Demo 4: List All Secrets
# ============================================
echo "📋 Step 4: List All Secrets"
echo "Command: keyp list"
echo ""
echo -e "$DEMO_PASSWORD" | keyp list
echo ""
sleep 2

# ============================================
# Demo 5: Get a Secret (Clipboard)
# ============================================
echo "📋 Step 5: Retrieve Secret to Clipboard"
echo "Command: keyp get github-token"
echo ""
echo -e "$DEMO_PASSWORD" | keyp get github-token
echo ""
sleep 2

# ============================================
# Demo 6: Vault Statistics
# ============================================
echo "📋 Step 6: View Vault Statistics"
echo "Command: keyp stats"
echo ""
echo -e "$DEMO_PASSWORD" | keyp stats
echo ""

# ============================================
# Cleanup
# ============================================
echo "════════════════════════════════════════"
echo "Demo Complete!"
echo "════════════════════════════════════════"
echo ""
echo "To repeat this demo or create variations:"
echo "  export KEYP_HOME=$DEMO_DIR"
echo ""
echo "To clean up:"
echo "  rm -rf $DEMO_DIR"
echo ""
