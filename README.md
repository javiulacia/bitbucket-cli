# bitbucket-cli

`bitbucket-cli` is a Go CLI for Bitbucket Cloud built on the official Bitbucket Cloud REST API v2.

This MVP implements:

- Authentication with Bitbucket Cloud using API token + email via HTTP Basic auth
- Authentication with access token via Bearer auth
- Repository commands for listing, viewing, and creating repositories
- Pull request commands for listing, viewing, and creating pull requests
- Human-readable output and `--json`

## Bitbucket API basis

The implementation follows the official Atlassian Bitbucket Cloud REST documentation:

- Auth intro: <https://developer.atlassian.com/cloud/bitbucket/rest/intro/#authentication>
- Current user: `GET /2.0/user`
- Repositories: `GET /2.0/repositories/{workspace}`, `GET/POST /2.0/repositories/{workspace}/{repo_slug}`
- Pull requests: `GET/POST /2.0/repositories/{workspace}/{repo_slug}/pullrequests`

## Build

```bash
go build ./cmd/bb
```

## Usage

Authenticate with an API token:

```bash
./bb auth login --email you@example.com --api-token <token> --workspace <workspace>
```

Authenticate with a bearer token:

```bash
./bb auth login --access-token <token> --workspace <workspace>
```

Check auth status:

```bash
./bb auth status
```

List repositories:

```bash
./bb repo list --workspace <workspace>
./bb repo list --json
```

View or create a repository:

```bash
./bb repo view <repo> --workspace <workspace>
./bb repo create <repo> --workspace <workspace>
```

List or view pull requests:

```bash
./bb pr list --repo <repo> --workspace <workspace>
./bb pr view 123 --repo <repo> --workspace <workspace>
```

Create a pull request:

```bash
./bb pr create \
  --repo <repo> \
  --workspace <workspace> \
  --title "My change" \
  --source feature/my-change \
  --destination main \
  --description "Details"
```

## Config

Credentials are stored in the user config directory:

- macOS: `~/Library/Application Support/bitbucket-cli/config.json`
- Linux: `~/.config/bitbucket-cli/config.json`
- Windows: `%AppData%\bitbucket-cli\config.json`

## Status

The repository started as a skeleton. The current implementation is a working MVP focused on the documented Bitbucket Cloud auth, repository, and pull request endpoints.
