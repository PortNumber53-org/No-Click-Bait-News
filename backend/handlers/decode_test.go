package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PortNumber53/no-click-bait-news/backend/models"
)

// Regression: the mobile app ships fields (e.g. remember_me) that the backend
// may predate. An unknown field must not turn a valid login into a 400.
func TestDecodeJSONIgnoresUnknownFields(t *testing.T) {
	body := `{"email":"u@example.com","password":"12345678","remember_me":true,"future_field":"x"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	var req models.LoginRequest
	if !DecodeJSON(w, r, &req) {
		t.Fatalf("DecodeJSON rejected a login body with unknown fields: %s", w.Body.String())
	}
	if req.Email != "u@example.com" || req.Password != "12345678" || !req.RememberMe {
		t.Fatalf("decoded unexpected values: %+v", req)
	}
}

func TestDecodeJSONRejectsMalformedBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":`))
	w := httptest.NewRecorder()
	var req models.LoginRequest
	if DecodeJSON(w, r, &req) {
		t.Fatal("DecodeJSON accepted malformed JSON")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDecodeJSONRejectsMultipleObjects(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"a"} {"email":"b"}`))
	w := httptest.NewRecorder()
	var req models.LoginRequest
	if DecodeJSON(w, r, &req) {
		t.Fatal("DecodeJSON accepted multiple JSON objects")
	}
}
