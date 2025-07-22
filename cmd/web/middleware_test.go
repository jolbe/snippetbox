package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.gregor-pifko/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	// Create a test request/response recorder
	req := httptest.NewRequest(http.MethodGet, "/secureHeadersMiddleware", nil)
	rec := httptest.NewRecorder()

	// Mock handler to test if middleware calls it
	mockNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock handler in the middleware and call the handler with our request/response recorders
	secureHeaders(mockNext).ServeHTTP(rec, req)

	// Get the resp
	resp := rec.Result()

	// secureHeaders() middleware sets all the expected headers on the HTTP response
	assert.Equal(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equal(t, resp.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, resp.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, resp.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, resp.Header.Get("X-XSS-Protection"), "0")

	// secureHeaders() middleware correctly calls the next handler in the chain
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
