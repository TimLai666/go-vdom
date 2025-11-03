package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Issue represents a detected problem in the code
type Issue struct {
	File    string
	Line    int
	Column  int
	Type    string
	Message string
	Context string
}

var (
	// Patterns to detect dangerous template usage
	quotedVarInExpr   = regexp.MustCompile(`\$\{[^}]*['"]{{[^}]+}}['"][^}]*\}`)
	stringBoolCompare = regexp.MustCompile(`['"]{{[^}]+}}['"][\s]*===[\s]*['"](?:true|false)['"]`)
	stringPlusVar     = regexp.MustCompile(`\$\{[^}]*['"][^'"]*['"][\s]*\+[\s]*{{[^}]+}}`)
	varPlusString     = regexp.MustCompile(`\$\{[^}]*{{[^}]+}}[\s]*\+[\s]*['"]`)
	doubleQuotedVar   = regexp.MustCompile(`"{{[^}]+}}"`)
	singleQuotedVar   = regexp.MustCompile(`'{{[^}]+}}'`)

	verbose = flag.Bool("v", false, "Verbose output")
	fix     = flag.Bool("fix", false, "Suggest fixes for issues")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."} // Default to current directory
	}

	issues := []Issue{}

	for _, path := range args {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip non-Go files and test files
			if info.IsDir() || !strings.HasSuffix(filePath, ".go") {
				return nil
			}

			// Skip the linter itself and vendor directories
			if strings.Contains(filePath, "template-linter") ||
				strings.Contains(filePath, "vendor") ||
				strings.Contains(filePath, ".git") {
				return nil
			}

			fileIssues := checkFile(filePath)
			issues = append(issues, fileIssues...)

			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking path %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	// Report issues
	if len(issues) == 0 {
		fmt.Println("✓ No template issues found!")
		os.Exit(0)
	}

	fmt.Printf("Found %d template issue(s):\n\n", len(issues))

	for _, issue := range issues {
		fmt.Printf("%s:%d:%d: %s\n", issue.File, issue.Line, issue.Column, issue.Type)
		fmt.Printf("  %s\n", issue.Message)
		if issue.Context != "" {
			fmt.Printf("  Context: %s\n", truncate(issue.Context, 100))
		}
		if *fix {
			suggestion := suggestFix(issue)
			if suggestion != "" {
				fmt.Printf("  💡 Suggestion: %s\n", suggestion)
			}
		}
		fmt.Println()
	}

	os.Exit(1)
}

func checkFile(filePath string) []Issue {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
		return nil
	}
	defer file.Close()

	issues := []Issue{}
	scanner := bufio.NewScanner(file)
	lineNum := 0
	inTemplate := false
	templateLines := []string{}
	templateStartLine := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Detect template start (common patterns)
		if strings.Contains(line, "Template:") ||
			strings.Contains(line, "template:") ||
			(strings.Contains(line, "`") && strings.Contains(line, "{{")) {
			inTemplate = true
			templateStartLine = lineNum
			templateLines = []string{line}
			continue
		}

		if inTemplate {
			templateLines = append(templateLines, line)

			// Check if template ends
			if strings.Contains(line, "`,") || strings.Contains(line, "`}") || strings.Contains(line, "` +") {
				// Check the accumulated template
				templateContent := strings.Join(templateLines, "\n")
				issues = append(issues, checkTemplate(filePath, templateStartLine, templateContent)...)
				inTemplate = false
				templateLines = []string{}
			}
		}

		// Also check individual lines for inline templates
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			issues = append(issues, checkLine(filePath, lineNum, line)...)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning file %s: %v\n", filePath, err)
	}

	return issues
}

func checkTemplate(filePath string, startLine int, content string) []Issue {
	issues := []Issue{}
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		lineNum := startLine + i
		issues = append(issues, checkLine(filePath, lineNum, line)...)
	}

	return issues
}

