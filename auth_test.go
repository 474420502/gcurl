package gcurl

import (
	"testing"
)

func TestAuthType_String(t *testing.T) {
	if AuthBasic.String() != "Basic" {
		t.Errorf("expected 'Basic', got '%s'", AuthBasic.String())
	}
	if AuthDigest.String() != "Digest" {
		t.Errorf("expected 'Digest', got '%s'", AuthDigest.String())
	}
	if AuthBearer.String() != "Bearer" {
		t.Errorf("expected 'Bearer', got '%s'", AuthBearer.String())
	}
	if AuthNTLM.String() != "NTLM" {
		t.Errorf("expected 'NTLM', got '%s'", AuthNTLM.String())
	}
	var unknown AuthType = 100
	if unknown.String() != "Unknown" {
		t.Errorf("expected 'Unknown', got '%s'", unknown.String())
	}
}

func TestNewBasicAuth(t *testing.T) {
	auth := NewBasicAuth("user", "pass")
	if auth.Type != AuthBasic {
		t.Errorf("expected AuthBasic, got %v", auth.Type)
	}
	if !auth.IsValid() {
		t.Error("expected valid basic auth")
	}
	header := auth.GetAuthHeader()
	if header != "" {
		t.Error("expected empty header for basic auth")
	}
}

func TestNewDigestAuth(t *testing.T) {
	auth := NewDigestAuth("user", "pass")
	if auth.Type != AuthDigest {
		t.Errorf("expected AuthDigest, got %v", auth.Type)
	}
	if !auth.IsValid() {
		t.Error("expected valid digest auth")
	}
	header := auth.GetAuthHeader()
	if header != "" {
		t.Error("expected empty header for digest auth")
	}
}

func TestNewBearerAuth(t *testing.T) {
	auth := NewBearerAuth("token")
	if auth.Type != AuthBearer {
		t.Errorf("expected AuthBearer, got %v", auth.Type)
	}
	if !auth.IsValid() {
		t.Error("expected valid bearer auth")
	}
	header := auth.GetAuthHeader()
	if header != "Bearer token" {
		t.Errorf("expected 'Bearer token', got '%s'", header)
	}
}

func TestAuthentication_String(t *testing.T) {
	basic := NewBasicAuth("user", "pass")
	digest := NewDigestAuth("user", "pass")
	bearer := NewBearerAuth("token")
	ntlm := &Authentication{Type: AuthNTLM, Username: "ntlmuser"}
	unknown := &Authentication{Type: 100, Username: "x"}
	_ = basic.String()
	_ = digest.String()
	_ = bearer.String()
	_ = ntlm.String()
	_ = unknown.String()
	var nilAuth *Authentication
	_ = nilAuth.String()
}
