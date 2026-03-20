package graphic

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ConvertColor(input string, fromFmt, toFmt string) (string, error) {
	if fromFmt == "auto" {
		fromFmt = detectColorFormat(input)
		if fromFmt == "" {
			return "", fmt.Errorf("cannot auto-detect color format for %q", input)
		}
	}

	r, g, b, err := parseToRGB(input, fromFmt)
	if err != nil {
		return "", err
	}

	if toFmt == "all" {
		var lines []string
		lines = append(lines, formatColor(r, g, b, "hex"))
		lines = append(lines, formatColor(r, g, b, "rgb"))
		lines = append(lines, formatColor(r, g, b, "hsl"))
		lines = append(lines, formatColor(r, g, b, "hsv"))
		lines = append(lines, formatColor(r, g, b, "cmyk"))
		return strings.Join(lines, "\n"), nil
	}

	return formatColor(r, g, b, toFmt), nil
}

func detectColorFormat(input string) string {
	if strings.HasPrefix(input, "#") {
		return "hex"
	}
	if strings.HasPrefix(input, "rgb(") {
		return "rgb"
	}
	if strings.HasPrefix(input, "hsl(") {
		return "hsl"
	}
	if strings.HasPrefix(input, "hsv(") {
		return "hsv"
	}
	if strings.HasPrefix(input, "cmyk(") {
		return "cmyk"
	}
	if len(input) == 6 {
		if _, err := strconv.ParseUint(input, 16, 32); err == nil {
			return "hex"
		}
	}
	if strings.Contains(input, ",") {
		return "rgb"
	}
	return ""
}

func parseToRGB(input string, format string) (uint8, uint8, uint8, error) {
	switch format {
	case "hex":
		return parseHex(input)
	case "rgb":
		return parseRGB(input)
	case "hsl":
		h, s, l, err := parseHSL(input)
		if err != nil {
			return 0, 0, 0, err
		}
		r, g, b := hslToRGB(h, s, l)
		return r, g, b, nil
	case "hsv":
		h, s, v, err := parseHSV(input)
		if err != nil {
			return 0, 0, 0, err
		}
		r, g, b := hsvToRGB(h, s, v)
		return r, g, b, nil
	case "cmyk":
		c, m, y, k, err := parseCMYK(input)
		if err != nil {
			return 0, 0, 0, err
		}
		r, g, b := cmykToRGB(c, m, y, k)
		return r, g, b, nil
	default:
		return 0, 0, 0, fmt.Errorf("unsupported color format %q", format)
	}
}

func parseHex(input string) (uint8, uint8, uint8, error) {
	hex := strings.TrimPrefix(input, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: must be 6 hex digits", input)
	}
	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hex color %q: %w", input, err)
	}
	r := uint8((val >> 16) & 0xFF)
	g := uint8((val >> 8) & 0xFF)
	b := uint8(val & 0xFF)
	return r, g, b, nil
}

func parseRGB(input string) (uint8, uint8, uint8, error) {
	s := input
	s = strings.TrimPrefix(s, "rgb(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid RGB color %q", input)
	}
	r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid RGB red value: %w", err)
	}
	g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid RGB green value: %w", err)
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid RGB blue value: %w", err)
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("RGB values must be 0-255, got (%d,%d,%d)", r, g, b)
	}
	return uint8(r), uint8(g), uint8(b), nil
}

func parseHSL(input string) (float64, float64, float64, error) {
	s := input
	s = strings.TrimPrefix(s, "hsl(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid HSL color %q", input)
	}
	h, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSL hue: %w", err)
	}
	sVal, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSL saturation: %w", err)
	}
	l, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSL lightness: %w", err)
	}
	if h < 0 || h > 360 || sVal < 0 || sVal > 100 || l < 0 || l > 100 {
		return 0, 0, 0, fmt.Errorf("HSL values out of range: hue must be 0-360, saturation/lightness must be 0-100")
	}
	return h, sVal / 100, l / 100, nil
}

