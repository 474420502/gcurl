package gcurl

import (
	"fmt"
	"testing"
)

func TestQuoteFixComparison(t *testing.T) {
	testCases := []struct {
		name           string
		curlCommand    string
		expectedHeader string
		expectedValue  string
	}{
		{
			name:           "sec-ch-ua with complex quotes",
			curlCommand:    `curl 'https://example.com' -H 'sec-ch-ua: "Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"'`,
			expectedHeader: "sec-ch-ua",
			expectedValue:  `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
		},
		{
			name:           "sec-ch-ua-platform with quotes",
			curlCommand:    `curl 'https://example.com' -H 'sec-ch-ua-platform: "Windows"'`,
			expectedHeader: "sec-ch-ua-platform",
			expectedValue:  `"Windows"`,
		},
		{
			name:           "regular header without quotes",
			curlCommand:    `curl 'https://example.com' -H 'accept: application/json'`,
			expectedHeader: "accept",
			expectedValue:  `application/json`,
		},
		{
			name:           "header with intentional quotes",
			curlCommand:    `curl 'https://example.com' -H 'custom-header: "value with spaces"'`,
			expectedHeader: "custom-header",
			expectedValue:  `"value with spaces"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cu, err := Parse(tc.curlCommand)
			if err != nil {
				t.Fatalf("Failed to parse curl command: %v", err)
			}

			actualValue := cu.Header.Get(tc.expectedHeader)
			if actualValue != tc.expectedValue {
				t.Errorf("Header %s mismatch:\nExpected: %q\nActual:   %q",
					tc.expectedHeader, tc.expectedValue, actualValue)
			} else {
				fmt.Printf("✓ %s: %s = %q\n", tc.name, tc.expectedHeader, actualValue)
			}
		})
	}
}
