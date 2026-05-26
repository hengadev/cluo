package validation

import "fmt"

// activeISO4217Codes is a hard-coded allowlist of all active ISO 4217
// currency codes.  This list covers ~170 codes as required by the issue.
var activeISO4217Codes = map[string]struct{}{
	"AED": {}, "AFN": {}, "ALL": {}, "AMD": {}, "ANG": {}, "AOA": {},
	"ARS": {}, "AUD": {}, "AWG": {}, "AZN": {}, "BAM": {}, "BBD": {},
	"BDT": {}, "BGN": {}, "BHD": {}, "BIF": {}, "BMD": {}, "BND": {},
	"BOB": {}, "BOV": {}, "BRL": {}, "BSD": {}, "BTN": {}, "BWP": {},
	"BYN": {}, "BZD": {}, "CAD": {}, "CDF": {}, "CHE": {}, "CHF": {},
	"CHW": {}, "CLF": {}, "CLP": {}, "CNY": {}, "COP": {}, "COU": {},
	"CRC": {}, "CUP": {}, "CVE": {}, "CZK": {}, "DJF": {}, "DKK": {},
	"DOP": {}, "DZD": {}, "EGP": {}, "ERN": {}, "ETB": {}, "EUR": {},
	"FJD": {}, "FKP": {}, "GBP": {}, "GEL": {}, "GHS": {}, "GIP": {},
	"GMD": {}, "GNF": {}, "GTQ": {}, "GYD": {}, "HKD": {}, "HNL": {},
	"HRK": {}, "HTG": {}, "HUF": {}, "IDR": {}, "ILS": {}, "INR": {},
	"IQD": {}, "IRR": {}, "ISK": {}, "JMD": {}, "JOD": {}, "JPY": {},
	"KES": {}, "KGS": {}, "KHR": {}, "KMF": {}, "KPW": {}, "KRW": {},
	"KWD": {}, "KYD": {}, "KZT": {}, "LAK": {}, "LBP": {}, "LKR": {},
	"LRD": {}, "LSL": {}, "LYD": {}, "MAD": {}, "MDL": {}, "MGA": {},
	"MKD": {}, "MMK": {}, "MNT": {}, "MOP": {}, "MRU": {}, "MUR": {},
	"MVR": {}, "MWK": {}, "MXN": {}, "MXV": {}, "MYR": {}, "MZN": {},
	"NAD": {}, "NGN": {}, "NIO": {}, "NOK": {}, "NPR": {}, "NZD": {},
	"OMR": {}, "PAB": {}, "PEN": {}, "PGK": {}, "PHP": {}, "PKR": {},
	"PLN": {}, "PYG": {}, "QAR": {}, "RON": {}, "RSD": {}, "RUB": {},
	"RWF": {}, "SAR": {}, "SBD": {}, "SCR": {}, "SDG": {}, "SEK": {},
	"SGD": {}, "SHP": {}, "SLE": {}, "SOS": {}, "SRD": {}, "SSP": {},
	"STN": {}, "SVC": {}, "SYP": {}, "SZL": {}, "THB": {}, "TJS": {},
	"TMT": {}, "TND": {}, "TOP": {}, "TRY": {}, "TTD": {}, "TWD": {},
	"TZS": {}, "UAH": {}, "UGX": {}, "USD": {}, "USN": {}, "UYI": {},
	"UYU": {}, "UYW": {}, "UZS": {}, "VED": {}, "VES": {}, "VND": {},
	"VUV": {}, "WST": {}, "XAF": {}, "XAG": {}, "XAU": {}, "XBA": {},
	"XBB": {}, "XBC": {}, "XBD": {}, "XCD": {}, "XDR": {}, "XOF": {},
	"XPD": {}, "XPF": {}, "XPT": {}, "XSU": {}, "XUA": {}, "YER": {},
	"ZAR": {}, "ZMW": {}, "ZWL": {},
}

// IsValidCurrency returns true if code is a recognised active ISO 4217 code.
func IsValidCurrency(code string) bool {
	_, ok := activeISO4217Codes[code]
	return ok
}

// ValidateCurrency returns an error if the code is not a recognised ISO 4217 code.
func ValidateCurrency(code string) error {
	if !IsValidCurrency(code) {
		return fmt.Errorf("invalid ISO 4217 currency code: %s", code)
	}
	return nil
}
