package main

import (
	"bytes"

	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.gregor-pifko/internal/models/mocks"
)

// Regular expression which captures the CSRF token value from the html pages signup/login/...
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)

	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

// Return an instance of our application struct containing mocked dependencies
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		infoLog:        log.New(io.Discard, "", 0),
		errorLog:       log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		user:           &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: scs.New(),
	}
}

// Define a custom testServer type which embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// Initialize and return a new instance of our custom testServer type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// Disable redirect-following for the test server client
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Makes a GET request to a given url path, and returns the response status code, headers and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	resp, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return resp.StatusCode, resp.Header, string(body)
}

// Makes a POST request to test server with specific form data
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	resp, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return resp.StatusCode, resp.Header, string(body)
}
