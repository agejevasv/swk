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

Requires Go 1.25+.

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
swk <category> <command> [input] [flags]
```

Every command reads from **stdin** when no argument is given, writes to **stdout**, and sends errors to **stderr**.

Print version:

```bash
swk --version
swk -v
```

## Commands

### Convert (`swk convert`)

| Command | Alias | Description |
|---------|-------|-------------|
| `convert base` | `c nb` | Convert between number bases (bin, oct, dec, hex) |
| `convert case` | `c case` | Convert between case conventions |
| `convert color` | `c col` | Convert between color formats |
| `convert date` | `c dt` | Convert between date/time formats |
| `convert image` | `c img` | Convert image formats, resize |
| `convert json` | `c j` | Convert and format JSON (yaml, csv) |
| `convert markdown` | `c md` | Render markdown to HTML or plain text |
| `convert xml` | `c x` | Format (prettify/minify) XML |

```bash
# Number base conversion
swk convert base --from dec --to hex 255

# Case conversion
echo 'helloWorld' | swk convert case --to snake    # hello_world

# Color format conversion
swk convert color '#FF5733'
swk convert color --from rgb --to hex '255,87,51'

# Date/time conversion
swk convert date --from unix --to iso 1700000000
swk convert date now

# Image conversion
swk convert image --to jpeg -i photo.png -o photo.jpg
swk convert image --to png --resize 200x200 -i large.png -o thumb.png

# JSON prettify/minify
echo '{"a":1,"b":2}' | swk convert json
cat config.json | swk convert json --minify

# JSON to YAML
echo '{"name":"swk"}' | swk convert json --to yaml

# YAML to JSON
echo 'name: swk' | swk convert json --from yaml

# JSON to CSV
echo '[{"name":"alice","age":30}]' | swk convert json --to csv

# CSV to JSON
echo 'name,age\nalice,30' | swk convert json --from csv

# Render markdown
cat README.md | swk convert markdown --html > preview.html

# XML format
cat messy.xml | swk convert xml
cat document.xml | swk convert xml --minify
```

### Encode (`swk encode`)

| Command | Alias | Description |
|---------|-------|-------------|
| `encode base64` | `enc b64` | Base64 encode/decode |
| `encode hash` | `enc h` | Generate hashes (MD5, SHA1, SHA256, SHA512) |
| `encode jwt` | `enc jwt` | Create or decode JWT tokens |
| `encode qr` | `enc qr` | Generate QR codes |

```bash
# Base64
echo 'hello world' | swk encode base64
echo 'aGVsbG8gd29ybGQ=' | swk encode base64 -d
echo 'data' | swk encode base64 --url-safe

# Hash
echo -n 'hello' | swk encode hash
echo -n 'hello' | swk encode hash --algo md5
echo -n 'hello' | swk encode hash --verify 2cf24dba...

# JWT
swk encode jwt --secret mykey '{"sub":"user1","role":"admin"}'
swk encode jwt -d 'eyJhbGciOiJIUzI1NiIs...'
swk encode jwt -d --secret mykey 'eyJhbGciOiJIUzI1NiIs...'

# QR code
swk encode qr 'https://github.com/agejevasv/swk'
swk encode qr --output png 'https://example.com' > qr.png
```

### Escape (`swk escape`)

| Command | Alias | Description |
|---------|-------|-------------|
| `escape html` | `esc html` | HTML entity escape/unescape |
| `escape json` | `esc json` | JSON string escape/unescape |
| `escape shell` | `esc shell` | Shell escape/unescape |
| `escape url` | `esc url` | URL percent-encode/decode |
| `escape xml` | `esc xml` | XML escape/unescape |

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

### Generate (`swk generate`)

| Command | Alias | Description |
|---------|-------|-------------|
| `generate image` | `gen image` | Generate placeholder images |
| `generate password` | `gen pw` | Generate random passwords |
| `generate text` | `gen text` | Generate lorem ipsum text |
| `generate uuid` | `gen uid` | Generate UUIDs (v1, v4, v5, v7) |

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

### Inspect (`swk inspect`)

| Command | Alias | Description |
|---------|-------|-------------|
| `inspect cert` | | Inspect X.509 PEM certificates |
| `inspect cron` | `inspect cr` | Explain cron expressions |
| `inspect text` | `inspect txt` | Character, word, line, byte counts |

```bash
# Certificate
cat cert.pem | swk inspect cert
swk inspect cert --check-expiry < cert.pem

# Cron
swk inspect cron '*/5 * * * *'
swk inspect cron --explain '0 9 * * 1-5'
swk inspect cron --next 3 '0 9 * * MON'

# Text stats
cat essay.txt | swk inspect text
echo 'hello world' | swk inspect text --json
```

### Query (`swk query`)

| Command | Alias | Description |
|---------|-------|-------------|
| `query html` | `q html` | Query HTML with CSS selectors |
| `query json` | `q jp` | Query JSON with JSONPath |
| `query regex` | `q re` | Match/replace with regular expressions |

```bash
# HTML (CSS selectors)
curl -s https://example.com | swk query html -q 'a' --attr href
cat page.html | swk query html -q 'div.content p'

# JSON (JSONPath)
echo '{"users":[{"name":"Alice"},{"name":"Bob"}]}' | swk query json -q '$.users[*].name'

# Regex
echo '2024-01-15 hello 2024-02-20' | swk query regex -p '\d{4}-\d{2}-\d{2}' -g
echo 'John:30' | swk query regex -p '(\w+):(\d+)' --groups
echo 'foo bar baz' | swk query regex -p 'bar' -r 'qux'
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

# Fetch API → format → inspect
curl -s https://api.example.com/data | swk convert json | swk inspect text

# CSV → JSON → query with JSONPath
cat data.csv | swk convert json --from csv | swk query json -q '$..[?(@.age>30)]'

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
make lint          # run go vet
```

## License

MIT
