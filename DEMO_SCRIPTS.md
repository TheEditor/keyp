# Demo Scripts for GIF/Screenshot Recording

These scripts simulate keyp workflows with demo credentials so you can safely record GIFs and screenshots without exposing any real passwords or sensitive information.

## Quick Start

Each script creates temporary demo vaults and uses fake credentials. Nothing is exposed in the output.

### Demo 1: Simple `keyp init` (Best for First GIF)

```bash
chmod +x demo-init.sh
./demo-init.sh
```

**What it shows:**
- `keyp init` command execution
- Password input prompts (masked with ●●●●●●)
- Password strength meter
- Success message

**Good for:** Quick demo of vault initialization

---

### Demo 2: Full Workflow (Comprehensive Demo)

```bash
chmod +x demo-full-workflow.sh
./demo-full-workflow.sh
```

**What it shows:**
- ✅ Initialize vault (`keyp init`)
- ✅ Store secrets (`keyp set`)
- ✅ List secrets (`keyp list`)
- ✅ Retrieve secret (`keyp get`)
- ✅ View statistics (`keyp stats`)

**Good for:** Complete feature walkthrough GIF (3-5 minutes)

---

### Demo 3: Git Sync Workflow

```bash
chmod +x demo-git-sync.sh
./demo-git-sync.sh
```

**What it shows:**
- ✅ Initialize vault
- ✅ Store secrets
- ✅ Initialize Git sync (`keyp sync init`)
- ✅ Check status (`keyp sync status`)
- ✅ Push to Git (`keyp sync push`)

**Good for:** Multi-machine sync and backup features

---

## How to Record GIFs

### Using asciinema (Recommended)

```bash
# Install asciinema if needed
npm install -g asciinema

# Record the demo
asciinema rec demo-keyp-init.cast
# Run: ./demo-init.sh
# Press Ctrl+D to stop recording

# Convert to GIF
npm install -g asciicast2gif
asciicast2gif demo-keyp-init.cast demo-keyp-init.gif
```

### Using Built-in Terminal Recording (macOS)

```bash
# Record with QuickTime
# File > New Screen Recording
# Run: ./demo-init.sh
# Convert MP4 to GIF with ffmpeg
ffmpeg -i recording.mp4 -vf "fps=10,scale=1280:-1:flags=lanczos" demo.gif
```

### Using Linux Tools

```bash
# Using SimpleScreenRecorder or OBS
# Record terminal while running: ./demo-init.sh
# Export as MP4, then convert to GIF

ffmpeg -i recording.mp4 -vf "fps=10,scale=1280:-1:flags=lanczos" demo.gif
```

---

## Demo Credentials

All scripts use these **fake demo credentials** (safe to show in recordings):

- **Master Password:** `SecureDemo123!`
- **GitHub Token:** `gh_1234567890abcdefghij`
- **API Key:** `sk_live_abc123xyz`
- **DB Password:** `db_secure_pass_123`

These are clearly fake and demonstrate the masked input (●●●●●●) without exposing real secrets.

---

## Important Notes

1. **Temporary Vaults:** Each script creates temporary directories that are deleted after running
2. **No Real Credentials:** Only demo/fake credentials are used
3. **Masked Input:** Terminal shows ●●●●●● for all password inputs (no actual characters visible)
4. **Safe to Share:** All output is safe to share in public GIFs/videos
5. **Customizable:** Edit the scripts to use different demo values if desired

---

## Tips for Better GIFs

1. **Clean Terminal:**
   ```bash
   PS1='$ '  # Clean prompt with no user/path info
   ```

2. **Appropriate Speed:** Run scripts at normal pace (not too fast)

3. **Clear Output:** Use terminal with good contrast (dark background recommended)

4. **Capture Just Output:** Show only relevant portions of the terminal

5. **Multiple Takes:** Don't worry if first attempt isn't perfect, script is repeatable

---

## File Locations

- `demo-init.sh` - Simple init demo
- `demo-full-workflow.sh` - Complete feature demo
- `demo-git-sync.sh` - Git sync demo
- `DEMO_SCRIPTS.md` - This file

---

## Examples

### Recording keyp init demo:

```bash
# 1. Open terminal with clean prompt
PS1='$ '

# 2. Start recording with asciinema
asciinema rec demo-init-recording.cast

# 3. Run the demo
./demo-init.sh

# 4. Stop recording (Ctrl+D)

# 5. Convert to GIF
asciicast2gif demo-init-recording.cast assets/demo-init.gif

# 6. Embed in README
# ![keyp init demo](assets/demo-init.gif)
```

---

**Remember:** These scripts are safe for recording because they use fake, clearly-fake credentials that are obviously for demonstration purposes. No real passwords or secrets are ever exposed.
