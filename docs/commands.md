# Commands

## Overview

bitbucket-cli provides a set of commands to manage Bitbucket workflows directly from the terminal.

Commands are grouped by domain:

- auth
- repo
- pr
- utils

---

## Auth

### Login

bb auth login

Authenticate with Bitbucket using app password or token.

---

### Status

bb auth status

Show current authentication status.

---

### Logout

bb auth logout

Remove stored credentials.

---

## Repositories

### List repositories

bb repo list

List repositories in the current workspace.

Options:

bb repo list --workspace myteam
bb repo list --json

---

### Clone repository

bb repo clone <repo>

Example:

bb repo clone my-repo

---

### Clone all repositories

bb repo clone-all

Options:

bb repo clone-all --workspace myteam
bb repo clone-all --dir ~/Documents/Deepcom
bb repo clone-all --parallel 5

---

### View repository

bb repo view <repo>

---

### Create repository

bb repo create <repo>

---

## Pull Requests

### List PRs

bb pr list

Options:

bb pr list --state open
bb pr list --json

---

### View PR

bb pr view <id>

---

### Create PR

bb pr create

Auto-detects:

- current repository
- current branch
- default target branch

---

### Checkout PR

bb pr checkout <id>

---

### Merge PR

bb pr merge <id>

---

### Comment on PR

bb pr comment <id> --body "message"

---

## Utilities

### Open in browser

bb browse

Opens the current repository in the browser.

---

### Sync repositories

bb sync

Updates local repositories:

- fetch
- prune
- pull default branch

---

### Output formats

Most commands support:

--json

Example:

bb repo list --json

---

## Examples

### Clone all repos into a directory

bb repo clone-all --workspace deepcom --dir ~/Documents/Deepcom

---

### Create a PR from current branch

bb pr create

---

### List open PRs in JSON

bb pr list --state open --json

---

## Notes

- Commands provide sensible defaults
- Context is inferred when possible
- Designed for scripting and automation

---

## Future Commands

Planned additions:

- pipeline commands
- workspace management
- Jira integration
- interactive mode