package identity

import (
	"errors"
	"testing"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToRole(t *testing.T) {
	tests := []struct {
		roleStr  string
		expected Role
		isOk     bool
		name     string
	}{
		{roleStr: "guest", expected: Guest, name: "Convert to GUEST", isOk: true},
		{roleStr: "client", expected: Client, name: "Convert to CLIENT", isOk: true},
		{roleStr: "administrator", expected: Administrator, name: "Convert to ADMINISTRATOR", isOk: true},
		{roleStr: "random_value", expected: Guest, name: "Convert to GUEST", isOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ConvertToRole(tt.roleStr)
			assert.Equal(t, tt.expected, got)
			assert.Equal(t, tt.isOk, ok)
		})
	}
}

func TestParseRole(t *testing.T) {
	t.Run("valid roles should parse correctly", func(t *testing.T) {
		tests := []struct {
			input    string
			expected Role
		}{
			{"guest", Guest},
			{"client", Client},
			{"administrator", Administrator},
			{"GUEST", Guest},                 // Case insensitive
			{"Client", Client},               // Mixed case
			{"  guest  ", Guest},             // With whitespace
			{"ADMINISTRATOR", Administrator}, // Upper case
		}

		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				role, err := ParseRole(test.input)
				require.NoError(t, err)
				assert.Equal(t, test.expected, role)
			})
		}
	})

	t.Run("invalid roles should return error", func(t *testing.T) {
		invalidRoles := []string{
			"invalid",
			"",
			"admin", // Not "administrator"
			"user",
			"123",
			"null",
		}

		for _, invalid := range invalidRoles {
			t.Run(invalid, func(t *testing.T) {
				role, err := ParseRole(invalid)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, errs.ErrInvalidValue))
				assert.Equal(t, Guest, role) // Should default to Guest on error
			})
		}
	})
}
