package fmt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// prefixScope tracks the xmlns declarations in scope so element and attribute
// names can be written back with the prefix the document used.
type prefixScope struct {
	stack []map[string]string // URI -> prefix ("" for the default namespace)
}

func newPrefixScope() *prefixScope {
	return &prefixScope{stack: []map[string]string{{}}}
}

func (p *prefixScope) lookup(uri string) (string, bool) {
	for i := len(p.stack) - 1; i >= 0; i-- {
		if prefix, ok := p.stack[i][uri]; ok {
			return prefix, true
		}
	}
	return "", false
}

func (p *prefixScope) rewrite(token xml.Token) xml.Token {
	switch t := token.(type) {
	case xml.StartElement:
		scope := map[string]string{}
		for _, a := range t.Attr {
			switch {
			case a.Name.Space == "xmlns":
				scope[a.Value] = a.Name.Local
			case a.Name.Space == "" && a.Name.Local == "xmlns":
				scope[a.Value] = ""
			}
		}
		p.stack = append(p.stack, scope)

		attrs := make([]xml.Attr, 0, len(t.Attr))
		for _, a := range t.Attr {
			a.Name = p.attrName(a.Name)
			attrs = append(attrs, a)
		}
		return xml.StartElement{Name: p.elementName(t.Name), Attr: attrs}

	case xml.EndElement:
		name := p.elementName(t.Name)
		if len(p.stack) > 1 {
			p.stack = p.stack[:len(p.stack)-1]
		}
		return xml.EndElement{Name: name}
	}

	return token
}

// elementName turns a resolved namespace URI back into the document's prefix.
// Leaving Name.Space set would make the encoder emit its own xmlns attribute
// alongside the one already present.
func (p *prefixScope) elementName(n xml.Name) xml.Name {
	if n.Space == "" {
		return xml.Name{Local: n.Local}
	}
	if prefix, ok := p.lookup(n.Space); ok {
		if prefix == "" {
			return xml.Name{Local: n.Local}
		}
		return xml.Name{Local: prefix + ":" + n.Local}
	}
	return xml.Name{Local: n.Local}
}

func (p *prefixScope) attrName(n xml.Name) xml.Name {
	switch n.Space {
	case "":
		return xml.Name{Local: n.Local}
	case "xmlns":
		return xml.Name{Local: "xmlns:" + n.Local}
	}
	if prefix, ok := p.lookup(n.Space); ok && prefix != "" {
		return xml.Name{Local: prefix + ":" + n.Local}
	}
	return xml.Name{Local: n.Local}
}

// XMLOptions holds options for XML formatting.
type XMLOptions struct {
	// Indent is the number of spaces per level; 0 produces unindented output.
	// Negative values are rejected.
	Indent int
	Minify bool
}

func FormatXML(input []byte, opts XMLOptions) ([]byte, error) {
	if opts.Indent < 0 {
		return nil, fmt.Errorf("indent must be >= 0, got %d", opts.Indent)
	}

	decoder := xml.NewDecoder(bytes.NewReader(input))
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)

	if opts.Minify {
		encoder.Indent("", "")
	} else {
		indent := strings.Repeat(" ", opts.Indent)
		encoder.Indent("", indent)
	}

	// The decoder reports namespaces as resolved URIs while also returning the
	// original xmlns attributes. Re-encoding that as-is emits a second xmlns
	// declaration and drops prefixes, so names are rebuilt from the
	// declarations actually in scope.
	prefixes := newPrefixScope()

	foundToken := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("invalid XML: %w", err)
		}
		foundToken = true

		// Drop whitespace that carries a newline: that is layout from an
		// already-formatted document, and it is being replaced. Spacing without
		// a newline sits between elements on one line and may be content.
		if cd, ok := token.(xml.CharData); ok {
			text := string(cd)
			if strings.TrimSpace(text) == "" && strings.ContainsAny(text, "\r\n") {
				continue
			}
		}

		if err := encoder.EncodeToken(prefixes.rewrite(token)); err != nil {
			return nil, err
		}
	}

	if !foundToken {
		return nil, fmt.Errorf("invalid XML: empty input")
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}

	result := buf.Bytes()

	if !opts.Minify && len(result) > 0 && result[len(result)-1] != '\n' {
		result = append(result, '\n')
	}

	return result, nil
}
