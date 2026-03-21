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
	root := filepath.Clean(opts.Root)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cleaned := path.Clean("/" + r.URL.Path)
		fp := filepath.Join(root, filepath.FromSlash(cleaned))

		// Resolve symlinks and verify the path is under root
		resolved, err := filepath.EvalSymlinks(fp)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !strings.HasPrefix(resolved, root) {
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(resolved)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if !info.IsDir() {
			f, err := os.Open(resolved)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer f.Close()
			http.ServeContent(w, r, info.Name(), info.ModTime(), f)
			return
		}

		// Directory: redirect if no trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
			return
		}

		// Try index files
		for _, index := range []string{"index.html", "index.htm"} {
			indexPath := filepath.Join(resolved, index)
			if fi, err := os.Stat(indexPath); err == nil && !fi.IsDir() {
				f, err := os.Open(indexPath)
				if err != nil {
					http.NotFound(w, r)
					return
				}
				defer f.Close()
				http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
				return
			}
		}

		// Directory listing
		if opts.NoIndex {
			http.NotFound(w, r)
			return
		}

		renderDirListing(w, cleaned, resolved)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	n, err := sw.ResponseWriter.Write(b)
	sw.size += n
	return n, err
}

func loggingMiddleware(next http.Handler, logW io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		fmt.Fprintf(logW, "%s %s %d %s %s\n",
			r.Method, r.URL.Path, sw.status, formatSize(sw.size), formatDuration(time.Since(start)))
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

func renderDirListing(w http.ResponseWriter, reqPath string, dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return entries[i].Name() < entries[j].Name()
	})

	data := dirListingData{
		Path:      reqPath,
		HasParent: reqPath != "/",
	}

	for _, entry := range entries {
		name := entry.Name()
		href := url.PathEscape(name)
		size := "-"
		modTime := "-"

		if info, err := entry.Info(); err == nil {
			modTime = info.ModTime().Format("2006-01-02 15:04")
			if !entry.IsDir() {
				size = formatSize(int(info.Size()))
			}
		}

		if entry.IsDir() {
			name += "/"
			href += "/"
		}

		data.Entries = append(data.Entries, dirEntryData{
			Name:    name,
			Href:    href,
			Size:    size,
			ModTime: modTime,
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dirListingTmpl.Execute(w, data)
}

func formatSize(n int) string {
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
