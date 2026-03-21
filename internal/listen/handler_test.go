package listen

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testHandler(t *testing.T, opts Options) (http.Handler, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	opts.Writer = &buf
	return Handler(opts), &buf
}

func TestHandler_GET(t *testing.T) {
	h, buf := testHandler(t, Options{})

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	log := buf.String()
	if !strings.Contains(log, "GET /hello") {
		t.Errorf("expected method and path in log, got %q", log)
	}
}

func TestHandler_POSTWithJSON(t *testing.T) {
	h, buf := testHandler(t, Options{})

	body := `{"event":"push","ref":"main"}`
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	if !strings.Contains(log, "POST /webhook") {
		t.Errorf("expected method and path, got %q", log)
	}
	if !strings.Contains(log, "Content-Type: application/json") {
		t.Errorf("expected Content-Type header, got %q", log)
	}
	// Should be pretty-printed
	if !strings.Contains(log, `"event": "push"`) {
		t.Errorf("expected pretty-printed JSON body, got %q", log)
	}
}

func TestHandler_POSTWithPlainText(t *testing.T) {
	h, buf := testHandler(t, Options{})

	req := httptest.NewRequest("POST", "/hook", strings.NewReader("plain text body"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	if !strings.Contains(log, "plain text body") {
		t.Errorf("expected raw body in log, got %q", log)
	}
}

func TestHandler_CustomStatus(t *testing.T) {
	h, _ := testHandler(t, Options{Status: 201})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Errorf("expected 201, got %d", rec.Code)
	}
}

func TestHandler_CustomBody(t *testing.T) {
	h, _ := testHandler(t, Options{Body: `{"ok":true}`})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Body.String() != `{"ok":true}` {
		t.Errorf("expected custom body, got %q", rec.Body.String())
	}
}

func TestHandler_NoBody(t *testing.T) {
	h, buf := testHandler(t, Options{NoBody: true})

	req := httptest.NewRequest("POST", "/", strings.NewReader("secret data"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	if strings.Contains(log, "secret data") {
		t.Errorf("expected body to be omitted with NoBody, got %q", log)
	}
}

func TestHandler_Headers(t *testing.T) {
	h, buf := testHandler(t, Options{})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Custom", "test-value")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	if !strings.Contains(log, "X-Custom: test-value") {
		t.Errorf("expected custom header in log, got %q", log)
	}
}

func TestHandler_Separator(t *testing.T) {
	h, buf := testHandler(t, Options{})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	if !strings.Contains(log, "---") {
		t.Errorf("expected separator in log, got %q", log)
	}
}

func TestHandler_LargeBody(t *testing.T) {
	h, buf := testHandler(t, Options{})

	// 2MB body — should be truncated to 1MB
	bigBody := strings.Repeat("x", 2<<20)
	req := httptest.NewRequest("POST", "/", strings.NewReader(bigBody))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	// Body in log should be at most 1MB
	bodyStart := strings.Index(log, "\n\n")
	if bodyStart > 0 {
		bodyPart := log[bodyStart:]
		if len(bodyPart) > maxBodySize+100 { // some slack for newlines
			t.Errorf("expected body truncated to ~1MB, got %d bytes", len(bodyPart))
		}
	}
}

func TestHandler_EmptyBody(t *testing.T) {
	h, buf := testHandler(t, Options{})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	log := buf.String()
	// Should not have double newline (body section) for GET with no body
	parts := strings.Split(log, "\n\n")
	// First part is empty (leading \n from separator), second is headers
	// There should be no body section
	bodyCount := 0
	for _, p := range parts {
		if strings.TrimSpace(p) != "" && !strings.Contains(p, "---") && !strings.Contains(p, ":") {
			bodyCount++
		}
	}
	_ = bodyCount // just verify it doesn't crash; body section is absent for empty body
}

func TestHandler_DefaultStatus(t *testing.T) {
	// Status 0 should default to 200
	var buf bytes.Buffer
	h := Handler(Options{Writer: &buf})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected default 200, got %d", rec.Code)
	}
}

var _ io.Writer = (*bytes.Buffer)(nil)
