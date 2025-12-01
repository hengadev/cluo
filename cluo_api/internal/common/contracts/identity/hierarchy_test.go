package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleLevel(t *testing.T) {
	tests := []struct {
		role          Role
		expectedLevel int
		name          string
	}{
		{Guest, 0, "guest has level 0"},
		{Client, 1, "client has level 1"},
		{Administrator, 2, "administrator has level 2"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedLevel, test.role.Level())
		})
	}

	t.Run("invalid role should return -1", func(t *testing.T) {
		invalidRole := Role(99)
		assert.Equal(t, -1, invalidRole.Level())
	})
}

func TestIsAtLeast(t *testing.T) {
	t.Run("role hierarchy comparisons", func(t *testing.T) {
		tests := []struct {
			current  Role
			target   Role
			expected bool
			name     string
		}{
			{Administrator, Guest, true, "admin >= guest"},
			{Administrator, Administrator, true, "admin >= admin"},
			{Client, Guest, true, "client >= guest"},
			{Guest, Client, false, "guest < client"},
			{Guest, Administrator, false, "guest < admin"},
			{Client, Administrator, false, "client < admin"},
			{Administrator, Client, true, "admin >= client"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				assert.Equal(t, test.expected, test.current.IsAtLeast(test.target))
			})
		}
	})
}

func TestHasPermission(t *testing.T) {
	t.Run("permission checks should work like IsAtLeast", func(t *testing.T) {
		assert.True(t, Administrator.HasPermission(Client))
		assert.True(t, Client.HasPermission(Guest))
		assert.False(t, Guest.HasPermission(Client))
		assert.True(t, Administrator.HasPermission(Guest))
		assert.True(t, Administrator.HasPermission(Administrator))
	})
}
