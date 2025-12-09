package sanitizer

import (
	"strings"
	"testing"

	"github.com/arnac-io/history-sanitizer/pkg/scanner"
)

func TestSanitize(t *testing.T) {
	content := `echo "hello"
export AWS_KEY=AKIAIOSFODNN7EXAMPLE
echo "world"`

	findings := []scanner.Finding{
		{
			Type:  "aws-access-token",
			Match: "AKIAIOSFODNN7EXAMPLE",
			Line:  2,
			Start: 15,
			End:   35, // Corrected: 15 + len("AKIAIOSFODNN7EXAMPLE") = 15 + 20 = 35
		},
	}

	result := Sanitize(content, findings)

	if strings.Contains(result, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("sanitized content still contains the original secret: %s", result)
	}

	if !strings.Contains(result, "[REDACTED_KEY_") {
		t.Errorf("sanitized content does not contain expected redaction marker: %s", result)
	}

	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines in output, got %d", len(lines))
	}

	if lines[0] != `echo "hello"` {
		t.Errorf("first line was modified unexpectedly: %s", lines[0])
	}
}

func TestSanitize_MultipleFindings(t *testing.T) {
	content := `export KEY1=secret123 KEY2=secret456`

	findings := []scanner.Finding{
		{
			Type:  "Generic Secret",
			Match: "secret123",
			Line:  1,
			Start: 12,
			End:   21,
		},
		{
			Type:  "Generic Secret",
			Match: "secret456",
			Line:  1,
			Start: 27,
			End:   36,
		},
	}

	result := Sanitize(content, findings)

	if strings.Contains(result, "secret123") || strings.Contains(result, "secret456") {
		t.Error("sanitized content still contains original secrets")
	}

	redactedCount := strings.Count(result, "[REDACTED_")
	if redactedCount < 2 {
		t.Errorf("expected at least 2 redaction markers, got %d", redactedCount)
	}
}

func TestSanitize_NoFindings(t *testing.T) {
	content := `echo "hello world"
ls -la
pwd`

	result := Sanitize(content, []scanner.Finding{})

	if result != content {
		t.Error("content should remain unchanged when there are no findings")
	}
}

func TestObfuscate(t *testing.T) {
	tests := []struct {
		name       string
		original   string
		secretType string
		wantPrefix string
	}{
		{
			name:       "API Key",
			original:   "AKIAIOSFODNN7EXAMPLE",
			secretType: "aws-access-token",
			wantPrefix: "[REDACTED_KEY_",
		},
		{
			name:       "Token",
			original:   "ghp_1234567890abcdefghijklmnop",
			secretType: "github-pat",
			wantPrefix: "[REDACTED_TOKEN_",
		},
		{
			name:       "Password",
			original:   "mysecretpassword",
			secretType: "generic-secret",
			wantPrefix: "[REDACTED_SECRET_",
		},
		{
			name:       "Credit Card",
			original:   "4532123456789012",
			secretType: "credit-card",
			wantPrefix: "[REDACTED_CC]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := obfuscate(tt.original, tt.secretType)
			if !strings.HasPrefix(result, tt.wantPrefix) {
				t.Errorf("obfuscate() = %v, want prefix %v", result, tt.wantPrefix)
			}
			if result == tt.original {
				t.Error("obfuscated value should not equal original")
			}
		})
	}
}

func TestObfuscate_Consistency(t *testing.T) {
	original := "AKIAIOSFODNN7EXAMPLE"
	secretType := "aws-access-token"

	result1 := obfuscate(original, secretType)
	result2 := obfuscate(original, secretType)

	if result1 != result2 {
		t.Error("obfuscate should produce consistent results for the same input")
	}
}

func TestObfuscatePreview(t *testing.T) {
	original := "AKIAIOSFODNN7EXAMPLE"
	secretType := "aws-access-token"

	preview := ObfuscatePreview(original, secretType)

	if len(preview) < 4 {
		t.Error("preview should contain at least first and last characters")
	}

	if !strings.Contains(preview, "***") && !strings.HasPrefix(preview, "[REDACTED_") {
		t.Error("preview should either show partial content or redaction marker")
	}
}
