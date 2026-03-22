package serve

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello world"), 0o644)
	os.WriteFile(filepath.Join(dir, "style.css"), []byte("body { color: red; }"), 0o644)

	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("nested content"), 0o644)

	return dir
}

func setupTestDirWithIndex(t *testing.T) string {
	t.Helper()
	dir := setupTestDir(t)
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>Home</h1>"), 0o644)
	return dir
}

func testHandler(t *testing.T, opts Options) http.Handler {
	t.Helper()
	if opts.LogWriter == nil {
		opts.LogWriter = io.Discard
	}
	return Handler(opts)
}

func request(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// --- File serving ---

func TestServeFile(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/hello.txt")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello world" {
		t.Errorf("expected 'hello world', got %q", rec.Body.String())
	}
}

func TestServeFile_ContentType(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/style.css")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/css") {
		t.Errorf("expected text/css content type, got %q", ct)
	}
}

func TestServeFile_Nested(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/sub/nested.txt")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "nested content" {
		t.Errorf("expected 'nested content', got %q", rec.Body.String())
	}
}

func TestServeFile_NotFound(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/missing.txt")
	if rec.Code != 404 {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestServeFile_HeadMethod(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "HEAD", "/hello.txt")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body for HEAD, got %d bytes", rec.Body.Len())
	}
}

func TestServeFile_MethodNotAllowed(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "POST", "/hello.txt")
	if rec.Code != 405 {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

// --- Index resolution ---

func TestIndex_HTML(t *testing.T) {
	dir := setupTestDirWithIndex(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<h1>Home</h1>") {
		t.Errorf("expected index.html content, got %q", rec.Body.String())
	}
}

func TestIndex_HTM_Fallback(t *testing.T) {
	dir := setupTestDir(t) // no index.html
	os.WriteFile(filepath.Join(dir, "index.htm"), []byte("<h1>HTM</h1>"), 0o644)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<h1>HTM</h1>") {
		t.Errorf("expected index.htm content, got %q", rec.Body.String())
	}
}

func TestIndex_HTML_TakesPrecedence(t *testing.T) {
	dir := setupTestDir(t)
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("HTML"), 0o644)
	os.WriteFile(filepath.Join(dir, "index.htm"), []byte("HTM"), 0o644)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	if rec.Body.String() != "HTML" {
		t.Errorf("expected index.html to take precedence, got %q", rec.Body.String())
	}
}

// --- Directory listing ---

func TestDirListing(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "hello.txt") {
		t.Errorf("expected listing to contain hello.txt")
	}
	if !strings.Contains(body, "sub/") {
		t.Errorf("expected listing to contain sub/")
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestDirListing_NoParentAtRoot(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	body := rec.Body.String()
	if strings.Contains(body, `href="../"`) {
		t.Errorf("expected no parent link at root")
	}
}

func TestDirListing_HasParentInSubdir(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/sub/")
	body := rec.Body.String()
	if !strings.Contains(body, `href="../"`) {
		t.Errorf("expected parent link in subdirectory listing")
	}
}

func TestDirListing_NoIndex(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir, NoIndex: true})

	rec := request(t, h, "GET", "/")
	if rec.Code != 404 {
		t.Errorf("expected 404 with NoIndex, got %d", rec.Code)
	}
}

func TestDirListing_DirsFirst(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/")
	body := rec.Body.String()
	subIdx := strings.Index(body, "sub/")
	helloIdx := strings.Index(body, "hello.txt")
	if subIdx > helloIdx {
		t.Errorf("expected directories listed before files")
	}
}

// --- Directory redirect ---

func TestDirRedirect(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/sub")
	if rec.Code != 301 {
		t.Errorf("expected 301 redirect, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc != "/sub/" {
		t.Errorf("expected redirect to /sub/, got %q", loc)
	}
}

// --- CORS ---

func TestCORS_Enabled(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir, CORS: true})

	rec := request(t, h, "GET", "/hello.txt")
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS header")
	}
}

func TestCORS_Options(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir, CORS: true})

	rec := request(t, h, "OPTIONS", "/hello.txt")
	if rec.Code != 204 {
		t.Errorf("expected 204 for OPTIONS, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS header on OPTIONS")
	}
}

