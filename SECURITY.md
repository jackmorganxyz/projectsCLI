# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest release | Yes |
| Older releases | No |

## Reporting a Vulnerability

If you discover a security vulnerability in projects, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please email the maintainer directly or use [GitHub's private vulnerability reporting](https://github.com/jackmorganxyz/projectsCLI/security/advisories/new).

## Scope

projects interacts with your system in the following ways:

- **Filesystem**: Reads and writes files under `~/.projects/` (configurable)
- **Git**: Executes `git` commands on project directories
- **GitHub**: Uses the `gh` CLI to create repositories (only when you run `push`)
- **No telemetry**: projects does not phone home or collect any data

## Dependencies

projects is built with Go and uses well-maintained open source libraries. Dependencies are tracked in `projectsCLI/go.mod` and can be audited with `go mod verify`.
