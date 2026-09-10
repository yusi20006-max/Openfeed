package telemirror

import (
	"net/http/httptest"
	"testing"
)

func TestSetImageHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/image.jpg", nil)
	setImageHeaders(req, "test-agent", true)

	if got := req.Header.Get("Accept"); got == "" || got == "text/html" {
		t.Fatalf("Accept = %q, want image-oriented Accept header", got)
	}
	if got := req.Header.Get("Sec-Fetch-Dest"); got != "image" {
		t.Fatalf("Sec-Fetch-Dest = %q, want image", got)
	}
	if got := req.Header.Get("Referer"); got != "https://translate.google.com/" {
		t.Fatalf("Referer = %q, want translate.google.com", got)
	}
}

func TestImageContentType(t *testing.T) {
	jpeg := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F'}
	if got := imageContentType(jpeg, ""); got == "" || got[:6] != "image/" {
		t.Fatalf("JPEG content type = %q, want image/*", got)
	}

	html := []byte("<!doctype html><html><body>error</body></html>")
	if got := imageContentType(html, "text/html"); got != "" {
		t.Fatalf("HTML content type = %q, want rejection", got)
	}

	if got := imageContentType([]byte("not an image"), "application/json"); got != "" {
		t.Fatalf("JSON/error body content type = %q, want rejection", got)
	}
}
