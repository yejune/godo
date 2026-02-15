# Worktree Commands Module

Purpose: Complete CLI command reference for Git worktree management with detailed usage examples and advanced options.

Version: 2.0.0
Last Updated: 2026-01-06

---

## Quick Reference (30 seconds)

Command Categories:
- Creation: new - Create isolated worktree
- Navigation: list, switch, go - Browse and navigate
- Management: sync, remove, clean - Maintain worktrees
- Status: status - Check worktree state
- Configuration: config - Manage settings

Quick Start:
1. Create worktree: {{slot:BRAND}}-worktree new SPEC-001 "User Authentication"
2. Switch to worktree: {{slot:BRAND}}-worktree switch SPEC-001
3. Or use shell eval: eval $({{slot:BRAND}}-worktree go SPEC-001)

---

## Creation Commands

### {{slot:BRAND}}-worktree new - Create Worktree

Create a new isolated Git worktree for SPEC development.

Syntax: {{slot:BRAND}}-worktree new <spec-id> [description] [options]

Arguments:
- spec-id: SPEC identifier (e.g., SPEC-001, SPEC-AUTH-001)
- description: Optional description for the worktree

Options:
- --branch <name>: Create specific branch instead of auto-generated
- --base <branch>: Base branch for new worktree (default: main)
- --template <name>: Use predefined template
- --shallow: Create shallow clone for faster setup
- --depth <number>: Clone depth for shallow clone
- --force: Force creation even if worktree exists

Examples:
- Basic creation: {{slot:BRAND}}-worktree new SPEC-001 "User Auth System"
- Custom branch: {{slot:BRAND}}-worktree new SPEC-002 "Payment" --branch feature/payment-gateway
- From develop: {{slot:BRAND}}-worktree new SPEC-003 "API Refactor" --base develop
- With template: {{slot:BRAND}}-worktree new SPEC-004 "Frontend" --template frontend
- Fast creation: {{slot:BRAND}}-worktree new SPEC-005 "Bug Fixes" --shallow --depth 1

Auto-Generated Branch Pattern:
- Format: feature/SPEC-{ID}-{description-kebab-case}
- Example: SPEC-001 becomes feature/SPEC-001-user-authentication

---

## Navigation Commands

### {{slot:BRAND}}-worktree list - List Worktrees

Display all registered worktrees with their status and metadata.

Syntax: {{slot:BRAND}}-worktree list [options]

Options:
- --format <format>: Output format (table, json, csv)
- --status <status>: Filter by status (active, merged, stale)
- --sort <field>: Sort by field (name, created, modified, status)
- --reverse: Reverse sort order
- --verbose: Show detailed information

Examples:
- Table format: {{slot:BRAND}}-worktree list
- JSON output: {{slot:BRAND}}-worktree list --format json
- Active only: {{slot:BRAND}}-worktree list --status active
- Sort by date: {{slot:BRAND}}-worktree list --sort created
- Detailed: {{slot:BRAND}}-worktree list --verbose

### {{slot:BRAND}}-worktree switch - Switch to Worktree

Change current working directory to the specified worktree.

Syntax: {{slot:BRAND}}-worktree switch <spec-id> [options]

Options:
- --auto-sync: Automatically sync before switching
- --force: Force switch even with uncommitted changes
- --new-terminal: Open in new terminal window

Examples:
- Basic switch: {{slot:BRAND}}-worktree switch SPEC-001
- With sync: {{slot:BRAND}}-worktree switch SPEC-002 --auto-sync
- Force switch: {{slot:BRAND}}-worktree switch SPEC-003 --force

### {{slot:BRAND}}-worktree go - Get Worktree Path

Output the cd command for shell integration.

Syntax: {{slot:BRAND}}-worktree go <spec-id> [options]

Options:
- --absolute: Show absolute path
- --relative: Show relative path from current directory
- --export: Export as environment variable

Shell Integration Methods:
- eval pattern (recommended): eval $({{slot:BRAND}}-worktree go SPEC-001)
- source pattern: {{slot:BRAND}}-worktree go SPEC-001 | source
- manual cd: cd $({{slot:BRAND}}-worktree go SPEC-001 --absolute)

---

## Management Commands

### {{slot:BRAND}}-worktree sync - Synchronize Worktree

Synchronize worktree with its base branch.

Syntax: {{slot:BRAND}}-worktree sync <spec-id> [options]

Arguments:
- spec-id: Worktree identifier (or --all for all worktrees)

Options:
- --auto-resolve: Automatically resolve simple conflicts
- --interactive: Interactive conflict resolution
- --dry-run: Show what would be synced without doing it
- --force: Force sync even with uncommitted changes
- --include <pattern>: Include only specific files
- --exclude <pattern>: Exclude specific files

Examples:
- Sync specific: {{slot:BRAND}}-worktree sync SPEC-001
- Sync all: {{slot:BRAND}}-worktree sync --all
- Interactive: {{slot:BRAND}}-worktree sync SPEC-001 --interactive
- Preview: {{slot:BRAND}}-worktree sync SPEC-001 --dry-run
- Include pattern: {{slot:BRAND}}-worktree sync SPEC-001 --include "src/"
- Exclude pattern: {{slot:BRAND}}-worktree sync SPEC-001 --exclude "node_modules/"

