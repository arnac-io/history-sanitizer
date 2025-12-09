package sanitizer

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/arnac-io/history-sanitizer/pkg/scanner"
)

// Sanitize obfuscates sensitive data in the content
func Sanitize(content string, findings []scanner.Finding) string {
	if len(findings) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")

	// Group findings by line for efficient processing
	findingsByLine := make(map[int][]scanner.Finding)
	for _, finding := range findings {
		findingsByLine[finding.Line] = append(findingsByLine[finding.Line], finding)
	}

	// Process each line that has findings
	for lineNum, lineFindings := range findingsByLine {
		if lineNum < 1 || lineNum > len(lines) {
			continue
		}

		lineIdx := lineNum - 1
		originalLine := lines[lineIdx]

		// Sort findings by start position (descending) to replace from end to start
		// This prevents position shifts during replacement
		sortedFindings := make([]scanner.Finding, len(lineFindings))
		copy(sortedFindings, lineFindings)

		// Simple bubble sort (descending by Start position)
		for i := 0; i < len(sortedFindings)-1; i++ {
			for j := 0; j < len(sortedFindings)-i-1; j++ {
				if sortedFindings[j].Start < sortedFindings[j+1].Start {
					sortedFindings[j], sortedFindings[j+1] = sortedFindings[j+1], sortedFindings[j]
				}
			}
		}

		// Replace each finding with obfuscated version
		newLine := originalLine
		for _, finding := range sortedFindings {
			if finding.Start < len(newLine) && finding.End <= len(newLine) {
				replacement := obfuscate(finding.Match, finding.Type)
				newLine = newLine[:finding.Start] + replacement + newLine[finding.End:]
			}
		}

		lines[lineIdx] = newLine
	}

	return strings.Join(lines, "\n")
}

// obfuscate generates a replacement for sensitive data
func obfuscate(original, secretType string) string {
	// Create a hash-based identifier that's consistent for the same value
	hash := sha256.Sum256([]byte(original))
	hashStr := fmt.Sprintf("%x", hash[:4])

	// Format based on type (matching Gitleaks pattern names)
	secretTypeLower := strings.ToLower(secretType)
	switch {
	case strings.Contains(secretTypeLower, "key") || strings.Contains(secretTypeLower, "aws"):
		return fmt.Sprintf("[REDACTED_KEY_%s]", hashStr)
	case strings.Contains(secretTypeLower, "token") || strings.Contains(secretTypeLower, "pat") || strings.Contains(secretTypeLower, "oauth"):
		return fmt.Sprintf("[REDACTED_TOKEN_%s]", hashStr)
	case strings.Contains(secretTypeLower, "password") || strings.Contains(secretTypeLower, "secret"):
		return fmt.Sprintf("[REDACTED_SECRET_%s]", hashStr)
	case strings.Contains(secretTypeLower, "url") || strings.Contains(secretTypeLower, "connection") || strings.Contains(secretTypeLower, "webhook"):
		return fmt.Sprintf("[REDACTED_URL_%s]", hashStr)
	case strings.Contains(secretTypeLower, "email"):
		return fmt.Sprintf("[REDACTED_EMAIL_%s]", hashStr)
	case strings.Contains(secretTypeLower, "credit") || strings.Contains(secretTypeLower, "card"):
		return "[REDACTED_CC]"
	default:
		return fmt.Sprintf("[REDACTED_%s]", hashStr)
	}
}

// ObfuscatePreview returns a preview of what the obfuscation would look like
func ObfuscatePreview(original, secretType string) string {
	if len(original) <= 8 {
		return obfuscate(original, secretType)
	}
	// Show first 2 and last 2 characters with asterisks in between
	return fmt.Sprintf("%s***%s", original[:2], original[len(original)-2:])
}
