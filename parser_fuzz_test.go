//go:build go1.18
// +build go1.18

package gcurl

import (
	"testing"
)

// FuzzParse a an
func FuzzParse(f *testing.F) {
	testcases := []string{
		`curl "http://example.com"`,
		`curl -X POST --data '{"a":1}' "http://example.com"`,
		`curl 'https://www.google.com' -H 'Accept: text/html' --compressed`,
		`curl -I -H 'X-Test: true' "http://httpbin.org/headers"`,
		`curl "https://a.b/c'd"`,
		`curl --user "name:pass"`,
		`curl -x socks5://user:pass@host:port/path`,
		`curl --data-binary @/etc/passwd`,
		`curl 'a'b'c'd'e'f'g'h'i'j'k'l'm'n'o'p'q'r's't'u'v'w'x'y'z`,
		`curl "a"b"c"d"e"f"g"h"i"j"k"l"m"n"o"p"q"r"s"t"u"v"w"x"y"z`,
		`curl a'b"c'd"e'f"g'h"i'j"k'l"m'n"o'p"q'r"s't"u'v"w'x"y'z`,
	}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		// The fuzzer will generate random inputs based on the seed corpus.
		// We just need to ensure that parsing these inputs does not cause a panic.
		// The function should either return a valid CURL object and no error,
		// or a nil CURL object and an error.
		_, _ = Parse(orig)
	})
}
