package sqlsafe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSafeSQLIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid identifier",
			input:    "valid_column_name",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "SQL injection attempt with semicolon",
			input:    "column_name;DROP TABLE users;",
			expected: false,
		},
		{
			name:     "SQL injection attempt with comment",
			input:    "column_name--comment",
			expected: true,
		},
		{
			name:     "SQL injection with quotes",
			input:    "column_name'",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSafeSQLString(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}