func TestCORS_Disabled(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir, CORS: false})

	rec := request(t, h, "GET", "/hello.txt")
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("expected no CORS header when disabled")
	}
}

// --- Logging ---

func TestLogging_Enabled(t *testing.T) {
	dir := setupTestDir(t)
	var buf bytes.Buffer
	h := testHandler(t, Options{Root: dir, LogWriter: &buf})

	request(t, h, "GET", "/hello.txt")
	log := buf.String()
	if !strings.Contains(log, "GET") {
		t.Errorf("expected log to contain method, got %q", log)
	}
	if !strings.Contains(log, "/hello.txt") {
		t.Errorf("expected log to contain path, got %q", log)
	}
	if !strings.Contains(log, "200") {
		t.Errorf("expected log to contain status, got %q", log)
	}
}

func TestLogging_Disabled(t *testing.T) {
	dir := setupTestDir(t)
	var buf bytes.Buffer
	h := testHandler(t, Options{Root: dir, NoLog: true, LogWriter: &buf})

	request(t, h, "GET", "/hello.txt")
	if buf.Len() != 0 {
		t.Errorf("expected no log output with NoLog, got %q", buf.String())
	}
}

// --- Security ---

func TestPathTraversal_DotDot(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/../../../etc/passwd")
	if rec.Code != 404 {
		t.Errorf("expected 404 for path traversal, got %d", rec.Code)
	}
}

func TestPathTraversal_SiblingPrefix(t *testing.T) {
	// Root is e.g. /tmp/foo; sibling /tmp/foobar should not be accessible.
	parent := t.TempDir()
	root := filepath.Join(parent, "serve")
	sibling := filepath.Join(parent, "servebar")
	os.MkdirAll(root, 0o755)
	os.MkdirAll(sibling, 0o755)
	os.WriteFile(filepath.Join(sibling, "secret.txt"), []byte("secret"), 0o644)

	h := testHandler(t, Options{Root: root})

	// A symlink inside root that points to the sibling directory
	os.Symlink(sibling, filepath.Join(root, "escape"))
	rec := request(t, h, "GET", "/escape/secret.txt")
	if rec.Code != 404 {
		t.Errorf("expected 404 for symlink to sibling with shared prefix, got %d", rec.Code)
	}
}

func TestServeFile_SymlinkedRoot(t *testing.T) {
	// When the root directory itself is a symlink, files should still be served.
	dir := setupTestDir(t)
	parent := t.TempDir()
	link := filepath.Join(parent, "link")
	os.Symlink(dir, link)

	h := testHandler(t, Options{Root: link})

	rec := request(t, h, "GET", "/hello.txt")
	if rec.Code != 200 {
		t.Errorf("expected 200 when root is a symlink, got %d", rec.Code)
	}
	if rec.Body.String() != "hello world" {
		t.Errorf("expected 'hello world', got %q", rec.Body.String())
	}
}

func TestPathTraversal_Encoded(t *testing.T) {
	dir := setupTestDir(t)
	h := testHandler(t, Options{Root: dir})

	rec := request(t, h, "GET", "/%2e%2e/%2e%2e/etc/passwd")
	if rec.Code != 404 {
		t.Errorf("expected 404 for encoded path traversal, got %d", rec.Code)
	}
}

// --- formatSize ---

func TestFormatSize(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0B"},
		{1, "1B"},
		{999, "999B"},
		{1000, "1.0kB"},
		{1500, "1.5kB"},
		{1000000, "1.0MB"},
		{1500000, "1.5MB"},
		{1000000000, "1.0GB"},
	}
	for _, tt := range tests {
		got := formatSize(tt.n)
		if got != tt.want {
			t.Errorf("formatSize(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

// --- formatDuration ---

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms   int
		want string
	}{
		{0, "0ms"},
		{5, "5ms"},
		{999, "999ms"},
		{1500, "1.5s"},
	}
	for _, tt := range tests {
		got := formatDuration(time.Duration(tt.ms) * time.Millisecond)
		if got != tt.want {
			t.Errorf("formatDuration(%dms) = %q, want %q", tt.ms, got, tt.want)
		}
	}
}
