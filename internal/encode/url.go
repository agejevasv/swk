package encode

import (
	"net/url"
)

func URLEncode(input string, component bool) string {
	if component {
		return url.QueryEscape(input)
	}
	return url.PathEscape(input)
}

func URLDecode(input string, component bool) (string, error) {
	if component {
		return url.QueryUnescape(input)
	}
	return url.PathUnescape(input)
}
