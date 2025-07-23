package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.gregor-pifko/internal/assert"
)

func TestPing(t *testing.T) {
	app := &application{
		infoLog:  log.New(io.Discard, "", 0),
		errorLog: log.New(io.Discard, "", 0),
	}

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
