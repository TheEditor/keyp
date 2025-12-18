# Guideline: Detail Level for `bd create` Descriptions

*Constructed by Claude Opus 4.5 for Haiku models executing task specifications*

---

## Purpose

When writing task specs for Haiku to execute, the detail level in `bd create` descriptions directly impacts execution efficiency. Too little detail forces Haiku to investigate and iterate (burning tokens, risking wrong solutions). Too much detail on obvious tasks wastes spec-writing time.

This guide helps determine how much detail to embed in issue descriptions.

---

## The Core Principle

**Trade diagnostic time for execution time.**

The spec-writing model (Opus) runs once to diagnose and prescribe. The executing model (Haiku) runs once to implement. If Haiku has to diagnose *and* implement, it may run multiple iterations -- or fail entirely on subtle fixes.

---

## Subtlety Spectrum

| Level | Characteristics | Description Needs | Example |
|-------|-----------------|-------------------|---------|
| **Obvious** | What implies how; single clear approach | Problem + location | "Add `disabled` prop to button during submit" |
| **Moderate** | Multiple valid approaches; some ambiguity | Problem + location + preferred approach | "Prevent double-submission -- use loading state flag, not debounce" |
| **Subtle** | Fix requires non-obvious knowledge; symptom != cause | Problem + location + exact code change | "Flex child won't truncate -- add `min-w-0` to parent" |
| **Arcane** | Browser quirks, obscure APIs, edge cases | Problem + location + exact code + explanation of why | "Safari grid gap inheritance bug -- requires explicit `gap: inherit` on nested grid" |

---

## Decision Framework

Ask these questions:

### 1. Does the symptom point to the cause?

- **Yes** -> Less detail needed
- **No** -> Include diagnosis in description

*Example*: "Button is blue, should be green" -- symptom is the cause. "Text bleeds outside card" -- symptom (overflow) doesn't reveal cause (flex min-width default).

### 2. Is there one obvious implementation?

- **Yes** -> Describe the what, Haiku handles the how
- **No** -> Specify the preferred approach

*Example*: "Add loading spinner" -- obvious. "Improve form performance" -- ambiguous, many approaches.

### 3. Does the fix require knowing a trick or quirk?

- **Yes** -> Provide the exact code
- **No** -> Description sufficient

*Example*: `min-w-0` for flex truncation is a trick. Adding a CSS class is not.

### 4. Would this require research to solve?

- **No** -> Light description
- **Yes, common pattern** -> Moderate description with hints
- **Yes, obscure solution** -> Exact code

---

## Task Granularity

### How Big Should a Task Be?

**One task = one verifiable unit of work.**

Guidelines:
- Can be tested/verified independently
- Description fits comfortably in a `bd create -d` field (under ~2000 chars)
- Completes in under 30 minutes of implementation time
- Touches a coherent set of related code

**Split when:**
- Changes are in unrelated areas of the codebase
- Each change can be verified independently
- Description becomes unwieldy (multiple scrolls to read)
- Different team members could work on pieces in parallel

**Combine when:**
- Changes are tightly coupled (type definition + usage)
- One change is meaningless without the other
- Verification requires all pieces together

### Multi-File Changes

When a logical change spans multiple files, decide: **one task or many?**

**Single task (tightly coupled):**
- Adding a type in `types.ts` and immediately using it in `component.tsx`
- Creating a utility function and calling it in the same PR
- Changes that would break the build if done partially

Format for multi-file single task:
```
PROBLEM: [What's wrong]

FILE 1: app/types.ts (around line 15)
FIND:
[code]

CHANGE TO:
[code]

FILE 2: app/component.tsx (around line 42)
FIND:
[code]

CHANGE TO:
[code]

ACCEPTANCE: [How to verify all changes work together]
```

**Multiple tasks with dependencies (loosely coupled):**
- Refactoring that touches many consumers
- Changes that can be verified file-by-file
- Work that could be parallelized

When splitting, always create dependencies:
```bash
# Task 1: Foundation work
bd create "Add UserProfile type" -t task -p 2 --parent <EPIC_ID> -d "..." --json
# Returns: {"id":"proj-042",...}

# Task 2: Depends on Task 1
bd create "Update ProfileCard to use UserProfile type" -t task -p 2 --parent <EPIC_ID> -d "..." --json
# Returns: {"id":"proj-043",...}

# Register the dependency
bd dep add proj-043 proj-042 --type depends-on
```

---

## Priority Levels

| Priority | When to Use | Examples |
|----------|-------------|----------|
| **P0** | Blocking issue, production broken, data loss risk | Build failure, security vulnerability, crash |
| **P1** | Important, affects user experience, should do soon | Visible bugs, broken features, performance issues |
| **P2** | Normal work, planned improvements | New features, refactoring, most tasks |
| **P3** | Nice to have, low urgency | Code cleanup, minor polish, documentation |

