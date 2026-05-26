package validation

import "testing"

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"USD", "USD", true},
		{"EUR", "EUR", true},
		{"GBP", "GBP", true},
		{"CAD", "CAD", true},
		{"JPY", "JPY", true},
		{"CHF", "CHF", true},
		{"lowercase rejected", "usd", false},
		{"empty rejected", "", false},
		{"unknown rejected", "XYZ", false},
		{"too short rejected", "US", false},
		{"too long rejected", "USDA", false},
		{"XAF", "XAF", true},
		{"ZWL", "ZWL", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCurrency(tt.code); got != tt.want {
				t.Errorf("IsValidCurrency(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	if err := ValidateCurrency("USD"); err != nil {
		t.Errorf("ValidateCurrency(\"USD\") returned unexpected error: %v", err)
	}
	if err := ValidateCurrency("XYZ"); err == nil {
		t.Error("ValidateCurrency(\"XYZ\") should have returned an error")
	}
}