Conflict Resolution:
When conflicts detected, choose from:
1. Keep worktree version
2. Accept base branch version
3. Open merge tool
4. Skip file
5. Abort sync

### {{slot:BRAND}}-worktree remove - Remove Worktree

Remove a worktree and clean up its registration.

Syntax: {{slot:BRAND}}-worktree remove <spec-id> [options]

Options:
- --force: Force removal without confirmation
- --keep-branch: Keep the branch after removing worktree
- --backup: Create backup before removal
- --dry-run: Show what would be removed without doing it

Examples:
- Interactive: {{slot:BRAND}}-worktree remove SPEC-001
- Force: {{slot:BRAND}}-worktree remove SPEC-001 --force
- Keep branch: {{slot:BRAND}}-worktree remove SPEC-001 --keep-branch
- With backup: {{slot:BRAND}}-worktree remove SPEC-001 --backup
- Preview: {{slot:BRAND}}-worktree remove SPEC-001 --dry-run

### {{slot:BRAND}}-worktree clean - Clean Up Worktrees

Remove worktrees for merged branches or stale worktrees.

Syntax: {{slot:BRAND}}-worktree clean [options]

Options:
- --merged-only: Only remove worktrees with merged branches
- --stale: Remove worktrees not updated in specified days
- --days <number>: Stale threshold in days (default: 30)
- --interactive: Interactive selection of worktrees to remove
- --dry-run: Show what would be cleaned without doing it
- --force: Skip confirmation prompts

Examples:
- Merged only: {{slot:BRAND}}-worktree clean --merged-only
- Stale (30 days): {{slot:BRAND}}-worktree clean --stale
- Custom threshold: {{slot:BRAND}}-worktree clean --stale --days 14
- Interactive: {{slot:BRAND}}-worktree clean --interactive
- Preview: {{slot:BRAND}}-worktree clean --dry-run
- Force: {{slot:BRAND}}-worktree clean --force

---

## Status and Configuration

### {{slot:BRAND}}-worktree status - Show Worktree Status

Display detailed status information about worktrees.

Syntax: {{slot:BRAND}}-worktree status [spec-id] [options]

Arguments:
- spec-id: Specific worktree (optional, shows current if not specified)

Options:
- --all: Show status of all worktrees
- --sync-check: Check if worktrees need sync
- --detailed: Show detailed Git status
- --format <format>: Output format (table, json)

Examples:
- Current worktree: {{slot:BRAND}}-worktree status
- Specific worktree: {{slot:BRAND}}-worktree status SPEC-001
- All with sync check: {{slot:BRAND}}-worktree status --all --sync-check
- Detailed Git status: {{slot:BRAND}}-worktree status SPEC-001 --detailed
- JSON output: {{slot:BRAND}}-worktree status --all --format json

Status Output Includes:
- Worktree path and branch
- Commits ahead/behind base
- Modified and untracked files
- Sync status and last sync time

### {{slot:BRAND}}-worktree config - Configuration Management

Manage {{slot:BRAND}}-worktree configuration settings.

Syntax: {{slot:BRAND}}-worktree config <action> [key] [value]

Actions:
- get [key]: Get configuration value
- set <key> <value>: Set configuration value
- list: List all configuration
- reset [key]: Reset to default value
- edit: Open configuration in editor

Configuration Keys:
- worktree_root: Root directory for worktrees
- auto_sync: Enable automatic sync (true/false)
- cleanup_merged: Auto-cleanup merged worktrees (true/false)
- default_base: Default base branch (main/develop)
- template_dir: Directory for worktree templates
- sync_strategy: Sync strategy (merge, rebase, squash)

Examples:
- List all: {{slot:BRAND}}-worktree config list
- Get value: {{slot:BRAND}}-worktree config get worktree_root
- Set value: {{slot:BRAND}}-worktree config set auto_sync true
- Reset: {{slot:BRAND}}-worktree config reset worktree_root
- Edit: {{slot:BRAND}}-worktree config edit

---

## Advanced Usage

### Batch Operations

Sync all active worktrees:
- Use shell loop with list --format json and jq to extract IDs
- Run sync for each ID in sequence or parallel

Clean all merged worktrees:
- {{slot:BRAND}}-worktree clean --merged-only --force

Create worktrees from SPEC list:
- Read SPEC IDs from file
- Run new command for each

### Shell Aliases

Recommended aliases for .bashrc or .zshrc:
- mw: Short for {{slot:BRAND}}-worktree
- mwl: List worktrees
- mws: Switch to worktree
- mwg: Navigate with eval pattern
- mwsync: Sync current worktree
- mwclean: Clean merged worktrees

### Git Hooks Integration

Post-checkout hook actions:
- Detect worktree environment
- Update last access time in registry
- Check if sync needed with base branch
- Load worktree-specific environment

Pre-push hook actions:
- Detect if pushing from worktree
- Check for uncommitted changes
- Verify sync status with base
- Update registry with push timestamp

---

Version: 2.0.0
Last Updated: 2026-01-06
Module: Complete CLI command reference with usage examples
