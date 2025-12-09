# Quick Start Guide

## Installation

### Method 1: Homebrew (Recommended)

```bash
brew tap arnac-io/tap
brew install history-sanitizer
```

### Method 2: From Source

1. Clone the repository:
```bash
git clone https://github.com/arnac-io/history-sanitizer.git
cd history-sanitizer
```

2. Download dependencies:
```bash
go mod download
```

3. Build the tool:
```bash
go build -o history-sanitizer
# or simply: make build
```

## First Run

### 1. Test with Sample Data

We've included a sample history file for testing:

```bash
# Run a dry-run on the sample file
./history-sanitizer -f examples/sample_history.txt --dry-run -v
```

This will show you what secrets were detected without modifying any files.

### 2. Check Your Own History

Before sanitizing your real history, **always backup first**:

```bash
# Backup your history
cp ~/.zsh_history ~/.zsh_history.backup

# Run dry-run to see what will be detected
./history-sanitizer -f ~/.zsh_history --dry-run -v
```

### 3. Sanitize Your History

Once you're comfortable with what will be changed:

**Option 1: Create a sanitized copy (safer)**
```bash
# Create sanitized version
./history-sanitizer -f ~/.zsh_history

# Review the sanitized file
less ~/.zsh_history.sanitized

# If satisfied, replace original (you already have a backup!)
mv ~/.zsh_history.sanitized ~/.zsh_history
```

**Option 2: In-place replacement (automatic backup)**
```bash
# Sanitize and replace original (creates .backup automatically)
./history-sanitizer -f ~/.zsh_history -i

# This creates ~/.zsh_history.backup and updates ~/.zsh_history
```

## Understanding the Output

### Example Output

```
🔍 Scanning history file: /Users/you/.zsh_history

⚠ Found 3 sensitive pattern(s)

Finding #1:
  Type: aws-access-token
  Line: 42

Finding #2:
  Type: generic-api-key
  Line: 108

Finding #3:
  Type: github-pat
  Line: 234

✓ Sanitized history saved to: /Users/you/.zsh_history.sanitized
```

### What Gets Redacted

Before:
```bash
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
```

After:
```bash
export AWS_ACCESS_KEY_ID=[REDACTED_KEY_a1b2c3d4]
```

The hash suffix (`a1b2c3d4`) is consistent for the same value, helping you identify if the same secret appears multiple times.

## Exploring Detection Rules

See all 30+ detection rules currently implemented:

```bash
./history-sanitizer list-rules
```

This displays all patterns extracted from Gitleaks that the tool actively uses for detection.

## Common Use Cases

### Regular Cleanup

Add to your dotfiles or cron:

```bash
# Weekly history sanitization (with automatic backup)
0 0 * * 0 /path/to/history-sanitizer -f ~/.zsh_history -i
```

### CI/CD Integration

Before committing dotfiles:

```bash
# In your pre-commit hook
history-sanitizer -f .zsh_history --dry-run || exit 1
```

### Multiple Shell Support

Sanitize bash, zsh, or fish history:

```bash
# Bash
./history-sanitizer -f ~/.bash_history

# Zsh
./history-sanitizer -f ~/.zsh_history

# Fish
./history-sanitizer -f ~/.local/share/fish/fish_history
```

## Tips

1. **Always backup first** - The tool is designed to be safe, but backups are essential
2. **Use dry-run mode** - Test before applying changes
3. **Review sanitized output** - Check the `.sanitized` file before replacing original
4. **Keep tool updated** - Gitleaks patterns are actively maintained and improved
5. **Check periodically** - Run scans regularly to catch new sensitive data

## Troubleshooting

### "No such file or directory"

Make sure the history file path is correct. Use `echo $HISTFILE` to find your shell's history location.

### "Permission denied"

History files often have restricted permissions. Use `sudo` if needed, or ensure you own the file.

### Too many false positives?

Some patterns might be overly sensitive. Consider:
- Reviewing the `list-rules` output
- Opening an issue if a specific rule causes problems
- The tool uses Gitleaks defaults, which prioritize security over convenience

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check the [examples/](examples/) directory for sample files
- Contribute detection patterns upstream to [Gitleaks](https://github.com/gitleaks/gitleaks)

## Getting Help

- 📖 [Full Documentation](README.md)
- 🐛 [Report Issues](https://github.com/arnac-io/history-sanitizer/issues)
- 💬 [Discussions](https://github.com/arnac-io/history-sanitizer/discussions)

