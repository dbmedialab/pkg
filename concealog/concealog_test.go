package concealog

import (
	"testing"
)

// TestReplacer - ensures that the replacer works as expected
func TestReplacer(t *testing.T) {
	input := `GET /capi/v1/client/customer/ticket/15497076 HTTP/1.1
Host: api.mediaconnect.no
User-Agent: Go-http-client/1.1
Accept: application/json
Authorization: Bearer 31e675cc-8ac7-4d18-a0fa-f4cd2e74a28a
Accept-Encoding: gzip`

	expected := `GET /capi/v1/client/customer/ticket/15497076 HTTP/1.1
Host: api.mediaconnect.no
User-Agent: Go-http-client/1.1
Accept: application/json
Authorization: Bearer *********
Accept-Encoding: gzip`

	ar, err := NewAuthReplacer()
	if err != nil {
		t.Errorf("Failed to create AuthReplacer: %w", err)
		return
	}

	res := ar.ReplaceString(input)
	if res != expected {
		t.Errorf("Unexpected output!\nExpected:\n%s\n\nGot:\n%s\n\n", expected, res)
		return
	}
}
