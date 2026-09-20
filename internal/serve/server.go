package serve

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Options configures the file server behavior.
type Options struct {
	Root      string
	Host      string
	Port      int
	CORS      bool
	NoIndex   bool
	NoLog     bool
	LogWriter io.Writer
}

// Handler returns an http.Handler configured per opts.
func Handler(opts Options) http.Handler {
	var h http.Handler = fileHandler(opts)
	if opts.CORS {
		h = corsMiddleware(h)
	}
	if !opts.NoLog {
		h = loggingMiddleware(h, opts.LogWriter)
	}
	return h
}

func fileHandler(opts Options) http.HandlerFunc {
	// os.Root confines every lookup to the directory: symlinks that leave it
	// cannot be followed, and the check cannot be raced by swapping a path
	// component between the check and the open.
	root, rootErr := os.OpenRoot(filepath.Clean(opts.Root))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if rootErr != nil {
			http.Error(w, "Server misconfigured", http.StatusInternalServerError)
			return
		}

		cleaned := path.Clean("/" + r.URL.Path)
		name := strings.TrimPrefix(cleaned, "/")
		if name == "" {
			name = "."
		}

		f, info, ok := openConfined(root, name)
		if !ok {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		if !info.IsDir() {
			http.ServeContent(w, r, info.Name(), info.ModTime(), f)
			return
		}

		// Directory: redirect if no trailing slash. Build the target from the
		// escaped path so reserved characters survive.
		if !strings.HasSuffix(r.URL.Path, "/") {
			target := (&url.URL{Path: cleaned + "/", RawQuery: r.URL.RawQuery}).String()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}

		for _, index := range []string{"index.html", "index.htm"} {
			idx, idxInfo, ok := openConfined(root, path.Join(name, index))
			if !ok {
				continue
			}
			defer idx.Close()
			if idxInfo.IsDir() {
				continue
			}
			http.ServeContent(w, r, idxInfo.Name(), idxInfo.ModTime(), idx)
			return
		}

		if opts.NoIndex {
			http.NotFound(w, r)
			return
		}

		renderDirListing(w, cleaned, root, f)
	}
}

// openConfined opens a path inside the root and reports whether it is servable.
// Anything that is not a regular file or a directory (FIFOs, devices, sockets)
// is refused: opening one can block the handler forever.
func openConfined(root *os.Root, name string) (*os.File, os.FileInfo, bool) {
	// Check the mode before opening: opening a FIFO blocks until a writer
	// appears, which would pin the handler goroutine indefinitely.
	info, err := root.Stat(name)
	if err != nil || !(info.Mode().IsRegular() || info.IsDir()) {
		return nil, nil, false
	}

	f, err := root.Open(name)
	if err != nil {
		return nil, nil, false
	}

	return f, info, true
}

type statusWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	n, err := sw.ResponseWriter.Write(b)
	sw.size += int64(n)
	return n, err
}

func loggingMiddleware(next http.Handler, logW io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		fmt.Fprintf(logW, "%s %s %d %s %s\n",
			r.Method, sanitizeForLog(r.URL.Path), sw.status, formatSize(sw.size), formatDuration(time.Since(start)))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Directory listing

type dirListingData struct {
	Path      string
	HasParent bool
	Entries   []dirEntryData
}

type dirEntryData struct {
	Name    string
	Href    string
	Size    string
	ModTime string
	IsDir   bool
}

var dirListingTmpl = template.Must(template.New("dirlist").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Index of {{.Path}}</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; max-width: 900px; margin: 40px auto; padding: 0 20px; color: #222; }
h1 { font-size: 1.4em; border-bottom: 1px solid #ddd; padding-bottom: 0.3em; }
table { border-collapse: collapse; width: 100%; }
th, td { text-align: left; padding: 6px 12px; border-bottom: 1px solid #eee; }
th { font-weight: 600; }
td.size, th.size { text-align: right; }
a { color: #0366d6; text-decoration: none; }
a:hover { text-decoration: underline; }
</style>
</head>
<body>
<h1>Index of {{.Path}}</h1>
<table>
<thead><tr><th>Name</th><th class="size">Size</th><th>Modified</th></tr></thead>
<tbody>
{{- if .HasParent}}
<tr><td><a href="../">../</a></td><td class="size">-</td><td>-</td></tr>
{{- end}}
{{- range .Entries}}
<tr><td><a href="{{.Href}}">{{.Name}}</a></td><td class="size">{{.Size}}</td><td>{{.ModTime}}</td></tr>
{{- end}}
</tbody>
</table>
</body>
</html>
`))

func renderDirListing(w http.ResponseWriter, reqPath string, root *os.Root, dir *os.File) {
	entries, err := dir.ReadDir(-1)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	relDir := strings.TrimPrefix(reqPath, "/")

	data := dirListingData{
		Path:      reqPath,
		HasParent: reqPath != "/",
	}

	for _, entry := range entries {
		// Stat through the root: it follows symlinks that stay inside and
		// fails for those that leave, so the listing shows what can be served
		// and reveals nothing about targets outside the root.
		info, err := root.Stat(path.Join(relDir, entry.Name()))
		if err != nil {
			continue
		}

		name := entry.Name()
		href := "./" + url.PathEscape(name)
		size := "-"
		if !info.IsDir() {
			size = formatSize(info.Size())
		} else {
			name += "/"
			href += "/"
		}

		data.Entries = append(data.Entries, dirEntryData{
			Name:    name,
			Href:    href,
			Size:    size,
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
			IsDir:   info.IsDir(),
		})
	}

	// Directories first, then by name.
	sort.Slice(data.Entries, func(i, j int) bool {
		if data.Entries[i].IsDir != data.Entries[j].IsDir {
			return data.Entries[i].IsDir
		}
		return data.Entries[i].Name < data.Entries[j].Name
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Nothing can be done about a write failure here: the header is already out.
	_ = dirListingTmpl.Execute(w, data)
}

// sanitizeForLog strips control characters so a request path cannot inject
// ANSI escapes or forge extra lines in the operator's terminal.
func sanitizeForLog(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return '\ufffd'
		}
		return r
	}, s)
}

func formatSize(n int64) string {
	switch {
	case n == 0:
		return "0B"
	case n < 1000:
		return fmt.Sprintf("%dB", n)
	case n < 1000000:
		return fmt.Sprintf("%.1fkB", float64(n)/1000)
	case n < 1000000000:
		return fmt.Sprintf("%.1fMB", float64(n)/1000000)
	default:
		return fmt.Sprintf("%.1fGB", float64(n)/1000000000)
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
