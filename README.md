# History Sanitizer 🔒

A powerful command-line tool written in Go that automatically scans your shell history files for sensitive information and obfuscates it to keep your data safe.

## Why?

Shell history files are incredibly useful for daily work, but they can inadvertently store sensitive information such as:
- API keys and tokens
- Passwords and secrets
- Database connection strings
- Private keys
- Credit card numbers
- Authentication headers

**history-sanitizer** helps you maintain the utility of your history while protecting sensitive data.

## Features

- 🔍 **Smart Detection**: Uses detection patterns from [Gitleaks](https://github.com/gitleaks/gitleaks) - industry-leading, community-maintained secret scanner (15k+ stars)
- 🎨 **Colored Output**: Clear, colored terminal output for easy reading
- 🔐 **Safe Obfuscation**: Replaces sensitive data with redacted placeholders
- 💾 **Non-Destructive**: Creates a new sanitized file, preserving your original
- 🌈 **Multi-Shell Support**: Works with bash, zsh, fish, and other shell history formats
- 🚀 **Fast & Efficient**: Built with Go for speed and reliability
- 🧪 **Dry Run Mode**: Preview changes before applying them
- 🔄 **Auto-Updated Patterns**: Leverages Gitleaks' actively maintained detection rules

## Detection Patterns

The tool uses detection patterns **sourced from Gitleaks** - a well-maintained, community-driven project with patterns for:

**Cloud Providers & Services:**
- AWS (Access Keys, Secret Keys, Session Tokens, MWS Keys)
- Google Cloud (API Keys, Service Account credentials)
- Azure (Storage Keys, Client Secrets, Connection Strings)
- DigitalOcean, Heroku, Stripe, Twilio tokens

**Version Control & CI/CD:**
- GitHub (Personal Access Tokens, App Tokens, OAuth tokens)
- GitLab (Personal/Project/Runner tokens)
- Bitbucket, CircleCI, Travis CI tokens

**Credentials & Secrets:**
- Private Keys (RSA, EC, DSA, PGP, SSH)
- JWT Tokens
- Database connection strings (MongoDB, MySQL, PostgreSQL)
- Generic passwords, API keys, and secrets

**Communication & Monitoring:**
- Slack (Bot/App/User/Webhook tokens)
- Discord webhooks
- SendGrid, Mailgun, Telegram API keys
- Datadog, New Relic tokens

...and **200+ more patterns** in the full Gitleaks config (embedded for reference)!

### Why Not Use Gitleaks Directly?

We use Gitleaks' **patterns** rather than their Go API because:
- ✅ **Gitleaks patterns are open source and well-maintained** by a large community
- ✅ Gitleaks CLI is designed as a standalone tool, not a Go library
- ✅ Their API lacks simple programmatic access (`config.GetDefault()` doesn't exist)
- ✅ Direct pattern usage avoids 50+ transitive dependencies
- ✅ We get the same detection quality with simpler, more maintainable code

The full `gitleaks.toml` (95KB, 200+ rules) is embedded in our codebase and updated from the official Gitleaks repository.

## Installation

### From Source

```bash
git clone https://github.com/arnac-io/history-sanitizer.git
cd history-sanitizer
go build -o history-sanitizer
```

### Using Go Install

```bash
go install github.com/arnac-io/history-sanitizer@latest
```

## Usage

### Basic Usage

Scan and sanitize your default shell history (zsh):

```bash
./history-sanitizer
```

### Specify a History File

```bash
./history-sanitizer -f ~/.bash_history
```

### Dry Run (Preview Only)

See what would be changed without modifying any files:

```bash
./history-sanitizer --dry-run
```

### Verbose Output

Show detailed information about each finding:

```bash
./history-sanitizer -v
```

### List Available Detection Rules

See all detection rules provided by Gitleaks:

```bash
./history-sanitizer list-rules
```

### Custom Output File

```bash
./history-sanitizer -f ~/.bash_history -o ~/safe_history.txt
```

### Complete Example

```bash
# Scan with dry run to see what will be found
./history-sanitizer -f ~/.zsh_history --dry-run -v

# If satisfied, run the actual sanitization
./history-sanitizer -f ~/.zsh_history -o ~/.zsh_history.clean

# Review the cleaned file
less ~/.zsh_history.clean

# Replace original (make sure to backup first!)
cp ~/.zsh_history ~/.zsh_history.backup
mv ~/.zsh_history.clean ~/.zsh_history
```

## Command-Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--file` | `-f` | Path to history file | `~/.zsh_history` |
| `--output` | `-o` | Output file path | `<input>.sanitized` |
| `--dry-run` | `-d` | Show changes without modifying files | `false` |
| `--verbose` | `-v` | Show detailed information | `false` |
| `--help` | `-h` | Show help message | - |

### Additional Commands

| Command | Description |
|---------|-------------|
| `list-rules` | Display all available Gitleaks detection rules |

## Example Output

```
🔍 Scanning history file: /Users/you/.zsh_history

⚠ Found 3 sensitive pattern(s)

Finding #1:
  Type: AWS Access Key
  Line: 42

Finding #2:
  Type: Generic Secret
  Line: 108

Finding #3:
  Type: GitHub Token
  Line: 234

✓ Sanitized history saved to: /Users/you/.zsh_history.sanitized

Original file preserved at: /Users/you/.zsh_history

To replace your history file, run:
  mv /Users/you/.zsh_history.sanitized /Users/you/.zsh_history
```

## How It Works

1. **Scan**: Reads your shell history file and scans each line against known patterns
2. **Detect**: Uses regular expressions to identify sensitive information
3. **Obfuscate**: Replaces sensitive data with safe placeholders like `[REDACTED_KEY_a1b2c3d4]`
4. **Save**: Writes the sanitized content to a new file

## Security Considerations

- ✅ Original files are never modified automatically
- ✅ Obfuscated values include a hash for consistency
- ✅ Output files are created with restrictive permissions (0600)
- ✅ All processing happens locally - no data is sent anywhere

## Development

### Project Structure

```
history-sanitizer/
├── main.go                 # Entry point
├── cmd/
│   └── root.go            # CLI command definitions
├── pkg/
│   ├── scanner/
│   │   └── scanner.go     # Pattern detection logic
│   └── sanitizer/
│       └── sanitizer.go   # Obfuscation logic
├── go.mod                  # Go module definition
└── README.md              # This file
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o history-sanitizer
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o history-sanitizer-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o history-sanitizer-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o history-sanitizer.exe
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- Powered by [Gitleaks](https://github.com/gitleaks/gitleaks) for secret detection - a well-maintained, industry-standard tool
- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [fatih/color](https://github.com/fatih/color) for colored terminal output

## Roadmap

- [ ] Add configuration file support for custom patterns
- [ ] Support for more shell history formats
- [ ] Integration with git hooks
- [ ] Cloud backup sanitization
- [ ] Machine learning-based detection

## Documentation

- 📘 [Quick Start Guide](QUICKSTART.md) - Get started in 5 minutes
- 🔍 [Pattern Sources](PATTERN_SOURCES.md) - How we use Gitleaks patterns (3rd party)
- 📁 [Examples](examples/) - Sample history files for testing

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/arnac-io/history-sanitizer/issues) on GitHub.

---

**⚠️ Remember**: Always backup your history files before running any sanitization tool!

