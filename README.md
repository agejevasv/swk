# swk — Developer's Swiss Army Knife

A command-line toolkit for everyday developer tasks. Encode, decode, format, convert, generate, and inspect data — all from your terminal.

## Install

### Quick install

```bash
curl -sL https://raw.githubusercontent.com/agejevasv/swk/main/install.sh | sh
```

Or specify a directory:

```bash
SWK_INSTALL_DIR=~/.local/bin curl -sL https://raw.githubusercontent.com/agejevasv/swk/main/install.sh | sh
```

### From source

Requires Go 1.26+.

```bash
git clone https://github.com/agejevasv/swk.git
cd swk
make build
```

This produces the `swk` binary in the current directory.

To install into your `$GOPATH/bin`:

```bash
make install
```

## Usage

```
swk <category> <command> [file|input] [flags]
```

Every command reads from **stdin** when no argument is given, writes to **stdout**, and sends errors to **stderr**. Document-oriented commands (json, xml, markdown, hash, etc.) accept a **file path** as the argument — if the file exists, its contents are read automatically. Use `-` to explicitly read stdin.

```bash
# These are equivalent:
cat data.json | swk convert json
swk convert json data.json
swk convert json - < data.json
```

Print version:

```bash
swk --version
swk -V
```

## Commands

### Convert (`swk convert`, `swk c`)

| Command | Alias | Description |
|---------|-------|-------------|
| `convert base` | | Convert between number bases (bin, oct, dec, hex) |
| `convert bytes` | | Convert between byte sizes and human-readable formats |
| `convert case` | | Convert between case conventions |
| `convert chmod` | | Convert between numeric and symbolic file permissions |
| `convert color` | | Convert between color formats |
| `convert date` | | Convert between date/time formats |
| `convert duration` | | Convert between seconds and human-readable durations |
| `convert image` | `c img` | Convert image formats, resize |
| `convert json` | | Convert and format JSON (yaml, csv) |
| `convert markdown` | `c md` | Render markdown to HTML or plain text |
| `convert table` | | Render a JSON array or CSV as a formatted table |
| `convert xml` | | Format (prettify/minify) XML |

```bash
# Number base conversion
swk convert base --from dec --to hex 255

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
swk convert color --from rgb --to hex '255,87,51'

# Date/time conversion
swk convert date --from unix --to iso 1700000000
swk convert date now
swk convert date --from unix --to '%Y-%m-%d' --tz UTC 1700000000   # 2023-11-14
swk convert date --from '%Y-%m-%d' --to unix '2023-11-14'

# Duration conversion
swk convert duration 86400             # 1d
swk convert duration '2d 5h 30m'      # 192600
swk convert duration 31536000          # 1y
swk convert duration '1y 6mo'         # 47088000

# Image conversion (accepts file path)
swk convert image --to jpeg photo.png -o photo.jpg
swk convert image --to png --resize 200x200 large.png -o thumb.png

# JSON prettify/minify (accepts file path)
swk convert json data.json
swk convert json --minify config.json

# JSON to YAML
echo '{"name":"swk"}' | swk convert json --to yaml

# YAML to JSON
echo 'name: swk' | swk convert json --from yaml

# JSON to CSV
echo '[{"name":"alice","age":30}]' | swk convert json --to csv

# CSV to JSON
echo 'name,age\nalice,30' | swk convert json --from csv

# Render markdown (accepts file path)
swk convert markdown --html README.md > preview.html
swk convert markdown --html --syntax-highlight README.md > highlighted.html

# JSON array as table
echo '[{"name":"alice","age":30},{"name":"bob","age":25}]' | swk convert table
echo '[{"name":"alice"}]' | swk convert table --style simple
printf 'name,age\nalice,30\n' | swk convert table --from csv

# Extract nested array, then render as table
echo '{"meta":"v1","users":[{"name":"alice"},{"name":"bob"}]}' \
  | swk query json -q '$.users' \
  | swk convert table

# XML format (accepts file path)
swk convert xml messy.xml
swk convert xml --minify document.xml
```

### Encode (`swk encode`, `swk enc`)

