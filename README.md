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

## Commands

### Converters (`swk convert`)

| Command | Alias | Description |
|---------|-------|-------------|
| `convert json-yaml` | `c jy` | Convert between JSON and YAML |
| `convert json-csv` | `c jc` | Convert between JSON and CSV |
| `convert numbase` | `c nb` | Convert between number bases (bin, oct, dec, hex) |
| `convert datetime` | `c dt` | Convert between date/time formats |
| `convert cron` | `c cr` | Parse and explain cron expressions |

```bash
# JSON to YAML
echo '{"name":"swk","version":"1.0"}' | swk convert json-yaml

# YAML back to JSON
echo -e "name: swk\nversion: '1.0'" | swk convert json-yaml --reverse

# Number base conversion
swk convert numbase --from dec --to hex 255

# Unix timestamp to ISO
swk convert datetime --from unix --to iso 1700000000

# Current time
swk convert datetime now

# Explain a cron expression
swk convert cron --explain '*/5 * * * *'

# Show next 3 runs (default is 5)
swk convert cron --next 3 '0 9 * * MON'
```

### Encoders / Decoders (`swk encode`)

| Command | Alias | Description |
|---------|-------|-------------|
| `encode base64` | `enc b64` | Base64 encode/decode |
| `encode url` | `enc url` | URL encode/decode |
| `encode html` | `enc html` | HTML entity encode/decode |
| `encode jwt` | `enc jwt` | Create or decode JWT tokens |
| `encode gzip` | `enc gz` | Gzip compress/decompress |
| `encode cert` | `enc cert` | Decode X.509 PEM certificates |
| `encode qr` | `enc qr` | Generate QR codes |

```bash
# Base64 encode
echo 'hello world' | swk encode base64

# Base64 decode
echo 'aGVsbG8gd29ybGQ=' | swk encode base64 -d

# URL-safe base64
echo 'binary data here' | swk encode base64 --url-safe

# URL encode
echo 'hello world & friends' | swk encode url

# HTML encode
echo '<script>alert("xss")</script>' | swk encode html

# Create a JWT
swk encode jwt --secret mykey '{"sub":"user1","role":"admin"}'

# Create with HS512
swk encode jwt --secret mykey --algo HS512 '{"sub":"user1"}'

# Decode a JWT
swk encode jwt -d 'eyJhbGciOiJIUzI1NiIs...'

# Decode and verify signature
swk encode jwt -d --secret mykey 'eyJhbGciOiJIUzI1NiIs...'

# Inspect a certificate
cat cert.pem | swk encode cert

# Check if certificate is expired (exit code 1 if expired)
swk encode cert --check-expiry < cert.pem

# Generate a QR code in terminal
swk encode qr 'https://github.com/agejevasv/swk'

# Generate QR as PNG
swk encode qr --output png 'https://example.com' > qr.png
```

### Formatters (`swk fmt`)

| Command | Alias | Description |
|---------|-------|-------------|
| `fmt json` | `f j` | Pretty-print or minify JSON |
| `fmt xml` | `f x` | Pretty-print or minify XML |
| `fmt sql` | `f sql` | Format SQL queries |

```bash
# Pretty-print JSON
echo '{"a":1,"b":2}' | swk fmt json

# Minify JSON
cat config.json | swk fmt json --minify

# Custom indent (4 spaces)
cat data.json | swk fmt json --indent 4

# Format XML
cat messy.xml | swk fmt xml

# Minify XML
cat document.xml | swk fmt xml --minify

# Format SQL
echo "SELECT id,name FROM users WHERE active=1 ORDER BY name" | swk fmt sql

# Uppercase SQL keywords
echo "select * from users where id = 1" | swk fmt sql --uppercase
```

### Generators (`swk gen`)

| Command | Alias | Description |
|---------|-------|-------------|
| `gen uuid` | `g uid` | Generate UUIDs (v1, v4, v5, v7) |
| `gen hash` | `g h` | Generate hashes (MD5, SHA1, SHA256, SHA512) |
| `gen password` | `g pw` | Generate random passwords |
| `gen lorem` | `g li` | Generate lorem ipsum text |

