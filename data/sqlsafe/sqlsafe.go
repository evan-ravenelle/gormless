package sqlsafe

import (
	"regexp"
	"strings"
)

func IsSafeSQLString(sqlText string) bool {
	if sqlText == "" {
		return false
	}

	// Check for common SQL injection patterns that would be malicious
	// even in a DDL context
	dangerousPatterns := []string{
		";\\s*DROP\\s+",           // Attempts to chain DROP commands
		";\\s*DELETE\\s+",         // Attempts to chain DELETE commands
		";\\s*TRUNCATE\\s+",       // Attempts to chain TRUNCATE commands
		"DROP\\s+DATABASE",        // Dropping an entire database is likely dangerous
		"DROP\\s+USER",            // Dropping users is dangerous
		"GRANT\\s+.*\\s+ALL",      // Granting all privileges is suspicious
		"--\\s*$",                 // Comment at end of line could hide malicious code
		"/\\*.*?\\*/",             // Multi-line comments that could hide code
		"EXECUTE\\s+", "EXEC\\s+", // Execution of dynamic SQL
		"xp_cmdshell", "sp_executesql", // SQL Server stored procedures for commands
		"INTO\\s+OUTFILE", "INTO\\s+DUMPFILE", // MySQL file operations
		"WAITFOR\\s+DELAY", "WAITFOR\\s+TIME", // Time-based attacks
		"0x[0-9A-F]{16,}", // Hex encoded strings (potential obfuscation)
	}

	// Check for dangerous patterns
	for _, pattern := range dangerousPatterns {
		re := regexp.MustCompile("(?i)" + pattern) // Case insensitive
		if re.MatchString(sqlText) {
			return false
		}
	}

	// Check for stacked queries (multiple statements)
	// Note: Some DDL operations legitimately use multiple statements,
	// so this check may need refinement for your specific use case
	if containsMultipleStatements(sqlText) && containsDangerousCombination(sqlText) {
		return false
	}

	sqlNoComments := removeComments(sqlText)

	// Check for attempts to modify system tables/views
	systemObjects := []string{
		"INFORMATION_SCHEMA\\.", "sys\\.", "pg_", "mysql\\.",
		"sqlite_master", "user_tables", "all_tables",
		"master\\.dbo", "msdb\\.dbo",
	}

	for _, obj := range systemObjects {
		re := regexp.MustCompile("(?i)" + obj)
		if re.MatchString(sqlNoComments) {
			// Allow read operations but block modifications
			lowerSQL := strings.ToLower(sqlNoComments)
			if strings.Contains(lowerSQL, "insert") ||
				strings.Contains(lowerSQL, "update") ||
				strings.Contains(lowerSQL, "delete") ||
				strings.Contains(lowerSQL, "drop") ||
				strings.Contains(lowerSQL, "alter") {
				return false
			}
		}
	}

	return true
}

// Helper function to check for multiple SQL statements
func containsMultipleStatements(sqlText string) bool {
	// Remove strings in single quotes to avoid false positives

	sqlNoStrings := removeQuotedStrings(sqlText)

	// Count semicolons not in comments
	sqlNoComments := removeComments(sqlNoStrings)

	// Split by semicolons and check if we have multiple non-empty statements
	statements := strings.Split(sqlNoComments, ";")
	count := 0
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			count++
		}
	}

	return count > 1
}

// Helper function to check for dangerous combinations of statements
func containsDangerousCombination(sqlText string) bool {
	lowerSQL := strings.ToLower(sqlText)

	// Check for combinations that might be suspicious
	return (strings.Contains(lowerSQL, "create") && strings.Contains(lowerSQL, "drop")) ||
		(strings.Contains(lowerSQL, "alter") && strings.Contains(lowerSQL, "drop")) ||
		(strings.Contains(lowerSQL, "select") && strings.Contains(lowerSQL, "into")) ||
		(strings.Contains(lowerSQL, "create") && strings.Contains(lowerSQL, "execute"))
}

// Helper function to remove quoted strings
func removeQuotedStrings(sqlText string) string {
	re := regexp.MustCompile("'[^']*'")
	noQuotesString := re.ReplaceAllString(sqlText, "")
	re = regexp.MustCompile("'")
	return re.ReplaceAllString(noQuotesString, "")
}

// Helper function to remove comments
func removeComments(sqlText string) string {
	// Remove single-line comments
	re1 := regexp.MustCompile("--.*$")
	noSingleComments := re1.ReplaceAllString(sqlText, "")

	// Remove multi-line comments
	re2 := regexp.MustCompile("/\\*[\\s\\S]*?\\*/")
	return re2.ReplaceAllString(noSingleComments, "")
}
