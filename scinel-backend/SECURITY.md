# Security Policy

## Supported Versions

We are currently supporting the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously at ScIntel. If you discover a security vulnerability, please **DO NOT** open a public GitHub issue.

**Instead, please follow these steps:**

1. **Email us directly**: Send a detailed report to **security@scintel.local** (or arif@example.com for this MVP).
2. **Include the following details**:
   - Description of the vulnerability
   - Steps to reproduce (with code snippets/curl commands)
   - Potential impact (e.g., data leak, unauthorized access)
   - Suggested fix (if available)
3. **Expected Response Time**:
   - Initial response: < 48 hours
   - Status update: < 72 hours
   - Resolution timeline: Depends on severity (Critical: 7 days, High: 14 days)

## Disclosure Policy

When a vulnerability is reported and confirmed:
- We will patch it in the next release
- We will credit the reporter (if they wish) in the release notes
- We will notify users via release notes and documentation updates

## Security Best Practices for Developers

- **Never commit secrets** (.env.local, JWT_SECRET, DB passwords)
- **Always use prepared statements** (No raw SQL concatenation)
- **Validate all inputs** (Use strict type checking and validation functions)
- **Log all authentication attempts** (Especially failures)
- **Keep dependencies updated** (Run `go get -u` regularly)

**Thank you for helping keep ScIntel secure!** 🔒
