# Architecture

## Overview

bitbucket-cli is a modular, developer-focused command-line tool written in Go.
The architecture is designed to be simple, testable, and scalable as new features are added.

The project follows a layered structure that separates:

- CLI commands
- Business logic
- External integrations (Bitbucket API, Git)
- Output rendering

---

## High-Level Structure

cmd/            → CLI entrypoints and Cobra commands
internal/       → Core application logic (private)
  ├── cli/      → Command implementations
  ├── bitbucket/→ Bitbucket API client
  ├── git/      → Local Git operations
  ├── config/   → Configuration management
  ├── output/   → Rendering (table, JSON)
  ├── browser/  → Open links in browser
  ├── prompt/   → Interactive prompts
  └── parallel/ → Concurrency utilities
pkg/            → Public reusable packages (optional)
testdata/       → Fixtures and golden files

---

## Design Principles

### 1. Separation of Concerns

Each layer has a single responsibility:

- CLI → parse input and flags
- Services → execute logic
- API client → interact with Bitbucket
- Output → render results

---

### 2. Thin CLI Layer

Commands should be lightweight and delegate logic.

The CLI should only:
- parse arguments
- call services
- handle errors
- print output

---

### 3. API Isolation

The Bitbucket API is encapsulated in:

internal/bitbucket/

Responsibilities:

- HTTP requests
- Authentication
- Pagination
- Error handling
- Data models

---

### 4. Git Operations Separation

All local Git interactions are isolated in:

internal/git/

Examples:

- clone repository
- detect current repo
- get current branch
- manage remotes

---

### 5. Config Management

User configuration is stored locally in:

internal/config/

Typical config example:

host: bitbucket.org
workspace: myteam
username: javier
git_protocol: ssh

---

### 6. Output Layer

All formatting is centralized in:

internal/output/

Supports:

- human-readable tables
- JSON output (--json)
- consistent formatting across commands

---

### 7. Concurrency

Parallel operations (e.g., cloning multiple repos) use:

internal/parallel/

This enables:

- controlled concurrency
- better performance
- predictable resource usage

---

## Execution Flow

Example: bb repo list

CLI → service → Bitbucket API → output

---

## Error Handling

- Errors should include context
- CLI prints user-friendly messages
- Internal errors remain structured

---

## Extensibility

The architecture allows easy addition of:

- new commands
- new API endpoints
- new output formats
- integrations (e.g., Jira)

---

## Future Improvements

- Plugin system
- Interactive mode (TUI)
- Offline caching
- Bitbucket Server support
- Advanced authentication (OAuth)

---

## Summary

The architecture is designed to:

- stay simple at MVP stage
- scale with new features
- remain maintainable and testable