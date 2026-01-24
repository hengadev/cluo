package config

// VaultConfig holds HashiCorp Vault configuration.
type VaultConfig struct {
	Address   string
	Token     string
	MountPath string
	// AppRole authentication (production)
	AppRoleID       string
	AppRoleSecretID string
}

// UseAppRole returns true if AppRole authentication should be used.
func (c VaultConfig) UseAppRole() bool {
	return c.AppRoleID != "" && c.AppRoleSecretID != ""
}

func loadVaultConfig() (VaultConfig, error) {
	return VaultConfig{
		Address:         getEnv("VAULT_ADDR", ""),
		Token:           getEnv("VAULT_TOKEN", ""),
		MountPath:       getEnv("CLUO_VAULT_MOUNT_PATH", "secret"),
		AppRoleID:       getEnv("VAULT_APPROLE_ROLE_ID", ""),
		AppRoleSecretID: getEnv("VAULT_APPROLE_SECRET_ID", ""),
	}, nil
}