**Default to P2** for most planned work. Reserve P0/P1 for genuine urgency.

---

## Dependencies

### When to Use `--parent`

Always use `--parent <EPIC_ID>` when:
- Task is part of a larger initiative
- Related tasks should be grouped
- You want to close all tasks before closing the epic

### When to Use `bd dep add`

Use `bd dep add <blocked> <blocker> --type <type>` when:
- Task B cannot start until Task A completes
- Task B uses something Task A creates
- Order of execution matters

Dependency types:
- `depends-on` -- Task cannot start until blocker is done
- `blocks` -- Inverse of depends-on
- `relates-to` -- Informational link, no blocking
- `discovered-from` -- New work found while doing another task

---

## Line Number Precision

When specifying locations:
- **"around line X"** means +/- 20 lines is acceptable
- **"line X"** means exact line (use sparingly, lines shift)
- Always include enough context in FIND block to locate uniquely

**Good:** "around line 180" with a unique 3-line FIND block
**Bad:** "line 180" with a generic single-line FIND block

If the file is small (<100 lines), line numbers are optional -- the FIND block is enough.

---

## Escape Character Rules

In `bd create -d "..."` descriptions, escape:
- Double quotes: `\"` 
- Backticks in template literals: `` \` ``
- Dollar signs in template literals: `\$`
- Backslashes: `\\`

Example:
```bash
bd create "Fix template" -t task -p 2 -d "FIND:
title={\`PROJECT: \${prompt.project}\`}

CHANGE TO:
title={\`Project: \${prompt.project}\`}"
```

When in doubt, use single quotes for the outer wrapper if your shell supports it, or put the description in a file and use `bd create -d @file.txt`.

---

## Description Templates

### Obvious Fix
```
PROBLEM: [What's wrong]
LOCATION: [File and approximate line]
CHANGE: [What to do in plain English]
```

### Moderate Fix
```
PROBLEM: [What's wrong]
LOCATION: [File and approximate line]  
APPROACH: [Preferred solution path]
ACCEPTANCE: [How to verify it works]
```

### Subtle Fix
```
PROBLEM: [What's wrong]
CAUSE: [Why it's happening -- the diagnosis]
LOCATION: [File and approximate line]

FIND:
[Exact current code]

CHANGE TO:
[Exact target code]

WHY: [Brief explanation of the fix]
ACCEPTANCE: [How to verify]
```

### Arcane Fix
```
PROBLEM: [What's wrong]
CAUSE: [Technical explanation of the underlying issue]
LOCATION: [File and approximate line]

FIND:
[Exact current code]

CHANGE TO:
[Exact target code]

TECHNICAL NOTE: [Why this specific fix works; what alternatives don't work and why]
ACCEPTANCE: [How to verify]
REGRESSION CHECK: [What else might break]
```

---

## Examples by Category

### Obvious -- No Code Needed
```bash
bd create "Add loading state to save button" -t task -p 2 -d "PROBLEM: Save button can be clicked multiple times during submission.
LOCATION: app/dashboard/prompt-editor.tsx, save button around line 180
CHANGE: Add isSaving state. Disable button and show spinner while saving.
ACCEPTANCE: Button shows spinner and is disabled during save operation."
```

### Moderate -- Approach Specified
```bash
bd create "Debounce search input" -t task -p 2 -d "PROBLEM: Search fires on every keystroke, causing excessive API calls.
LOCATION: app/dashboard/dashboard-client.tsx, search input around line 320
APPROACH: Use 300ms debounce on setSearchQuery. Don't add lodash -- use setTimeout/clearTimeout pattern.
ACCEPTANCE: Search only fires 300ms after user stops typing. No new dependencies."
```

### Subtle -- Exact Code Provided
```bash
bd create "Fix project badge overflow" -t task -p 1 -d "PROBLEM: Project badge text extends beyond card boundary.
CAUSE: Badge has no width constraint; flex-wrap doesn't limit individual item width.
LOCATION: app/dashboard/dashboard-client.tsx around line 475

FIND:
<Badge variant=\"outline\" className=\"border-2 border-secondary/30 text-secondary font-bold text-xs tracking-wide\">
  PROJECT: {prompt.project}
</Badge>

CHANGE TO:
<Badge 
  variant=\"outline\" 
  className=\"border-2 border-secondary/30 text-secondary font-bold text-xs tracking-wide max-w-full truncate\"
  title={\`PROJECT: \${prompt.project}\`}
>
  PROJECT: {prompt.project}
</Badge>

WHY: max-w-full constrains to parent; truncate adds ellipsis; title provides hover tooltip for full text.
ACCEPTANCE: Long project names truncate with ellipsis. Hover shows full name."
```

### Arcane -- Full Explanation
```bash
bd create "Fix flex container preventing text truncation" -t task -p 1 -d "PROBLEM: SelectValue text overflows 220px container despite truncate class.
CAUSE: Flex containers have implicit min-width:auto which prevents children from shrinking below content size. The truncate class on the child can't work because the parent won't let it shrink.
LOCATION: app/dashboard/dashboard-client.tsx around line 335

FIND:
<SelectTrigger className=\"w-full sm:w-[220px] !h-12 bg-input ...\">
  <FolderKanban className=\"w-4 h-4 mr-2 text-muted-foreground\" />

CHANGE TO:
<SelectTrigger className=\"w-full sm:w-[220px] !h-12 bg-input ... min-w-0\">
  <FolderKanban className=\"w-4 h-4 mr-2 shrink-0 text-muted-foreground\" />

TECHNICAL NOTE: min-w-0 overrides the flex default, allowing the container to shrink. shrink-0 on the icon prevents it from shrinking when the text truncates. Without both, either the text won't truncate or the icon will compress.
ACCEPTANCE: Long project names truncate. Icon stays full size. Dropdown arrow visible."
```

---

## What NOT to Do

### Bad: Vague Description
```bash
# DON'T
bd create "Fix the card" -t task -p 2 -d "The card looks wrong. Please fix it."
```
Problem: No location, no specifics, no acceptance criteria.

### Bad: Missing Context
```bash
# DON'T
bd create "Add min-w-0" -t task -p 2 -d "Add min-w-0 to the div on line 335."
```
Problem: No explanation of why. If line numbers shift, Haiku can't find it.

### Bad: Over-specified Obvious Task
```bash
# DON'T
bd create "Change button color" -t task -p 2 -d "PROBLEM: Button is blue.
CAUSE: The className uses bg-blue-500.
LOCATION: app/page.tsx line 42

FIND:
<Button className=\"bg-blue-500\">

CHANGE TO:
<Button className=\"bg-green-500\">

TECHNICAL NOTE: CSS background colors are applied via Tailwind utility classes. The bg-* prefix sets background-color. Green is chosen for semantic meaning of success...

ACCEPTANCE: Button is green."
```
Problem: Massively over-detailed for a trivial change. Wastes spec-writing time.

### Bad: Monolithic Multi-File Task
```bash
# DON'T
bd create "Refactor user system" -t task -p 2 -d "Update types.ts, then user-service.ts, then profile.tsx, then settings.tsx, then header.tsx, then sidebar.tsx, then api/user/route.ts..."
```
Problem: Too many files, impossible to verify incrementally, huge description.

---

## When Things Go Wrong: Abort vs Adapt

### Haiku's Decision Tree When FIND Block Doesn't Match

**1. Minor drift (whitespace, formatting):**
- ADAPT: Make the equivalent change to the actual code
- Proceed with task

**2. Code restructured but intent is clear:**
- ADAPT: Apply the spirit of the change to new structure
- Note the adaptation in commit message

**3. Code completely different or missing:**
- ABORT: Do not guess
- Comment on why the task couldn't be completed
- Leave task open for human review

**4. New work discovered during implementation:**
- Complete current task if possible
- File new issue immediately:
  ```bash
  bd create "Discovered: [problem]" -t bug -p 1 --json
  # Capture returned ID
  bd dep add <new-issue-id> <current-task-id> --type discovered-from
  ```
- Note discovery in current task's completion message

### Examples

**Adapt (OK):**
```
Spec said line 335, found equivalent code at line 342.
Applied same change. Lines shifted due to earlier commits.
```

**Abort (correct):**
```
FIND block not found. The SelectTrigger component appears to have 
been replaced with a custom component. Cannot safely apply change.
Task left open for review.
```

**Discovered work:**
```
While implementing card spacing fix, found that the Badge component 
has hardcoded styles preventing truncation. Filed proj-047 as 
discovered-from this task.
```

---

## When In Doubt

**Over-specify.** 

The cost of Haiku iterating (or failing) exceeds the cost of writing detailed descriptions. A few extra minutes writing the spec saves token budget and guarantees correct implementation.

---

## Summary

| If the fix is... | Then provide... |
|------------------|-----------------|
| Obvious | Problem + location |
| Moderate | Problem + location + approach |
| Subtle | Problem + location + exact code |
| Arcane | Problem + location + exact code + technical explanation |

| If changes span... | Then... |
|--------------------|---------|
| Single file | One task |
| Multiple files, tightly coupled | One task, multiple FILE sections |
| Multiple files, loosely coupled | Multiple tasks with `bd dep add` |

The goal: **Haiku executes once, correctly.**
