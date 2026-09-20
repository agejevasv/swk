package graphic

import "testing"

func FuzzConvertColor(f *testing.F) {
	for _, s := range []string{"#FF5733", "FF5733", "rgb(255,87,51)", "hsl(11,100%,60%)", "hsv(11,80%,100%)",
		"cmyk(0,66,80,0)", "255,87,51", "", "#", "#FFF", "rgb()", "hsl(400,200%,-5%)"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		for _, to := range []string{"hex", "rgb", "hsl", "hsv", "cmyk", "all"} {
			_, _ = ConvertColor(s, "auto", to)
		}
	})
}

// maxRoundTripError is the largest per-channel error each colour space can
// introduce, measured exhaustively over all 16,777,216 RGB values. The loss is
// inherent: hsl/hsv/cmyk are rendered with integer degrees and percentages, so
// they cannot address every 24-bit colour. Any change that makes a space worse
// than this fails the round trip below.
var maxRoundTripError = map[string]int{
	"rgb":  0,
	"hsl":  5,
	"hsv":  3,
	"cmyk": 2,
}

// Every colour must survive a trip through every other colour space.
func FuzzColorRoundTrip(f *testing.F) {
	f.Add(uint8(255), uint8(87), uint8(51))
	f.Add(uint8(0), uint8(0), uint8(0))
	f.Add(uint8(255), uint8(255), uint8(255))
	f.Add(uint8(1), uint8(2), uint8(3))
	f.Fuzz(func(t *testing.T, r, g, b uint8) {
		hex := formatColor(r, g, b, "hex")

		for _, space := range []string{"rgb", "hsl", "hsv", "cmyk"} {
			rendered := formatColor(r, g, b, space)
			gotR, gotG, gotB, err := parseToRGB(rendered, space)
			if err != nil {
				t.Fatalf("%s rendering %q of %s does not parse back: %v", space, rendered, hex, err)
			}
			tolerance := maxRoundTripError[space]
			if diff(r, gotR) > tolerance || diff(g, gotG) > tolerance || diff(b, gotB) > tolerance {
				t.Fatalf("%s round trip of %s gave %s (via %q), beyond the measured tolerance of %d",
					space, hex, formatColor(gotR, gotG, gotB, "hex"), rendered, tolerance)
			}
		}
	})
}

func diff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}
