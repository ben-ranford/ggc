# Upgrade Guide: ggc v6 CLI Unification

This release unifies the CLI into a consistent subcommand structure:

  ggc <command> <subcommand> [modifier] [arguments]

Legacy flags and hyphenated commands are replaced with space‑separated subcommands. When legacy syntax is used, ggc prints a friendly suggestion and exits without executing.

## Common Migrations

- Rebase:
  - ggc rebase -i → ggc rebase interactive
  - ggc rebase --interactive → ggc rebase interactive

- Add:
  - ggc add -p → ggc add patch
  - ggc add -i / --interactive → ggc add interactive

- Restore:
  - ggc restore --staged <file> → ggc restore staged <file>

- Fetch:
  - ggc fetch --prune → ggc fetch prune

- Commit:
  - ggc commit --allow-empty → ggc commit allow empty
  - ggc commit --amend → ggc commit amend
  - ggc commit amend --no-edit (or --amend --no-edit) → ggc commit amend no-edit

- Branch:
  - ggc branch checkout-remote → ggc branch checkout remote
  - ggc branch delete-merged → ggc branch delete merged
  - ggc branch set-upstream → ggc branch set upstream

- Clean:
  - ggc clean-interactive → ggc clean interactive

## Notes

- Help output, interactive mode, and completion scripts now reflect the unified syntax.
- No execution fallback for legacy syntax; use the suggested new command.