| Command | Alias | Description |
|---------|-------|-------------|
| `encode base64` | `enc b64` | Base64 encode/decode |
| `encode hash` | `enc sum` | Generate hashes (MD5, SHA1, SHA256, SHA512) |
| `encode jwt` | | Create, decode, or verify JWT tokens |
| `encode qr` | | Generate QR codes |

```bash
# Base64
echo 'hello world' | swk encode base64
echo 'aGVsbG8gd29ybGQ=' | swk encode base64 -d
echo 'data' | swk encode base64 --url-safe

# Hash (accepts file path)
swk encode hash README.md
swk encode hash --algo md5 README.md
echo -n 'hello' | swk encode hash --verify 2cf24dba...

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

| Command | Alias | Description |
|---------|-------|-------------|
| `escape html` | | HTML entity escape/unescape |
| `escape json` | | JSON string escape/unescape |
| `escape shell` | `esc sh` | Shell escape/unescape |
| `escape url` | | URL percent-encode/decode |
| `escape xml` | | XML escape/unescape |

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

| Command | Alias | Description |
|---------|-------|-------------|
| `generate image` | | Generate placeholder images |
| `generate password` | `g pw` | Generate random passwords |
| `generate text` | | Generate lorem ipsum text |
| `generate uuid` | | Generate UUIDs (v1, v4, v5, v7) |

```bash
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

| Command | Alias | Description |
|---------|-------|-------------|
| `inspect cert` | | Inspect X.509 PEM certificates |
| `inspect cron` | | Explain cron expressions |
| `inspect text` | `i txt` | Character, word, line, byte counts |
| `inspect url` | | Parse URL into components |

```bash
# Certificate (accepts file path)
swk inspect cert cert.pem
swk inspect cert --check-expiry cert.pem

# Cron
swk inspect cron '*/5 * * * *'
swk inspect cron --explain '0 9 * * 1-5'
swk inspect cron --next 3 '0 9 * * MON'

# Text stats (accepts file path)
swk inspect text essay.txt
echo 'hello world' | swk inspect text --json

# URL parsing
swk inspect url 'https://example.com:8080/api/v1/users?page=1&limit=10#section'
```

### Query (`swk query`, `swk q`)

| Command | Alias | Description |
|---------|-------|-------------|
| `query html` | | Query HTML with CSS selectors |
| `query json` | | Query JSON with JSONPath |
| `query regex` | `q re` | Match/replace with regular expressions |

```bash
# HTML (CSS selectors — accepts file path)
curl -s https://example.com | swk query html -q 'a' --attr href
swk query html -q 'div.content p' page.html

# JSON (JSONPath — accepts file path)
swk query json -q '$.users[*].name' data.json

# Regex (pattern is first argument — accepts file path)
echo '2024-01-15 hello 2024-02-20' | swk query regex -g '\d{4}-\d{2}-\d{2}'
echo 'John:30' | swk query regex --groups '(\w+):(\d+)'
echo 'foo bar baz' | swk query regex -r 'qux' 'bar'
```

## Piping and chaining

Every command reads stdin and writes stdout:

```bash
# JSON → format → base64
echo '{"a":1}' | swk convert json | swk encode base64

# Generate password → hash it
swk generate password | swk encode hash

# YAML → JSON → minify → clipboard (macOS)
cat config.yaml | swk convert json --from yaml | swk convert json -m | pbcopy

# Fetch API → format → as table
curl -s https://api.example.com/users | swk convert table

# Extract nested array → table
echo '{"status":"ok","items":[{"id":1,"name":"foo"},{"id":2,"name":"bar"}]}' \
  | swk query json -q '$.items' \
  | swk convert table

# Duration roundtrip
swk convert duration '2d 5h' --to seconds | swk convert duration --to human

# CSV → JSON → query with JSONPath
swk convert json --from csv data.csv | swk query json -q '$..[?(@.age>30)]'

# Scrape links from a webpage
curl -s https://example.com | swk query html -q 'a' --attr href
```

## Shell completion

```bash
# Bash
swk completion bash > /etc/bash_completion.d/swk

# Zsh
swk completion zsh > "${fpath[1]}/_swk"

# Fish
swk completion fish > ~/.config/fish/completions/swk.fish
```

## Testing

```bash
make test          # run all tests
make test-verbose  # run with verbose output
make lint          # run go vet + staticcheck
```

## License

MIT