```bash
# Generate a UUID (v4 by default)
swk gen uuid

# Generate 5 UUIDs
swk gen uuid --count 5

# UUID v7 (time-ordered)
swk gen uuid --version 7

# SHA256 hash (default)
echo -n 'hello' | swk gen hash

# MD5 hash
echo -n 'hello' | swk gen hash --algo md5

# Verify a hash
echo -n 'hello' | swk gen hash --verify 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824

# Generate a 32-character password
swk gen password --length 32

# Password with only letters and digits
swk gen password --no-symbols

# Generate 10 passwords
swk gen password --count 10

# Lorem ipsum paragraphs
swk gen lorem --paragraphs 3

# Lorem ipsum words
swk gen lorem --words 50
```

### Testers (`swk test`)

| Command | Alias | Description |
|---------|-------|-------------|
| `test regex` | `t re` | Test regex patterns against input |
| `test jsonpath` | `t jp` | Query JSON with JSONPath |
| `test xmlval` | `t xv` | Validate XML well-formedness |

```bash
# Test a regex
echo '2024-01-15 hello 2024-02-20' | swk test regex -p '\d{4}-\d{2}-\d{2}' -g

# Regex with capture groups
echo 'John:30' | swk test regex -p '(\w+):(\d+)' --groups

# Regex replace
echo 'foo bar baz' | swk test regex -p 'bar' -r 'qux'

# JSONPath query
echo '{"users":[{"name":"Alice"},{"name":"Bob"}]}' | swk test jsonpath -q '$.users[*].name'

# Validate XML
cat document.xml | swk test xmlval
```

### Text Utilities (`swk text`)

| Command | Alias | Description |
|---------|-------|-------------|
| `text inspect` | `txt info` | Character, word, line, byte counts |
| `text escape` | `txt esc` | Escape/unescape strings |
| `text case` | `txt case` | Convert between case conventions |
| `text diff` | `txt diff` | Unified diff between two files |
| `text md` | `txt md` | Render markdown |

```bash
# Inspect text
cat essay.txt | swk text inspect

# Inspect as JSON
echo 'hello world' | swk text inspect --json

# Escape for JSON
echo 'line1\nline2' | swk text escape --mode json

# Unescape
echo '\"hello\"' | swk text escape --mode json --unescape

# Shell escape
echo "it's a test" | swk text escape --mode shell

# Case conversion (preserves whitespace and structure)
echo 'hello world' | swk text case --to camel    # helloWorld
echo 'hello world' | swk text case --to pascal   # HelloWorld
echo 'hello world' | swk text case --to snake    # hello_world
echo 'hello world' | swk text case --to kebab    # hello-world
echo 'helloWorld'  | swk text case --to snake    # hello_world
cat file.go | swk text case --to upper           # uppercases file, preserves structure

# Diff two files
swk text diff -a old.txt -b new.txt

# Render markdown to styled HTML page
cat README.md | swk text md --html > preview.html

# Use a different syntax highlighting theme
cat README.md | swk text md --html --theme monokai > preview.html
```

### Graphic Tools (`swk graphic`)

| Command | Alias | Description |
|---------|-------|-------------|
| `graphic color` | `gfx col` | Convert between color formats |
| `graphic image` | `gfx img` | Convert image formats, resize |
| `graphic generate` | `gfx gen` | Generate placeholder images |

```bash
# Convert hex to all formats
swk graphic color '#FF5733'

# Convert RGB to hex
swk graphic color --from rgb --to hex '255,87,51'

# Convert image format
swk graphic image --to jpeg -i photo.png -o photo.jpg

# Resize image
swk graphic image --to png --resize 200x200 -i large.png -o thumb.png

# Convert with quality setting
swk graphic image --to jpeg --quality 75 -i input.png -o output.jpg

# Generate placeholder image
swk graphic generate -o placeholder.png

# Specific style and dimensions
swk graphic generate --style circles --width 1920 --height 1080 -o wallpaper.png

# Available styles: circles, squares, lines, mixed (default)
```

## Piping and chaining

Every command reads stdin and writes stdout:

```bash
# JSON → format → base64
echo '{"a":1}' | swk fmt json | swk encode base64

# Generate password → hash it
swk gen password | swk gen hash

# YAML → JSON → minify → clipboard (macOS)
cat config.yaml | swk convert json-yaml -r | swk fmt json -m | pbcopy

# Fetch API → format → inspect
curl -s https://www.cloudflarestatus.com/api/v2/status.json | swk fmt json | swk text inspect

# Uppercase a file preserving structure
cat main.go | swk text case --to upper

# CSV → JSON → query with JSONPath
cat data.csv | swk convert json-csv -r | swk test jsonpath -q '$..[?(@.age>30)]'
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
