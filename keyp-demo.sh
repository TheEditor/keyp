#!/bin/bash
# keyp asciinema demo script
# 
# PREP (add to ~/.bashrc before recording):
#   export PS1='\[\e[32m\]keyp-demo\[\e[0m\]:\[\e[34m\]\W\[\e[0m\]$ '
#   export HOSTNAME=demo
#   export USER=user
#
# RECORD:
#   asciinema rec keyp-demo.cast -c ./keyp-demo.sh
#
# Or record interactively and run commands manually for natural pacing:
#   asciinema rec keyp-demo.cast
#
# TIPS:
# - Use a clean terminal (clear history, simple prompt)
# - Delete ~/.keyp before recording for fresh vault
# - Practice the password typing (use "demo1234" for consistency)

set -e

# Colors for comments (printed but not typed)
comment() {
    echo -e "\n\033[90m# $1\033[0m"
    sleep 1
}

# Simulate typing with delay
type_cmd() {
    echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
    for (( i=0; i<${#1}; i++ )); do
        echo -n "${1:$i:1}"
        sleep 0.05
    done
    echo
    sleep 0.3
}

# Run command (for automated version)
run() {
    type_cmd "$1"
    eval "$1"
    sleep 1
}

# Clean slate
rm -rf ~/.keyp 2>/dev/null || true
clear

comment "Initialize a new vault"
sleep 0.5

echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp init"
sleep 0.5
# Note: Can't automate password input well - do this interactively
# For scripted demo, would need expect or similar

cat << 'EOF'
Enter master password: ••••••••
Confirm password: ••••••••
✓ Vault created at ~/.keyp/vault.db
EOF
sleep 2

comment "Store a simple secret (API token)"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp set github-token"
sleep 0.5
cat << 'EOF'
Enter value: ••••••••
✓ Secret 'github-token' saved
EOF
sleep 2

comment "Store a structured secret (multiple fields)"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo 'keyp add "AT&T"'
sleep 0.5
cat << 'EOF'
Enter fields (empty label to finish):
  Label: Account PIN
  Value: ••••••••
  Label: Support PIN
  Value: ••••••••
  Label: Email
  Value: billing@example.com
  Label: 
✓ Secret 'AT&T' created with 3 field(s)
EOF
sleep 2

comment "List all secrets"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp list"
sleep 0.5
cat << 'EOF'
Secrets (2):
  AT&T
  github-token
EOF
sleep 2

comment "View a structured secret (values hidden)"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo 'keyp show "AT&T"'
sleep 0.5
cat << 'EOF'

AT&T
────────────────────
Account PIN: ********
Support PIN: ********
Email: billing@example.com

Tags: (none)
EOF
sleep 2

comment "Copy a secret to clipboard"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp get github-token"
sleep 0.5
cat << 'EOF'
✓ Copied to clipboard (clears in 45s)
EOF
sleep 2

comment "Search across all secrets"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo 'keyp search "PIN"'
sleep 0.5
cat << 'EOF'
Found 1 secret:
  AT&T (Account PIN, Support PIN)
EOF
sleep 2

comment "Add tags for organization"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo 'keyp tag add "AT&T" telecom family'
sleep 0.5
cat << 'EOF'
✓ Tags updated
EOF
sleep 1

echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp list"
sleep 0.5
cat << 'EOF'
Secrets (2):
  AT&T [telecom, family]
  github-token
EOF
sleep 2

comment "Lock the vault when done"
echo -ne "\033[32mkeyp-demo\033[0m:\033[34m~\033[0m$ "
echo "keyp lock"
sleep 0.5
cat << 'EOF'
✓ Vault locked
EOF
sleep 2

echo
echo -e "\033[90m# That's keyp — local-first secrets for developers and families\033[0m"
echo -e "\033[90m# https://github.com/TheEditor/keyp\033[0m"
sleep 3
