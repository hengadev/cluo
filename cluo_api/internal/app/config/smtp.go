package config

// SMTPConfig holds SMTP mail server configuration.
type SMTPConfig struct {
	Host            string
	Port            string
	From            string
	Username        string
	Password        string
	PortalPublicURL string
}

func loadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:            getEnv("SMTP_HOST", ""),
		Port:            getEnv("SMTP_PORT", "587"),
		From:            getEnv("SMTP_FROM", ""),
		Username:        getEnv("SMTP_USERNAME", ""),
		Password:        getEnv("SMTP_PASSWORD", ""),
		PortalPublicURL: getEnv("PORTAL_PUBLIC_URL", ""),
	}
}

// IsConfigured returns true when all required SMTP fields are present.
func (c SMTPConfig) IsConfigured() bool {
	return c.Host != "" && c.From != ""
}
