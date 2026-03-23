# swk — Developer's Swiss Army Knife

[![CI](https://github.com/agejevasv/swk/actions/workflows/ci.yml/badge.svg)](https://github.com/agejevasv/swk/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/agejevasv/swk/graph/badge.svg)](https://codecov.io/gh/agejevasv/swk)
[![Go Report Card](https://goreportcard.com/badge/github.com/agejevasv/swk)](https://goreportcard.com/report/github.com/agejevasv/swk)
[![Go Reference](https://pkg.go.dev/badge/github.com/agejevasv/swk.svg)](https://pkg.go.dev/github.com/agejevasv/swk)

A single-binary, zero-dependency CLI toolkit for everyday developer tasks.

## Install

### Quick install (macOS/Linux)

```bash
curl -sL https://raw.githubusercontent.com/agejevasv/swk/main/install.sh | sh
```

Or specify a directory:

```bash
SWK_INSTALL_DIR=~/.local/bin curl -sL https://raw.githubusercontent.com/agejevasv/swk/main/install.sh | sh
```

**Windows:** Download the binary from [GitHub Releases](https://github.com/agejevasv/swk/releases) and add it to your PATH. On WSL, the install script above works as-is.

### From source

```bash
git clone https://github.com/agejevasv/swk.git
cd swk
make build
```

### Shell completion

```bash
swk completion bash > /etc/bash_completion.d/swk
swk completion zsh > "${fpath[1]}/_swk"
swk completion fish > ~/.config/fish/completions/swk.fish
```
Source your shell or reopen.

## Usage

```
swk <category> <command> [file|input] [flags]
```

Commands read from **stdin** when no argument is given, write to **stdout**, and send errors to **stderr** (`serve`, `listen`, and `diff` are exceptions). Document-oriented commands (json, xml, markdown, hash, etc.) accept a **file path** as the argument — if the file exists, its contents are read automatically. Use `-` to explicitly read stdin.

```bash
cat data.json | swk format json
swk format json data.json
swk format json - < data.json
```

## Commands

| Category | Description |
|----------|-------------|
| [**convert**](#convert-swk-convert-swk-c) | Data format converters |
| [**format**](#format-swk-format-swk-fmt-swk-f) | Prettify, minify, and render data |
| [**encode**](#encode-swk-encode-swk-enc) | Encoders and decoders |
| [**escape**](#escape-swk-escape-swk-esc) | Escape and unescape strings |
| [**generate**](#generate-swk-generate-swk-g) | Data generators |
| [**inspect**](#inspect-swk-inspect-swk-i) | Inspect and analyze data |
| [**query**](#query-swk-query-swk-q) | Query and search data |
| [**diff**](#diff-swk-diff-swk-d) | Compare files |
| [**serve**](#serve-swk-serve) | Local static file server |
| [**listen**](#listen-swk-listen) | Log incoming HTTP requests |

---

### Convert (`swk convert`, `swk c`)

| Command | Description |
|---------|-------------|
| `convert base` | Convert between number bases (bin, oct, dec, hex) |
| `convert bytes` | Convert between byte sizes and human-readable formats |
| `convert case` | Convert between case conventions |
| `convert chmod` | Convert between numeric and symbolic file permissions |
| `convert color` | Convert between color formats |
| `convert csv2json` | Convert CSV to JSON |
| `convert date` | Convert between date/time formats |
| `convert duration` | Convert between seconds and human-readable durations |
| `convert image` / `img` | Convert image formats, resize |
| `convert json2csv` | Convert JSON array to CSV |
| `convert json2yaml` | Convert JSON to YAML |
| `convert markdown` / `md` | Render markdown to HTML or plain text |
| `convert yaml2json` | Convert YAML to JSON |

```bash
# Number base conversion
swk convert base 255 --from dec --to hex

# Byte sizes (default: 1024-based with IEC labels)
swk convert bytes 1073741824           # 1 GiB
swk convert bytes '1.5GiB'            # 1610612736
swk convert bytes -d 1000000000        # 1 GB (decimal, 1000-based)
swk convert bytes '1.5GB'             # 1500000000

# Case conversion
echo 'helloWorld' | swk convert case --to snake    # hello_world
echo 'hello world' | swk convert case --to camel   # helloWorld

# File permissions
swk convert chmod 755                  # shows rwxr-xr-x + breakdown
swk convert chmod rwxr-xr-x           # shows 755 + breakdown
swk convert chmod 4755                 # setuid support

# Color format conversion
swk convert color '#FF5733'
swk convert color '255,87,51' --from rgb --to hex

# Date/time conversion
swk convert date 1700000000 --from unix --to iso
swk convert date 1700000000 --from unix --to human --tz UTC
swk convert date now --to unix
swk convert date 1700000000 --from unix --to '%Y-%m-%d' --tz UTC   # 2023-11-14
swk convert date '2023-11-14' --from '%Y-%m-%d' --to unix

# Duration conversion
swk convert duration 86400             # 1d
swk convert duration '2d 5h 30m'      # 192600
swk convert duration 31536000          # 1y
swk convert duration '1y 6mo'         # 47088000

# Image conversion (accepts file path)
swk convert image photo.png --to jpeg -o photo.jpg
swk convert image large.png --to png --resize 200x200 -o thumb.png

# JSON to YAML
echo '{"name":"swk"}' | swk convert json2yaml
swk convert json2yaml data.json

# YAML to JSON
echo 'name: swk' | swk convert yaml2json
swk convert yaml2json config.yaml

# JSON to CSV
echo '[{"name":"alice","age":30}]' | swk convert json2csv

# CSV to JSON
printf 'name,age\nalice,30\n' | swk convert csv2json

# Render markdown (accepts file path)
swk convert markdown README.md --html > preview.html
swk convert markdown README.md --html --syntax-highlight > highlighted.html
```

### Format (`swk format`, `swk fmt`, `swk f`)

| Command | Description |
|---------|-------------|
| `format json` | Prettify or minify JSON |
| `format xml` | Prettify or minify XML |
| `format json2table` | Format JSON array as a table |
| `format csv2table` | Format CSV as a table |

```bash
# JSON prettify/minify (accepts file path)
swk format json data.json
swk format json config.json --minify
swk format json data.json --indent 4

# XML format (accepts file path)
swk format xml messy.xml
swk format xml document.xml --minify

# JSON array as table
echo '[{"name":"alice","age":30},{"name":"bob","age":25}]' | swk format json2table
echo '[{"name":"alice"}]' | swk format json2table --style simple

# CSV as table
printf 'name,age\nalice,30\n' | swk format csv2table

# Extract nested array, then render as table
echo '{"meta":"v1","users":[{"name":"alice"},{"name":"bob"}]}' \
  | swk query json '$.users' \
  | swk format json2table
```

### Encode (`swk encode`, `swk enc`)

| Command | Description |
|---------|-------------|
| `encode base64` / `b64` | Base64 encode/decode |
| `encode hash` / `sum` | Generate hashes (MD5, SHA1, SHA256, SHA512) |
| `encode html` | HTML entity encode/decode (alias of `escape html`) |
| `encode jwt` | Create, decode, or verify JWT tokens |
| `encode qr` | Generate QR codes |
| `encode url` | URL percent-encode/decode (alias of `escape url`) |

```bash
# Base64
echo 'hello world' | swk encode base64
echo 'aGVsbG8gd29ybGQ=' | swk encode base64 -d
echo 'data' | swk encode base64 --url-safe

# Hash (accepts file path)
swk encode hash README.md
swk encode hash README.md --algo md5
echo -n 'hello' | swk encode hash --verify 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824

# JWT — create with HMAC
swk encode jwt --secret mykey '{"sub":"user1","role":"admin"}'

# JWT — create with RSA/EC/Ed25519
swk encode jwt --algo RS256 --key private.pem '{"sub":"user1"}'
swk encode jwt --algo ES256 --key ec-private.pem '{"sub":"user1"}'

# JWT — decode (no verification, works with any algorithm)
swk encode jwt -d 'eyJhbGciOiJIUzI1NiIs...'

# JWT — verify with HMAC secret
swk encode jwt -d --secret mykey 'eyJhbGciOiJIUzI1NiIs...'

# JWT — verify with public key (RSA/EC/Ed25519)
swk encode jwt -d --key public.pem 'eyJhbGciOiJSUzI1NiIs...'

# QR code
swk encode qr 'https://github.com/agejevasv/swk'
swk encode qr --output png 'https://example.com' > qr.png
```

### Escape (`swk escape`, `swk esc`)

| Command | Description |
|---------|-------------|
| `escape html` | HTML entity escape/unescape |
| `escape json` | JSON string escape/unescape |
| `escape shell` / `sh` | Shell escape/unescape |
| `escape url` | URL percent-encode/decode |
| `escape xml` | XML escape/unescape |

```bash
# HTML
echo '<script>alert("xss")</script>' | swk escape html
echo '&lt;div&gt;' | swk escape html -u

# JSON
echo 'line1\nline2' | swk escape json
echo '\"hello\"' | swk escape json -u

# Shell
echo "it's a test" | swk escape shell

# URL
echo 'hello world & friends' | swk escape url
echo 'hello%20world' | swk escape url -u

# XML
echo '<tag attr="val">' | swk escape xml
```

### Generate (`swk generate`, `swk g`)

| Command | Description |
|---------|-------------|
| `generate cert` | Generate self-signed TLS certificates |
| `generate cron` | Generate cron expressions from flags |
| `generate image` | Generate placeholder images |
| `generate password` / `pw` | Generate random passwords |
| `generate text` | Generate lorem ipsum text |
| `generate uuid` | Generate UUIDs (v1, v4, v5, v7) |

```bash
# Self-signed TLS certificate
swk generate cert
swk generate cert --cn myapp.local --days 30
swk generate cert --dns localhost --dns myapp.local --ip 127.0.0.1
swk generate cert --key-type rsa -o ./certs/server

# Cron expressions
swk generate cron --every 5m                     # */5 * * * *
swk generate cron --daily --at 9:00              # 0 9 * * *
swk generate cron --weekdays --at 9:00           # 0 9 * * 1-5
swk generate cron --weekly --day MON --at 9:00   # 0 9 * * 1
swk generate cron --monthly --day 15             # 0 0 15 * *
swk generate cron --yearly --month JUN --day 1   # 0 0 1 6 *

# Generate and verify
swk generate cron --daily --at 9:00 | swk inspect cron

# Placeholder image
swk generate image -o placeholder.png
swk generate image --style circles --width 1920 --height 1080 -o wallpaper.png

# Passwords
swk generate password
swk generate password --length 32 --no-symbols
swk generate password --count 10

# Lorem ipsum
swk generate text --paragraphs 3
swk generate text --words 50

# UUIDs
swk generate uuid
swk generate uuid --count 5
swk generate uuid --version 7
```

### Inspect (`swk inspect`, `swk i`)

| Command | Description |
|---------|-------------|
| `inspect cert` | Inspect X.509 PEM certificates |
| `inspect cron` | Explain cron expressions |
| `inspect dns` | DNS lookups (A, AAAA, MX, NS, TXT, CNAME) |
| `inspect domain` | Domain registration info via RDAP |
| `inspect ip` | Show your public IP address |
| `inspect jwt` | Inspect JWT token claims and expiry |
| `inspect net` | List processes listening on network ports (Linux/macOS) |
| `inspect subnet` | Calculate subnet information from CIDR |
| `inspect text` / `txt` | Character, word, line, byte counts |
| `inspect url` | Parse URL into components |

```bash
# Certificate (accepts file path)
swk inspect cert cert.pem
swk inspect cert cert.pem --check-expiry

# Cron
swk inspect cron '*/5 * * * *'
swk inspect cron --explain '0 9 * * 1-5'
swk inspect cron --next 3 '0 9 * * MON'

# DNS lookups
swk inspect dns example.com
swk inspect dns example.com --type MX
swk inspect dns 8.8.8.8                # reverse lookup
swk inspect dns --json example.com

# Domain registration info
swk inspect domain example.com
swk inspect domain --json example.com

# Public IP
swk inspect ip

# JWT token inspection
swk inspect jwt 'eyJhbGciOiJIUzI1NiIs...'
echo 'eyJhbGciOiJIUzI1NiIs...' | swk inspect jwt
swk inspect jwt --check-expiry 'eyJhbGciOiJIUzI1NiIs...'
swk inspect jwt --json 'eyJhbGciOiJIUzI1NiIs...'

# Network ports (Linux/macOS)
swk inspect net
swk inspect net --all                   # include established connections
swk inspect net --tcp --port 8080       # filter by protocol and port
swk inspect net --json

# Subnet calculator
swk inspect subnet 192.168.1.0/24
swk inspect subnet 10.0.0.0/16

# Text stats (accepts file path)
swk inspect text essay.txt
echo 'hello world' | swk inspect text --json

# URL parsing
swk inspect url 'https://example.com:8080/api/v1/users?page=1&limit=10#section'
```

### Query (`swk query`, `swk q`)

| Command | Description |
|---------|-------------|
| `query html` | Query HTML with CSS selectors |
| `query json` | Query JSON with JSONPath |
| `query regex` / `re` | Match/replace with regular expressions |

```bash
# HTML (CSS selectors — accepts file path)
curl -s http://example.com | swk query html 'a' --attr href
swk query html 'div.content p' page.html

# JSON (JSONPath — accepts file path)
swk query json '$.users[*].name' data.json

# Regex (pattern is first argument — accepts file path)
echo '2024-01-15 hello 2024-02-20' | swk query regex -o -g '\d{4}-\d{2}-\d{2}'
echo 'John:30' | swk query regex --groups '(\w+):(\d+)'
echo 'foo bar baz' | swk query regex -r 'qux' 'bar'
```

### Diff (`swk diff`, `swk d`)

| Command | Description |
|---------|-------------|
| `diff text` / `txt` | Unified text diff |
| `diff json` | Semantic JSON diff (normalizes key order) |

```bash
# Text diff
swk diff text old.txt new.txt

# JSON diff (ignores key order)
swk diff json old.json new.json
swk diff json <(curl -s api/v1) <(curl -s api/v2)
curl -s api/v1 | swk diff json - saved.json
```

### Serve (`swk serve`)

| Command | Description |
|---------|-------------|
| `serve` | Start a local static file server |

```bash
# Serve current directory on port 8080
swk serve

# Serve a specific directory on a custom port
swk serve ./dist --port 3000

# Enable CORS headers for local API development
swk serve --cors

# Disable directory listing
swk serve --no-index

# Bind to localhost only, random port
swk serve --host 127.0.0.1 --port 0

# HTTPS with self-signed cert
swk generate cert
swk serve --tls
```

### Listen (`swk listen`)

| Command | Description |
|---------|-------------|
| `listen` | Log incoming HTTP requests (webhook receiver) |

```bash
# Log all incoming requests on port 8080
swk listen

# Custom port and response
swk listen --port 9000 --status 201 --body '{"ok":true}'

# Don't log request bodies
swk listen --no-body
```

## Piping

Commands read stdin and write stdout, making them composable:

```bash
# JSON → format → base64
echo '{"a":1}' | swk format json | swk encode base64

# Generate password → tee to stderr → hash
swk generate password | tee /dev/stderr | swk encode hash

# YAML → JSON → minify → clipboard (macOS)
cat config.yaml | swk convert yaml2json | swk format json -m | pbcopy

# Fetch API → format as table
curl -s https://api.example.com/users | swk format json2table

# Extract nested array → table
echo '{"status":"ok","items":[{"id":1,"name":"foo"},{"id":2,"name":"bar"}]}' \
  | swk query json '$.items' \
  | swk format json2table

# Duration roundtrip
swk convert duration '2d 5h' --to seconds | swk convert duration --to human

# Generate cron → verify
swk generate cron --daily --at 9:00 | swk inspect cron

# Scrape links from a webpage
curl -s http://example.com | swk query html 'a' --attr href
```

## Testing

```bash
make test          # run all tests
make test-verbose  # run with verbose output
make lint          # run go vet + staticcheck
```

## License

MIT
