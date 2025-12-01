package identity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleConstants(t *testing.T) {
	t.Run("string constants should match role strings", func(t *testing.T) {
		assert.Equal(t, "guest", GuestStr)
		assert.Equal(t, "client", ClientStr)
		assert.Equal(t, "administrator", AdministratorStr)
	})

	t.Run("role string should match constants", func(t *testing.T) {
		assert.Equal(t, GuestStr, Guest.String())
		assert.Equal(t, ClientStr, Client.String())
		assert.Equal(t, AdministratorStr, Administrator.String())
	})
}
