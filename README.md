# swk — Developer's Swiss Army Knife

A CLI toolkit for everyday developer tasks.

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

```bash
git clone https://github.com/agejevasv/swk.git
cd swk
make build
```

### Shell completion

```bash
# Bash
swk completion bash > /etc/bash_completion.d/swk

# Zsh
swk completion zsh > "${fpath[1]}/_swk"

# Fish
swk completion fish > ~/.config/fish/completions/swk.fish
```

## Usage

```
swk <category> <command> [file|input] [flags]
```

Every command reads from **stdin** when no argument is given, writes to **stdout**, and sends errors to **stderr**. Document-oriented commands (json, xml, markdown, hash, etc.) accept a **file path** as the argument — if the file exists, its contents are read automatically. Use `-` to explicitly read stdin.

```bash
# These are equivalent:
cat data.json | swk format json
swk format json data.json
swk format json - < data.json
```

## Commands

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
swk convert date --from unix --to '2006-01-02' --tz UTC 1700000000   # 2023-11-14
swk convert date --from '2006-01-02' --to unix '2023-11-14'

# Duration conversion
swk convert duration 86400             # 1d
swk convert duration '2d 5h 30m'      # 192600
swk convert duration 31536000          # 1y
swk convert duration '1y 6mo'         # 47088000

# Image conversion (accepts file path)
swk convert image --to jpeg photo.png -o photo.jpg
swk convert image --to png --resize 200x200 large.png -o thumb.png

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
swk convert markdown --html README.md > preview.html
swk convert markdown --html --syntax-highlight README.md > highlighted.html
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
swk format json --minify config.json
swk format json --indent 4 data.json

# XML format (accepts file path)
swk format xml messy.xml
swk format xml --minify document.xml

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
| `generate image` | Generate placeholder images |
| `generate password` / `pw` | Generate random passwords |
| `generate text` | Generate lorem ipsum text |
| `generate uuid` | Generate UUIDs (v1, v4, v5, v7) |

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

| Command | Description |
|---------|-------------|
| `inspect cert` | Inspect X.509 PEM certificates |
| `inspect cron` | Explain cron expressions |
| `inspect text` / `txt` | Character, word, line, byte counts |
| `inspect url` | Parse URL into components |

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

| Command | Description |
|---------|-------------|
| `query html` | Query HTML with CSS selectors |
| `query json` | Query JSON with JSONPath |
| `query regex` / `re` | Match/replace with regular expressions |

```bash
# HTML (CSS selectors — accepts file path)
curl -s https://example.com | swk query html 'a' --attr href
swk query html 'div.content p' page.html

# JSON (JSONPath — accepts file path)
swk query json '$.users[*].name' data.json

# Regex (pattern is first argument — accepts file path)
echo '2024-01-15 hello 2024-02-20' | swk query regex -o -g '\d{4}-\d{2}-\d{2}'
echo 'John:30' | swk query regex --groups '(\w+):(\d+)'
echo 'foo bar baz' | swk query regex -r 'qux' 'bar'
```

## Piping and chaining

Every command reads stdin and writes stdout:

```bash
# JSON → format → base64
echo '{"a":1}' | swk format json | swk encode base64

# Generate password → hash it
swk generate password | swk encode hash

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

# CSV → JSON → query with JSONPath
swk convert csv2json data.csv | swk query json '$..[?(@.age>30)]'

# Scrape links from a webpage
curl -s https://example.com | swk query html 'a' --attr href
```

## Testing

```bash
make test          # run all tests
make test-verbose  # run with verbose output
make lint          # run go vet + staticcheck
```

## License

MIT
