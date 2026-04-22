# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| main    | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue.
2. Email security concerns to the project maintainer.
3. Include a description of the vulnerability, steps to reproduce, and potential impact.
4. Allow reasonable time for a fix before public disclosure.

## Security Practices

- All API endpoints require authentication (JWT or API key) except `/health` and auth routes.
- JWT secrets must be changed from defaults in production (`STATIC_CRON_TOKEN`, `AUTH_JWT_SECRET`).
- Rate limiting is enabled by default.
- RBAC with 21 group roles and 3 wallet roles enforces least-privilege access.
- Database credentials are never logged.