func parseHSV(input string) (float64, float64, float64, error) {
	s := input
	s = strings.TrimPrefix(s, "hsv(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid HSV color %q", input)
	}
	h, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSV hue: %w", err)
	}
	sVal, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[1]), "%"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSV saturation: %w", err)
	}
	v, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(parts[2]), "%"), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid HSV value: %w", err)
	}
	if h < 0 || h > 360 || sVal < 0 || sVal > 100 || v < 0 || v > 100 {
		return 0, 0, 0, fmt.Errorf("HSV values out of range: hue must be 0-360, saturation/value must be 0-100")
	}
	return h, sVal / 100, v / 100, nil
}

func parseCMYK(input string) (float64, float64, float64, float64, error) {
	s := input
	s = strings.TrimPrefix(s, "cmyk(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("invalid CMYK color %q", input)
	}
	c, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid CMYK cyan: %w", err)
	}
	m, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid CMYK magenta: %w", err)
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid CMYK yellow: %w", err)
	}
	k, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid CMYK key: %w", err)
	}
	if c < 0 || c > 100 || m < 0 || m > 100 || y < 0 || y > 100 || k < 0 || k > 100 {
		return 0, 0, 0, 0, fmt.Errorf("CMYK values must be 0-100")
	}
	return c / 100, m / 100, y / 100, k / 100, nil
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		v := uint8(math.Round(l * 255))
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hk := h / 360

	tr := hk + 1.0/3.0
	tg := hk
	tb := hk - 1.0/3.0

	r := hueToRGB(p, q, tr)
	g := hueToRGB(p, q, tg)
	b := hueToRGB(p, q, tb)

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(b * 255))
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func rgbToHSL(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	l := (max + min) / 2

	if delta == 0 {
		return 0, 0, math.Round(l * 100)
	}

	var s float64
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}

	var h float64
	switch max {
	case rf:
		h = (gf - bf) / delta
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/delta + 2
	case bf:
		h = (rf-gf)/delta + 4
	}
	h *= 60

	return math.Round(h), math.Round(s * 100), math.Round(l * 100)
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	if s == 0 {
		val := uint8(math.Round(v * 255))
		return val, val, val
	}

	h = math.Mod(h, 360) / 60
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float64
	switch int(i) {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(b * 255))
}

func rgbToHSV(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	var h float64
	if delta == 0 {
		h = 0
	} else {
		switch max {
		case rf:
			h = (gf - bf) / delta
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/delta + 2
		case bf:
			h = (rf-gf)/delta + 4
		}
		h *= 60
	}

	var s float64
	if max != 0 {
		s = delta / max
	}

	return math.Round(h), math.Round(s * 100), math.Round(max * 100)
}

func cmykToRGB(c, m, y, k float64) (uint8, uint8, uint8) {
	r := 255 * (1 - c) * (1 - k)
	g := 255 * (1 - m) * (1 - k)
	b := 255 * (1 - y) * (1 - k)
	return uint8(math.Round(r)), uint8(math.Round(g)), uint8(math.Round(b))
}

func rgbToCMYK(r, g, b uint8) (float64, float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	k := 1 - math.Max(rf, math.Max(gf, bf))
	if k == 1 {
		return 0, 0, 0, 100
	}

	c := (1 - rf - k) / (1 - k)
	m := (1 - gf - k) / (1 - k)
	y := (1 - bf - k) / (1 - k)

	return math.Round(c * 100), math.Round(m * 100), math.Round(y * 100), math.Round(k * 100)
}

func formatColor(r, g, b uint8, format string) string {
	switch format {
	case "hex":
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	case "rgb":
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	case "hsl":
		h, s, l := rgbToHSL(r, g, b)
		return fmt.Sprintf("hsl(%.0f,%.0f%%,%.0f%%)", h, s, l)
	case "hsv":
		h, s, v := rgbToHSV(r, g, b)
		return fmt.Sprintf("hsv(%.0f,%.0f%%,%.0f%%)", h, s, v)
	case "cmyk":
		c, m, y, k := rgbToCMYK(r, g, b)
		return fmt.Sprintf("cmyk(%.0f,%.0f,%.0f,%.0f)", c, m, y, k)
	default:
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
}
