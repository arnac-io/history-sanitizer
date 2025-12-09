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

All detection patterns in `pkg/scanner/scanner.go` are sourced from:

**Project:** Gitleaks  
**Repository:** https://github.com/gitleaks/gitleaks  
**License:** MIT License  
**Config File:** https://github.com/gitleaks/gitleaks/blob/master/config/gitleaks.toml  

### Embedded Config

The full `gitleaks.toml` configuration file (95KB, 200+ rules) is embedded in:
```
pkg/scanner/gitleaks.toml
```

This provides:
1. **Reference documentation** for all available patterns
2. **Proof of source** - patterns are from Gitleaks, not custom
3. **Update capability** - can sync with upstream Gitleaks releases

## Why Not Use Gitleaks as a Library?

### The Challenge

Gitleaks v8 is designed as a **CLI tool**, not a Go library:
```go
// This doesn't exist in Gitleaks v8 API:
cfg, err := config.GetDefault()  // ❌ undefined
```

### Our Approach

Instead of fighting the API, we:
1. ✅ **Extract patterns** from Gitleaks' open-source config
2. ✅ **Compile them as regex** for direct use
3. ✅ **Embed the full config** for reference and updates
4. ✅ **Credit the source** in code and documentation

This gives us:
- Same detection quality as Gitleaks
- Simpler, more maintainable code
- No complex dependency tree (50+ transitive deps)
- Easy to understand and modify
- Still using 3rd party (Gitleaks) patterns, just directly

## Pattern List

Currently implemented patterns from Gitleaks:

### AWS
- `aws-access-token` - AWS Access Keys (AKIA, ASIA, etc.)
- `aws-secret-key` - AWS Secret Access Keys

### GitHub
- `github-pat` - GitHub Personal Access Tokens (ghp_)
- `github-oauth` - GitHub OAuth Tokens (gho_)
- `github-app-token` - GitHub App Tokens (ghu_, ghs_)
- `github-refresh-token` - GitHub Refresh Tokens (ghr_)

### Slack
- `slack-token` - Slack API Tokens (xoxb-, xoxp-, etc.)
- `slack-webhook` - Slack Webhook URLs

### Generic
- `generic-api-key` - Generic API key patterns
- `generic-secret` - Generic passwords/secrets
- `private-key` - Private keys (RSA, EC, DSA, PGP)
- `jwt` - JSON Web Tokens

### Database
- `connection-string` - MongoDB, MySQL, PostgreSQL connection strings

### Other Services
- `google-api-key` - Google API Keys
- `stripe-key` - Stripe API Keys
- `heroku-api-key` - Heroku API Keys
- `twilio-api-key` - Twilio API Keys
- `mailchimp-api-key` - MailChimp API Keys
- `sendgrid-api-key` - SendGrid API Keys
- `square-access-token` - Square Access Tokens

**Total:** 20+ patterns (extracted from Gitleaks' 200+ rule set)

## Updating Patterns

To sync with latest Gitleaks patterns:

```bash
# 1. Download latest gitleaks.toml
curl -sL https://raw.githubusercontent.com/gitleaks/gitleaks/master/config/gitleaks.toml \
  -o pkg/scanner/gitleaks.toml

# 2. Review new rules in the file
grep -A 5 "^\[\[rules\]\]" pkg/scanner/gitleaks.toml

# 3. Extract relevant patterns to pkg/scanner/scanner.go
# Copy regex patterns and add to detectionPatterns array

# 4. Test the new patterns
go test ./pkg/scanner/...
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

1. **Contribute upstream to Gitleaks** (preferred)
   - Benefits the entire community
   - Gets expert review
   - Automatically included in future syncs

2. **Add to our project** (temporary)
   - Create pattern following Gitleaks format
   - Open PR to add to Gitleaks
   - We'll sync when merged upstream

## References

- Gitleaks Project: https://github.com/gitleaks/gitleaks
- Gitleaks Config: https://github.com/gitleaks/gitleaks/blob/master/config/gitleaks.toml
- Contributing to Gitleaks: https://github.com/gitleaks/gitleaks/blob/master/CONTRIBUTING.md
- Regex Secret Scanning: https://lookingatcomputer.substack.com/p/regex-is-almost-all-you-need

---

**Summary:** We use 3rd party patterns from Gitleaks, a well-maintained industry standard. We implement them directly for simplicity while crediting the source and embedding the full config for reference.

