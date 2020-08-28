package gcurl

import "testing"

func TestDomain(t *testing.T) {
	if !isCookieDomainName("www.baidu.com") {
		t.Error("isCookieDomainName error")
	}

	if !validCookieDomain("127.0.0.1") {
		t.Error("validCookieDomain error")
	}
}
