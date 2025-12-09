# Pattern Sources & Attribution

## Overview

This project uses secret detection patterns **sourced from the Gitleaks project**, a well-maintained, community-driven secret scanner.

## Why Use 3rd Party Patterns?

Maintaining high-quality secret detection patterns requires:
- Deep knowledge of various service authentication schemes
- Regular updates as providers change their formats
- Community testing and validation
- Low false-positive rates

Rather than maintaining our own patterns, we leverage **Gitleaks** - the industry standard with:
- ⭐ **15,000+ GitHub stars**
- 👥 **Active community** contributing patterns
- 🔄 **Regular updates** for new services
- 🎯 **200+ detection rules** covering major providers
- ✅ **Battle-tested** by thousands of organizations

## Pattern Attribution

All detection patterns are sourced from:

**Project:** Gitleaks  
**Repository:** https://github.com/gitleaks/gitleaks  
**License:** MIT License  
**Config File:** https://github.com/gitleaks/gitleaks/blob/master/config/gitleaks.toml  

### Active Patterns

Our active detection patterns (36) are stored in:
```
pkg/scanner/patterns.toml
```

These are extracted from Gitleaks and optimized for shell history scanning.

### Reference Config

The full `gitleaks.toml` configuration file (95KB, 200+ rules) is embedded in:
```
pkg/scanner/gitleaks.toml
```

This provides:
1. **Reference documentation** for all available Gitleaks patterns
2. **Proof of source** - patterns are from Gitleaks, not custom
3. **Update capability** - can sync with upstream Gitleaks releases
4. **Pattern discovery** - find new patterns to add to patterns.toml

## Why Not Use Gitleaks as a Library?

### The Challenge

Gitleaks v8 is designed as a **CLI tool**, not a Go library:
```go
// This doesn't exist in Gitleaks v8 API:
cfg, err := config.GetDefault()  // ❌ undefined
```

### Our Approach

Instead of using the Gitleaks API, we:
1. ✅ **Extract patterns** from Gitleaks' open-source config
2. ✅ **Store in TOML** at `pkg/scanner/patterns.toml`
3. ✅ **Load and compile at runtime** using Go's regexp package
4. ✅ **Embed the full config** for reference at `pkg/scanner/gitleaks.toml`
5. ✅ **Credit the source** in code and documentation

This gives us:
- Same detection quality as Gitleaks
- Simpler, more maintainable code
- No complex dependency tree (50+ transitive deps)
- Easy to add/modify patterns without code changes
- Still using 3rd party (Gitleaks) patterns, just more directly
- Clear separation between patterns (TOML) and code (Go)

## Pattern List

Currently implemented patterns in `pkg/scanner/patterns.toml` (extracted from Gitleaks):

### AWS (5 patterns)
- `aws-access-token` - AWS Access Keys (AKIA, ASIA, AGPA, etc.)
- `aws-secret-key` - AWS Secret Access Keys
- `aws-session-token` - AWS Session Tokens
- `aws-sso-token` - AWS SSO Access Tokens
- `cli-access-token` - CLI --access-token flags (generic)

### GitHub (5 patterns)
- `github-pat` - GitHub Personal Access Tokens (ghp_)
- `github-oauth` - GitHub OAuth Tokens (gho_)
- `github-app-token` - GitHub App Tokens (ghu_, ghs_)
- `github-refresh-token` - GitHub Refresh Tokens (ghr_)
- `github-fine-grained-pat` - GitHub Fine-Grained PATs

### Slack (2 patterns)
- `slack-token` - Slack API Tokens (xoxb-, xoxp-, etc.)
- `slack-webhook` - Slack Webhook URLs

### Generic (2 patterns)
- `generic-api-key` - Generic API key patterns
- `generic-secret` - Generic passwords/secrets

### Cryptographic (2 patterns)
- `private-key` - Private keys (RSA, EC, DSA, PGP, SSH)
- `jwt` - JSON Web Tokens

