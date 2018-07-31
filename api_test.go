package mpesa

import (
	"testing"
)

var testMPESAService = New("MY4vqhkg0Rlkj8eGTmyWmlUm9whiIKho", "jUhEL2xGhGLz3tua", "")

func TestService_GenerateNewAccessToken(t *testing.T) {
	token, err := testMPESAService.GenerateNewAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}
