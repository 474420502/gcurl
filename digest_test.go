package gcurl

import (
	"testing"
)

func TestDigestAuthentication(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantUser string
		wantPass string
		wantErr  bool
	}{
		{
			name:     "Valid digest credentials",
			input:    "testuser:testpass",
			wantUser: "testuser",
			wantPass: "testpass",
			wantErr:  false,
		},
		{
			name:     "Password with colon",
			input:    "user:pass:with:colons",
			wantUser: "user",
			wantPass: "pass:with:colons",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			input:    "user:",
			wantUser: "user",
			wantPass: "",
			wantErr:  false,
		},
		{
			name:     "Invalid format - no colon",
			input:    "userpass",
			wantUser: "",
			wantPass: "",
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    "",
			wantUser: "",
			wantPass: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{}
			err := handleDigest(c, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("handleDigest() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("handleDigest() unexpected error: %v", err)
				return
			}

			if c.AuthV2 == nil {
				t.Errorf("handleDigest() AuthV2 is nil")
				return
			}

			if c.AuthV2.Type != AuthDigest {
				t.Errorf("handleDigest() AuthType = %v, want %v", c.AuthV2.Type, AuthDigest)
			}

			if c.AuthV2.Username != tt.wantUser {
				t.Errorf("handleDigest() Username = %v, want %v", c.AuthV2.Username, tt.wantUser)
			}

			if c.AuthV2.Password != tt.wantPass {
				t.Errorf("handleDigest() Password = %v, want %v", c.AuthV2.Password, tt.wantPass)
			}
		})
	}
}

func TestDigestOptionParsing(t *testing.T) {
	tests := []struct {
		name       string
		curlString string
		wantUser   string
		wantPass   string
		wantErr    bool
	}{
		{
			name:       "Digest with long option",
			curlString: `curl --digest user:pass https://example.com`,
			wantUser:   "user",
			wantPass:   "pass",
			wantErr:    false,
		},
		{
			name:       "Digest with complex password",
			curlString: `curl --digest "admin:p@ssw0rd:with:symbols" https://example.com`,
			wantUser:   "admin",
			wantPass:   "p@ssw0rd:with:symbols",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.curlString)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if c.AuthV2 == nil {
				t.Errorf("Parse() AuthV2 is nil")
				return
			}

			if c.AuthV2.Type != AuthDigest {
				t.Errorf("Parse() AuthType = %v, want %v", c.AuthV2.Type, AuthDigest)
			}

			if c.AuthV2.Username != tt.wantUser {
				t.Errorf("Parse() Username = %v, want %v", c.AuthV2.Username, tt.wantUser)
			}

			if c.AuthV2.Password != tt.wantPass {
				t.Errorf("Parse() Password = %v, want %v", c.AuthV2.Password, tt.wantPass)
			}
		})
	}
}

func TestDigestAuthenticationMethods(t *testing.T) {
	auth := NewDigestAuth("testuser", "testpass")

	// Test type
	if auth.Type != AuthDigest {
		t.Errorf("NewDigestAuth() Type = %v, want %v", auth.Type, AuthDigest)
	}

	// Test credentials
	if auth.Username != "testuser" {
		t.Errorf("NewDigestAuth() Username = %v, want %v", auth.Username, "testuser")
	}

	if auth.Password != "testpass" {
		t.Errorf("NewDigestAuth() Password = %v, want %v", auth.Password, "testpass")
	}

	// Test header generation
	header := auth.GetAuthHeader()
	if header != "" {
		t.Errorf("NewDigestAuth() GetAuthHeader() = %v, want empty string (digest requires challenge)", header)
	}
}
