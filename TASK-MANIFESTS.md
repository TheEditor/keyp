# Task Manifests

A task manifest is a compact document that references Beads issues without embedding their full descriptions. It serves as a "customs declaration" for a phase of work - listing what's in the shipment and their relationships, while the full cargo (descriptions) lives in the Beads database.

## Why Manifests?

Traditional task specs embed full descriptions in Markdown, causing:
- **File size bloat** - 30k+ files strain AI context windows
- **Duplication** - Descriptions exist in spec AND Beads database
- **Staleness** - Spec can drift from actual issue content

Manifests solve this by keeping specs tiny. Full context lives in Beads where it belongs.

## Manifest Format

```markdown
# [Project] Phase [N]: [Title]

## Preface
[Brief context - 3-5 lines max]
- Read bd-issue-tracking skill first
- Run `bd ready --json` to see available work
- Use `bd show <id>` to read full descriptions

## Issues
<id>: <type> - <title>
<id>: <title> [depends: <id>, <id>]
<id>: <title> [depends: <id>]
...

## Completion
- All issues closed: `bd ready --json` returns empty
- Tests pass: `make test`
- Final commit includes `(bd:<epic-id>)`
```

## Example Manifest

```markdown
# keyp Phase 3: CLI Commands

## Preface
Implement core CLI commands. Phase 2 (storage layer) must be complete.
- Read bd-issue-tracking skill: `view /mnt/skills/user/bd-issue-tracking/SKILL.md`
- Check ready work: `bd ready --json`
- View issue details: `bd show <id>`

## Issues
keyp-a1b2: epic - Phase 3 CLI Commands
keyp-c3d4: Add password prompt utility
keyp-e5f6: Add clipboard utility
keyp-g7h8: Implement keyp init [depends: keyp-c3d4]
keyp-i9j0: Implement keyp set [depends: keyp-c3d4, keyp-g7h8]
keyp-k1l2: Implement keyp get [depends: keyp-e5f6, keyp-i9j0]
keyp-m3n4: Implement keyp list [depends: keyp-g7h8]
keyp-o5p6: Implement keyp delete [depends: keyp-c3d4, keyp-g7h8]

## Completion
- `bd ready --json` returns empty
- `make build && make test` passes
- `./keyp init && ./keyp set test val && ./keyp get test` works
```

## Workflow

### Planning Session (Human + AI)

1. **Discuss scope** - What's in this phase?
2. **AI creates issues** - Runs `bd create` with full descriptions
3. **AI captures IDs** - Beads assigns them automatically
4. **AI sets dependencies** - Runs `bd dep add`
5. **AI builds manifest** - Compact subject lines with IDs

### Execution Session (Haiku or other agent)

1. **Read manifest** - Understand scope and dependencies
2. **Run `bd ready`** - See what's unblocked
3. **Run `bd show <id>`** - Get full description for current task
4. **Implement** - Do the work
5. **Close issue** - `bd close <id> --reason "..."`
6. **Repeat** - Until `bd ready` returns empty

### Visualization (Human)

Use Jeffrey Emanuel's `bv` (Beads Viewer) for rich visualization:
- `bv` - Interactive TUI with dependency graphs
- `bv --robot-plan` - JSON execution plan
- `bv --robot-insights` - Graph metrics (PageRank, bottlenecks)

## Key Principles

1. **Descriptions live in Beads** - Not in the manifest
2. **Manifest is just references** - IDs, titles, dependencies
3. **One line per issue** - Keep it scannable
4. **Dependencies inline** - `[depends: id1, id2]` format
5. **Preface is minimal** - Just enough to orient the agent

## Scaling

A manifest at 30k characters can reference 200+ issues comfortably. The actual descriptions (which can be arbitrarily long) live in Beads and are fetched on-demand via `bd show`.

This enables planning sessions that scope entire projects (100+ issues) in a single conversation, then hand off to execution agents in focused chunks.

## Relationship to Other Tools

| Tool | Purpose |
|------|---------|
| `bd` | CLI for creating/managing issues |
| `bv` | TUI for visualizing issues and dependencies |
| Manifest | Compact reference document for AI handoff |
| AGENTS.md | Project-level instructions and context |

The manifest complements these tools - it's the bridge between planning (where issues are created) and execution (where issues are worked).
