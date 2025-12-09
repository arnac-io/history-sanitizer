package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arnac-io/history-sanitizer/pkg/sanitizer"
	"github.com/arnac-io/history-sanitizer/pkg/scanner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	historyFile string
	outputFile  string
	dryRun      bool
	verbose     bool
	inPlace     bool
)

var rootCmd = &cobra.Command{
	Use:   "history-sanitizer",
	Short: "Scan and sanitize sensitive information from shell history",
	Long: `history-sanitizer automatically scans your shell history files 
for sensitive information like passwords, API keys, tokens, and secrets,
then obfuscates them to keep your history safe.`,
	RunE: runSanitizer,
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultHistory := filepath.Join(homeDir, ".zsh_history")

	rootCmd.Flags().StringVarP(&historyFile, "file", "f", defaultHistory, "Path to history file")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: <input>.sanitized)")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be changed without modifying files")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")
	rootCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Replace the original file (creates backup with .backup suffix)")
}

func Execute() error {
	return rootCmd.Execute()
}

func runSanitizer(cmd *cobra.Command, args []string) error {
	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return fmt.Errorf("history file not found: %s", historyFile)
	}

	// Set default output file
	if outputFile == "" {
		outputFile = historyFile + ".sanitized"
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("🔍 Scanning history file: %s\n", yellow(historyFile))

	// Read history file
	content, err := os.ReadFile(historyFile)
	if err != nil {
		return fmt.Errorf("failed to read history file: %w", err)
	}

	// Scan for sensitive data
	findings, err := scanner.ScanContent(string(content))
	if err != nil {
		return fmt.Errorf("failed to scan content: %w", err)
	}

	if len(findings) == 0 {
		fmt.Println(green("✓ No sensitive information found!"))
		return nil
	}

	fmt.Printf("\n%s Found %d sensitive pattern(s)\n\n", red("⚠"), len(findings))

	// Display findings with obfuscation
	lines := strings.Split(string(content), "\n")
	for i, finding := range findings {
		fmt.Printf("Finding #%d:\n", i+1)
		fmt.Printf("  Type: %s\n", yellow(finding.Type))
		fmt.Printf("  Line: %d\n", finding.Line)

		// Show the full line with obfuscated secret
		if finding.Line > 0 && finding.Line <= len(lines) {
			fullLine := lines[finding.Line-1]
			obfuscated := sanitizer.ObfuscatePreview(finding.Match, finding.Type)
			sanitizedLine := strings.Replace(fullLine, finding.Match, red(obfuscated), 1)
			fmt.Printf("  Command: %s\n", sanitizedLine)
		}

		if verbose {
			fmt.Printf("  Secret: %s\n", red(finding.Match))
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Println(yellow("🔸 Dry run mode - no files will be modified"))
		return nil
	}

	// Sanitize content
	sanitized := sanitizer.Sanitize(string(content), findings)

	if inPlace {
		// Create backup with timestamp
		backupFile := historyFile + ".backup"
		err = os.WriteFile(backupFile, content, 0600)
		if err != nil {
			return fmt.Errorf("failed to create backup file: %w", err)
		}
		fmt.Printf("%s Backup created: %s\n", green("✓"), backupFile)

		// Replace original file
		err = os.WriteFile(historyFile, []byte(sanitized), 0600)
		if err != nil {
			return fmt.Errorf("failed to write sanitized file: %w", err)
		}
		fmt.Printf("%s History file sanitized: %s\n", green("✓"), green(historyFile))
	} else {
		// Write to output file
		err = os.WriteFile(outputFile, []byte(sanitized), 0600)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("%s Sanitized history saved to: %s\n", green("✓"), green(outputFile))
		fmt.Printf("\nOriginal file preserved at: %s\n", historyFile)
		fmt.Println("\nTo automatically replace your history file, use:")
		fmt.Printf("  %s -f %s -i\n", os.Args[0], historyFile)
	}

	return nil
}
