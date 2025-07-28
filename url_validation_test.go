package gcurl

import (
	"net/url"
	"testing"
)

// TestURLValidation 测试URL验证逻辑
func TestURLValidation(t *testing.T) {
	testCases := []struct {
		urlStr      string
		shouldError bool
		description string
	}{
		{"http://httpbin.org/get", false, "valid HTTP URL"},
		{"https://example.com", false, "valid HTTPS URL"},
		{"ftp://example.com/file.txt", false, "valid FTP URL"},
		{"not-a-valid-url", false, "Go's url.Parse is very permissive"}, // Go认为这是有效的
		{"", false, "empty URL - Go accepts it"},
		{" ", false, "space is considered valid by Go"},
		{"://invalid", true, "malformed scheme"},
		{"http://", false, "incomplete but valid by Go standards"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parsedURL, err := url.Parse(tc.urlStr)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error for URL '%s', but got none. Parsed: %+v", tc.urlStr, parsedURL)
				} else {
					t.Logf("✓ Correctly rejected URL '%s': %v", tc.urlStr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success for URL '%s', but got error: %v", tc.urlStr, err)
				} else {
					t.Logf("✓ Accepted URL '%s' - Scheme: '%s', Host: '%s', Path: '%s'",
						tc.urlStr, parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
				}
			}
		})
	}
}

// TestEnhancedURLValidation 测试更严格的URL验证
func TestEnhancedURLValidation(t *testing.T) {
	testCases := []struct {
		urlStr        string
		shouldBeValid bool
		description   string
	}{
		{"http://httpbin.org/get", true, "valid HTTP URL"},
		{"https://example.com", true, "valid HTTPS URL"},
		{"ftp://example.com/file.txt", true, "valid FTP URL"},
		{"not-a-valid-url", false, "should be rejected - no scheme"},
		{"", false, "empty URL"},
		{" ", false, "whitespace URL"},
		{"://invalid", false, "malformed scheme"},
		{"http://", false, "incomplete URL - no host"},
		{"file:///path/to/file", true, "valid file URL"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			isValid := isValidURL(tc.urlStr)

			if tc.shouldBeValid && !isValid {
				t.Errorf("URL '%s' should be valid but was rejected", tc.urlStr)
			} else if !tc.shouldBeValid && isValid {
				t.Errorf("URL '%s' should be invalid but was accepted", tc.urlStr)
			} else {
				t.Logf("✓ URL '%s' validation result: %v (expected: %v)", tc.urlStr, isValid, tc.shouldBeValid)
			}
		})
	}
}
