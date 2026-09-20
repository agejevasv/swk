package listen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxBodySize = 1 << 20 // 1MB

// Options configures the listen handler.
type Options struct {
	Status int
	Body   string
	NoBody bool
	Writer io.Writer
}

// Handler returns an http.Handler that logs all incoming requests.
func Handler(opts Options) http.Handler {
	if opts.Status == 0 {
		opts.Status = 200
	}

	var mu sync.Mutex

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		truncated := false
		if r.Body != nil {
			// Read one byte past the cap so truncation can be reported.
			body, _ = io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
			r.Body.Close()
			if len(body) > maxBodySize {
				body, truncated = body[:maxBodySize], true
			}
		}

		mu.Lock()
		logRequest(opts.Writer, r, body, opts.NoBody, truncated)
		mu.Unlock()

		w.WriteHeader(opts.Status)
		if opts.Body != "" {
			fmt.Fprint(w, opts.Body)
		}
	})
}

func logRequest(w io.Writer, r *http.Request, body []byte, noBody, truncated bool) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	// RequestURI() cannot represent an authority-form target such as CONNECT.
	target := r.RequestURI
	if target == "" {
		target = r.URL.RequestURI()
	}
	fmt.Fprintf(w, "\n--- %s %s %s ---\n", r.Method, sanitize(target), ts)

	// Headers, sorted
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range r.Header[k] {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	}

	if r.Host != "" {
		// Host is not in Header map, print separately if not already there
		if r.Header.Get("Host") == "" {
			fmt.Fprintf(w, "Host: %s\n", r.Host)
		}
	}

	if noBody || len(body) == 0 {
		return
	}

	fmt.Fprintln(w)

	// Pretty-print JSON bodies
	if isJSON(r.Header.Get("Content-Type")) {
		var buf bytes.Buffer
		if json.Indent(&buf, body, "", "  ") == nil {
			fmt.Fprintln(w, buf.String())
			return
		}
	}

	fmt.Fprintln(w, sanitize(string(body)))
	if truncated {
		fmt.Fprintf(w, "[body truncated at %d MiB]\n", maxBodySize>>20)
	}
}

// sanitize strips control characters so a request cannot inject ANSI escapes
// or forge extra log records on the operator's terminal.
func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return r
		}
		if r < 0x20 || r == 0x7f {
			return '\ufffd'
		}
		return r
	}, s)
}

func isJSON(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}
