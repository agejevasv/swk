package listen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(io.LimitReader(r.Body, maxBodySize))
			r.Body.Close()
		}

		logRequest(opts.Writer, r, body, opts.NoBody)

		w.WriteHeader(opts.Status)
		if opts.Body != "" {
			fmt.Fprint(w, opts.Body)
		}
	})
}

func logRequest(w io.Writer, r *http.Request, body []byte, noBody bool) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "\n--- %s %s %s ---\n", r.Method, r.URL.RequestURI(), ts)

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

	fmt.Fprintln(w, string(body))
}

func isJSON(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}
