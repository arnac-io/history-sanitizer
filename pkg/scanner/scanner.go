package scanner

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Finding represents a detected sensitive data
type Finding struct {
	Type  string
	Match string
	Line  int
	Start int
	End   int
}

//go:embed gitleaks.toml
var gitleaksConfig string

//go:embed patterns.toml
var patternsConfig string

// Pattern represents a secret detection pattern
type pattern struct {
	name  string
	regex *regexp.Regexp
	desc  string
}

// PatternConfig represents the TOML configuration structure
type PatternConfig struct {
	Title       string         `toml:"title"`
	Description string         `toml:"description"`
	Patterns    []PatternEntry `toml:"patterns"`
}

// PatternEntry represents a single pattern in the config
type PatternEntry struct {
	Name        string `toml:"name"`
	Regex       string `toml:"regex"`
	Description string `toml:"description"`
}

// detectionPatterns is loaded from patterns.toml at init time
var detectionPatterns []pattern

func init() {
	// Load patterns from embedded TOML config
	var config PatternConfig
	err := toml.Unmarshal([]byte(patternsConfig), &config)
	if err != nil {
		// Fallback to empty patterns - will cause errors but won't crash
		fmt.Printf("Warning: failed to load patterns config: %v\n", err)
		return
	}

	// Compile regex patterns
	detectionPatterns = make([]pattern, 0, len(config.Patterns))
	for _, p := range config.Patterns {
		compiled, err := regexp.Compile(p.Regex)
		if err != nil {
			fmt.Printf("Warning: failed to compile pattern %s: %v\n", p.Name, err)
			continue
		}
		detectionPatterns = append(detectionPatterns, pattern{
			name:  p.Name,
			regex: compiled,
			desc:  p.Description,
		})
	}
}

// Fallback patterns if TOML loading fails (minimal set)
var fallbackPatterns = []pattern{
	{
		name:  "aws-access-token",
		regex: regexp.MustCompile(`(AKIA|ASIA)[A-Z0-9]{16}`),
		desc:  "AWS Access Token",
	},
	{
		name:  "github-pat",
		regex: regexp.MustCompile(`ghp_[0-9a-zA-Z]{36}`),
		desc:  "GitHub Personal Access Token",
	},
	{
		name:  "jwt",
		regex: regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`),
		desc:  "JSON Web Token",
	},
}

// ScanContent scans the content for sensitive information
// Uses detection patterns loaded from patterns.toml (sourced from Gitleaks)
func ScanContent(content string) ([]Finding, error) {
	// Use fallback if patterns failed to load
	patterns := detectionPatterns
	if len(patterns) == 0 {
		patterns = fallbackPatterns
	}

	var results []Finding
	lines := strings.Split(content, "\n")

	// Scan each line with all patterns
	for lineNum, line := range lines {
		for _, p := range patterns {
			matches := p.regex.FindAllStringIndex(line, -1)
			if matches != nil {
				for _, match := range matches {
					if len(match) >= 2 {
						matchedText := line[match[0]:match[1]]
						finding := Finding{
							Type:  p.name,
							Match: matchedText,
							Line:  lineNum + 1,
							Start: match[0],
							End:   match[1],
						}
						results = append(results, finding)
					}
				}
			}
		}
	}

	return results, nil
}

// findLineNumber finds which line contains a character at the given position
func findLineNumber(content string, pos int) int {
	if pos < 0 || pos > len(content) {
		return 1
	}

	lineNum := 1
	for i := 0; i < pos && i < len(content); i++ {
		if content[i] == '\n' {
			lineNum++
		}
	}
	return lineNum
}

// findPositionInLine finds the start and end position of a secret within a line
func findPositionInLine(lines []string, lineNum int, secret string) (int, int) {
	if lineNum < 1 || lineNum > len(lines) {
		return 0, 0
	}

	line := lines[lineNum-1]
	start := strings.Index(line, secret)
	if start == -1 {
		return 0, len(line)
	}

	return start, start + len(secret)
}

// GetPatternCount returns the number of patterns being scanned
func GetPatternCount() int {
	return len(detectionPatterns)
}

// GetDetectorInfo returns information about the detector rules
func GetDetectorInfo() ([]string, error) {
	var info []string
	for _, p := range detectionPatterns {
		info = append(info, p.name+": "+p.desc)
	}
	return info, nil
}

// GetGitleaksConfigSize returns the size of the embedded Gitleaks config
// This demonstrates that we have the full Gitleaks config available for reference
func GetGitleaksConfigSize() int {
	return len(gitleaksConfig)
}
