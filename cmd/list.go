package cmd

import (
	"fmt"
	"sort"

	"github.com/arnac-io/history-sanitizer/pkg/scanner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list-rules",
	Short: "List all available detection rules",
	Long:  `Display all detection rules provided by Gitleaks that are used for scanning.`,
	RunE:  runListRules,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListRules(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s Detection Rules (powered by Gitleaks)\n\n", green("🔍"))

	rules, err := scanner.GetDetectorInfo()
	if err != nil {
		return fmt.Errorf("failed to get detector rules: %w", err)
	}

	// Sort rules alphabetically
	sort.Strings(rules)

	fmt.Printf("Total rules: %s\n\n", cyan(len(rules)))

	// Group rules by category (based on prefix)
	categories := make(map[string][]string)
	for _, rule := range rules {
		// Simple categorization based on common prefixes
		category := "Other"
		if len(rule) > 0 {
			// Extract prefix (e.g., "aws-", "github-", etc.)
			for i, char := range rule {
				if char == '-' && i > 0 {
					category = rule[:i]
					break
				}
			}
		}
		categories[category] = append(categories[category], rule)
	}

	// Sort categories
	var categoryNames []string
	for cat := range categories {
		categoryNames = append(categoryNames, cat)
	}
	sort.Strings(categoryNames)

	// Display by category
	for _, cat := range categoryNames {
		if len(categories[cat]) > 0 {
			fmt.Printf("%s:\n", green(cat))
			for _, rule := range categories[cat] {
				fmt.Printf("  • %s\n", rule)
			}
			fmt.Println()
		}
	}

	return nil
}
