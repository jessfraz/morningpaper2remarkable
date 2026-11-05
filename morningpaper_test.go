package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTryToFindPDFLinkHTMLDepth(t *testing.T) {
	cases := []struct {
		name           string
		body           string
		wantParseError bool
	}{
		{name: "ordinary HTML", body: "<html><body><p>A paper summary.</p></body></html>"},
		{name: "excessive HTML depth", body: strings.Repeat("<div>", 4096) + strings.Repeat("</div>", 4096), wantParseError: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if _, err := io.WriteString(w, tc.body); err != nil {
					t.Errorf("writing fixture HTML: %v", err)
				}
			}))
			defer server.Close()

			_, err := tryToFindPDFLink(server.URL)
			if tc.wantParseError {
				if err == nil || !strings.Contains(err.Error(), "parsing link") {
					t.Fatalf("excessively nested HTML should return a parser error, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("ordinary HTML should parse successfully: %v", err)
			}
		})
	}
}
