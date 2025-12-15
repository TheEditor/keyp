#!/bin/bash
# Demo script for keyp init GIF
# This script simulates the keyp init flow with demo credentials
# No real passwords are used or exposed in the output

# Demo password (fake, just for the demo)
DEMO_PASSWORD="SecureDemo123!"

# Create a temporary demo vault directory
DEMO_DIR=$(mktemp -d)
export KEYP_HOME="$DEMO_DIR"

echo "Running keyp init demo..."
echo "Demo vault location: $DEMO_DIR"
echo ""

# Run keyp init with piped input (password masked in output)
echo -e "$DEMO_PASSWORD\n$DEMO_PASSWORD" | keyp init

# Keep the demo vault available for follow-up demos
echo ""
echo "Demo vault ready for additional commands!"
echo "To use this vault for other demos, set:"
echo "  export KEYP_HOME=$DEMO_DIR"
echo ""
echo "Or clean up with:"
echo "  rm -rf $DEMO_DIR"
