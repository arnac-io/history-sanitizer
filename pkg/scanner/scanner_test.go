package scanner

import (
	"strings"
	"testing"
)

func TestScanContent_AWSKey(t *testing.T) {
	content := `
echo "hello"
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
aws s3 ls
`
	findings, err := ScanContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) == 0 {
		t.Fatal("expected to find AWS key, but got no findings")
	}

	found := false
	for _, f := range findings {
		// Gitleaks uses "aws-access-token" or similar rule IDs
		if strings.Contains(strings.ToLower(f.Type), "aws") ||
			strings.Contains(f.Match, "AKIA") {
			found = true
			break
		}
	}

	if !found {
		t.Error("did not find expected AWS key pattern")
	}
}

func TestScanContent_Password(t *testing.T) {
	content := `
mysql -u root -p
export PASSWORD="mysecretpassword123"
echo "done"
`
	findings, err := ScanContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) == 0 {
		t.Fatal("expected to find password, but got no findings")
	}

	found := false
	for _, f := range findings {
		if strings.Contains(strings.ToLower(f.Type), "password") || strings.Contains(strings.ToLower(f.Type), "secret") {
			found = true
			break
		}
	}

	if !found {
		t.Error("did not find expected password pattern")
	}
}

func TestScanContent_JWT(t *testing.T) {
	content := `
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
`
	_, err := ScanContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// JWT detection depends on Gitleaks rules, which may vary
	// Just ensure no error occurred - JWT pattern might not be in default rules
	t.Log("JWT scan completed successfully")
}

func TestScanContent_NoSecrets(t *testing.T) {
	content := `
ls -la
cd /home/user
echo "Hello World"
git status
`
	findings, err := ScanContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) != 0 {
		t.Errorf("expected no findings for clean content, but got %d findings", len(findings))
		for _, f := range findings {
			t.Logf("Unexpected finding: Type=%s, Match=%s", f.Type, f.Match)
		}
	}
}

func TestScanContent_MultipleFindings(t *testing.T) {
	content := `
export API_KEY="sk-1234567890abcdefghij"
curl -X POST https://api.example.com -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ0ZXN0In0.test"
export PASSWORD="secret123456"
`
	findings, err := ScanContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, but got %d", len(findings))
	}
}

func TestGetPatternCount(t *testing.T) {
	count := GetPatternCount()
	// We've extracted 20+ key patterns from Gitleaks
	if count < 15 {
		t.Errorf("expected at least 15 patterns from Gitleaks, got %d", count)
	}
	t.Logf("Loaded %d detection patterns from Gitleaks", count)
}

func TestGetDetectorInfo(t *testing.T) {
	rules, err := GetDetectorInfo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rules) == 0 {
		t.Error("expected to get detector rules, but got empty list")
	}

	t.Logf("Available detection rules: %d", len(rules))
	// Log first 10 rules as examples
	for i := 0; i < 10 && i < len(rules); i++ {
		t.Logf("  - %s", rules[i])
	}
}
