package dictionarygen

import "testing"

func TestIdentifier(t *testing.T) {
	tbl := []struct {
		Word       string
		Identifier string
	}{
		{"", ""},
		{"User-Password", "UserPassword"},
		{"User_Password", "UserPassword"},
		{"user_password", "UserPassword"},
		{"expiry", "Expiry"},
		{"3Com-URL", "ThreeComURL"},
		{"3GPP-RAT-Type", "ThreeGPPRATType"},
	}

	for _, tt := range tbl {
		ident := identifier(tt.Word)
		if ident != tt.Identifier {
			t.Errorf("identifier(%s) = %s; expected %s", tt.Word, ident, tt.Identifier)
		}
	}
}
