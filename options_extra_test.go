package gcurl

import (
	"testing"
)

func TestHandleVerbose(t *testing.T) {
	c := New()
	err := handleVerbose(c)
	if err != nil {
		t.Errorf("handleVerbose should not error, got %v", err)
	}
}

func TestHandleInclude(t *testing.T) {
	c := New()
	err := handleInclude(c)
	if err != nil {
		t.Errorf("handleInclude should not error, got %v", err)
	}
}

func TestHandleSilent(t *testing.T) {
	c := New()
	err := handleSilent(c)
	if err != nil {
		t.Errorf("handleSilent should not error, got %v", err)
	}
}

func TestHandleTrace(t *testing.T) {
	c := New()
	err := handleTrace(c)
	if err != nil {
		t.Errorf("handleTrace should not error, got %v", err)
	}
}