### Database (2 patterns)
- `connection-string` - MongoDB, MySQL, PostgreSQL connection strings
- `mongodb-password` - MongoDB passwords in connection strings

### Cloud & Services (7 patterns)
- `google-api-key` - Google API Keys
- `stripe-key` - Stripe API Keys
- `heroku-api-key` - Heroku API Keys / UUID patterns
- `twilio-api-key` - Twilio API Keys
- `mailchimp-api-key` - MailChimp API Keys
- `sendgrid-api-key` - SendGrid API Keys
- `square-access-token` - Square Access Tokens

### HTTP/Curl Authentication (4 patterns)
- `proxy-user-credentials` - Curl proxy credentials (--proxy-user)
- `curl-user-credentials` - Curl HTTP basic auth (-u, --user)
- `authorization-bearer` - Authorization Bearer tokens
- `authorization-basic` - Authorization Basic tokens

### Other (7 patterns)
- `cli-secret-value` - CLI secret value flags
- `datadog-api-key` - Datadog/Generic 32-char hex keys
- `pagerduty-token` - PagerDuty tokens
- `proxy-password-url` - Proxy URLs with passwords
- `1password-service-token` - 1Password service tokens
- `env-var-jwt` - Environment variables with JWTs
- `env-var-secret` - Environment variables with secrets

**Total:** 36 patterns (extracted from Gitleaks' 200+ rule set, plus custom patterns for shell history)

## Updating Patterns

To sync with latest Gitleaks patterns:

```bash
# 1. Download latest gitleaks.toml (reference)
curl -sL https://raw.githubusercontent.com/gitleaks/gitleaks/master/config/gitleaks.toml \
  -o pkg/scanner/gitleaks.toml

# 2. Review new rules in the file
grep -A 5 "^\[\[rules\]\]" pkg/scanner/gitleaks.toml | less

# 3. Extract relevant patterns to patterns.toml
# Edit pkg/scanner/patterns.toml and add new patterns in TOML format:
# [[patterns]]
# name = "pattern-name"
# regex = '''pattern-regex'''
# description = "Pattern description"

# 4. Test the new patterns
go test ./pkg/scanner/...

# 5. Verify patterns are loaded
go build
./history-sanitizer list-rules
```

## License Compliance

### Gitleaks License
MIT License - Copyright (c) 2019 Zachary Rice

Permission is hereby granted, free of charge, to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software.

### Our Usage
We comply with the MIT license by:
- ✅ Crediting Gitleaks as the pattern source
- ✅ Including license attribution in documentation
- ✅ Maintaining original copyright notices
- ✅ Not claiming patterns as our own work

## Contributing Patterns

If you discover a missing secret type:

1. **Add to patterns.toml** (immediate)
   - Edit `pkg/scanner/patterns.toml`
   - Add pattern in TOML format
   - Test with `go test ./pkg/scanner/...`
   - Submit PR to this project

2. **Contribute upstream to Gitleaks** (preferred for general patterns)
   - Benefits the entire security community
   - Gets expert review from Gitleaks maintainers
   - Will be included in future syncs
   - Helps improve the industry standard

3. **Workflow for new patterns**
   - Add to our `patterns.toml` first (quick)
   - Submit to Gitleaks (community benefit)
   - Sync back when merged upstream

## References

- Gitleaks Project: https://github.com/gitleaks/gitleaks
- Gitleaks Config: https://github.com/gitleaks/gitleaks/blob/master/config/gitleaks.toml
- Contributing to Gitleaks: https://github.com/gitleaks/gitleaks/blob/master/CONTRIBUTING.md
- Regex Secret Scanning: https://lookingatcomputer.substack.com/p/regex-is-almost-all-you-need

---

**Summary:** We use 3rd party patterns from Gitleaks, a well-maintained industry standard. We implement them directly for simplicity while crediting the source and embedding the full config for reference.

