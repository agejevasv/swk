package encode

import (
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
)

var qrLevelMap = map[string]qrcode.RecoveryLevel{
	"L": qrcode.Low,
	"M": qrcode.Medium,
	"Q": qrcode.High,
	"H": qrcode.Highest,
}

func parseQRLevel(level string) (qrcode.RecoveryLevel, error) {
	lvl, ok := qrLevelMap[strings.ToUpper(level)]
	if !ok {
		return 0, fmt.Errorf("invalid QR error correction level %q: must be L, M, Q, or H", level)
	}
	return lvl, nil
}

func QRGenerate(input string, size int, level string) ([]byte, error) {
	lvl, err := parseQRLevel(level)
	if err != nil {
		return nil, err
	}
	return qrcode.Encode(input, lvl, size)
}

func QRTerminal(input string, level string) (string, error) {
	lvl, err := parseQRLevel(level)
	if err != nil {
		return "", err
	}

	q, err := qrcode.New(input, lvl)
	if err != nil {
		return "", err
	}

	bitmap := q.Bitmap()
	rows := len(bitmap)
	cols := 0
	if rows > 0 {
		cols = len(bitmap[0])
	}

	var sb strings.Builder

	for y := 0; y < rows; y += 2 {
		for x := 0; x < cols; x++ {
			top := bitmap[y][x]
			bottom := false
			if y+1 < rows {
				bottom = bitmap[y+1][x]
			}

			switch {
			case top && bottom:
				sb.WriteRune('█')
			case top && !bottom:
				sb.WriteRune('▀')
			case !top && bottom:
				sb.WriteRune('▄')
			default:
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String(), nil
}