func checkLine(filePath string, lineNum int, line string) []Issue {
	issues := []Issue{}

	// Check 1: Quoted template variables inside ${}
	if match := quotedVarInExpr.FindString(line); match != "" {
		col := strings.Index(line, match) + 1
		issues = append(issues, Issue{
			File:    filePath,
			Line:    lineNum,
			Column:  col,
			Type:    "quoted-var-in-expression",
			Message: "Template variable is quoted inside ${} expression. This will be treated as a string literal.",
			Context: match,
		})
	}

	// Check 2: String comparison with boolean literals
	if match := stringBoolCompare.FindString(line); match != "" {
		col := strings.Index(line, match) + 1
		issues = append(issues, Issue{
			File:    filePath,
			Line:    lineNum,
			Column:  col,
			Type:    "string-bool-comparison",
			Message: "Comparing template variable as string with boolean literal. Use boolean comparison instead.",
			Context: match,
		})
	}

	// Check 3: String concatenation with template variables in ${}
	if match := stringPlusVar.FindString(line); match != "" {
		col := strings.Index(line, match) + 1
		issues = append(issues, Issue{
			File:    filePath,
			Line:    lineNum,
			Column:  col,
			Type:    "string-concatenation-in-expression",
			Message: "String concatenation in ${} expression may not work as expected. Consider handling in Go code.",
			Context: match,
		})
	}

	if match := varPlusString.FindString(line); match != "" {
		col := strings.Index(line, match) + 1
		issues = append(issues, Issue{
			File:    filePath,
			Line:    lineNum,
			Column:  col,
			Type:    "string-concatenation-in-expression",
			Message: "String concatenation in ${} expression may not work as expected. Consider handling in Go code.",
			Context: match,
		})
	}

	// Check 4: Double-quoted template variables (potential issue in JS context)
	if strings.Contains(line, "${") && doubleQuotedVar.FindString(line) != "" {
		// This might be intentional in some cases, so make it a warning
		match := doubleQuotedVar.FindString(line)
		col := strings.Index(line, match) + 1

		// Only warn if it's not inside a ${} expression (where it would be caught by check 1)
		if !isInsideExpression(line, col) {
			issues = append(issues, Issue{
				File:    filePath,
				Line:    lineNum,
				Column:  col,
				Type:    "double-quoted-var",
				Message: "Template variable is double-quoted. Verify this is intentional for the context.",
				Context: match,
			})
		}
	}

	// Check 5: Look for old-style string comparisons like '{{var}}' === 'value'
	if strings.Contains(line, "'{{") && strings.Contains(line, "}}'") && strings.Contains(line, "===") {
		match := singleQuotedVar.FindString(line)
		if match != "" {
			col := strings.Index(line, match) + 1
			issues = append(issues, Issue{
				File:    filePath,
				Line:    lineNum,
				Column:  col,
				Type:    "quoted-var-comparison",
				Message: "Template variable is quoted in comparison. This treats it as a string.",
				Context: line,
			})
		}
	}

	return issues
}

func isInsideExpression(line string, pos int) bool {
	// Check if position is inside ${}
	lastExprStart := strings.LastIndex(line[:pos], "${")
	lastExprEnd := strings.LastIndex(line[:pos], "}")

	return lastExprStart > lastExprEnd && lastExprStart != -1
}

func suggestFix(issue Issue) string {
	switch issue.Type {
	case "quoted-var-in-expression":
		return "Remove quotes around {{...}} inside ${} expressions. Use: ${{{var}}} instead of ${'{{var}}'}"
	case "string-bool-comparison":
		return "Change to boolean comparison: {{flag}} === true instead of '{{flag}}' === 'true'"
	case "string-concatenation-in-expression":
		return "Move string concatenation to Go code and pass the final value as a prop"
	case "double-quoted-var":
		return "If in HTML attribute context, this might be correct. If in JS context, verify it's intentional."
	case "quoted-var-comparison":
		return "Remove quotes for non-string comparisons: {{var}} === value instead of '{{var}}' === 'value'"
	default:
		return ""
	}
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
