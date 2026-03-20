package graphic

import (
	"strings"
	"testing"
)

func TestConvertColorHexToFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
		toFmt string
		want  string
	}{
		// #FF0000 (red)
		{name: "red to rgb", input: "#FF0000", toFmt: "rgb", want: "rgb(255,0,0)"},
		{name: "red to hsl", input: "#FF0000", toFmt: "hsl", want: "hsl(0,100%,50%)"},
		{name: "red to cmyk", input: "#FF0000", toFmt: "cmyk", want: "cmyk(0,100,100,0)"},

		// #00FF00 (green)
		{name: "green to rgb", input: "#00FF00", toFmt: "rgb", want: "rgb(0,255,0)"},
		{name: "green to hsl", input: "#00FF00", toFmt: "hsl", want: "hsl(120,100%,50%)"},

		// #0000FF (blue)
		{name: "blue to rgb", input: "#0000FF", toFmt: "rgb", want: "rgb(0,0,255)"},
		{name: "blue to hsl", input: "#0000FF", toFmt: "hsl", want: "hsl(240,100%,50%)"},

		// #FFFFFF (white)
		{name: "white to rgb", input: "#FFFFFF", toFmt: "rgb", want: "rgb(255,255,255)"},
		{name: "white to hsl", input: "#FFFFFF", toFmt: "hsl", want: "hsl(0,0%,100%)"},
		{name: "white to cmyk", input: "#FFFFFF", toFmt: "cmyk", want: "cmyk(0,0,0,0)"},

		// #000000 (black)
		{name: "black to rgb", input: "#000000", toFmt: "rgb", want: "rgb(0,0,0)"},
		{name: "black to cmyk", input: "#000000", toFmt: "cmyk", want: "cmyk(0,0,0,100)"},
		{name: "black to hsl", input: "#000000", toFmt: "hsl", want: "hsl(0,0%,0%)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertColor(tt.input, "hex", tt.toFmt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertColor(%q, hex, %q) = %q, want %q", tt.input, tt.toFmt, got, tt.want)
			}
		})
	}
}

func TestConvertColorRGBToHex(t *testing.T) {
	got, err := ConvertColor("255,0,0", "rgb", "hex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "#FF0000" {
		t.Errorf("got %q, want #FF0000", got)
	}
}

func TestConvertColorAutoDetect(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "detect hex", input: "#FF0000", want: "#FF0000"},
		{name: "detect rgb func", input: "rgb(255,0,0)", want: "#FF0000"},
		{name: "detect hsl func", input: "hsl(240,100%,50%)", want: "rgb(0,0,255)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toFmt := "hex"
			if tt.name == "detect hsl func" {
				toFmt = "rgb"
			}
			got, err := ConvertColor(tt.input, "auto", toFmt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertColorInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		from  string
	}{
		{name: "invalid color string auto", input: "notacolor", from: "auto"},
		{name: "invalid hex chars", input: "#ZZZZZZ", from: "hex"},
		{name: "too short hex", input: "#FF", from: "hex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertColor(tt.input, tt.from, "rgb")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestConvertColorAllFormats(t *testing.T) {
	got, err := ConvertColor("#FF0000", "hex", "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"#FF0000", "rgb(255,0,0)", "hsl(", "hsv(", "cmyk("} {
		if !strings.Contains(got, want) {
			t.Errorf("'all' output missing %q.\nGot:\n%s", want, got)
		}
	}
}

func TestConvertColorHSLRoundtrip(t *testing.T) {
	// Convert red: hex -> hsl -> hex
	hsl, err := ConvertColor("#FF0000", "hex", "hsl")
	if err != nil {
		t.Fatalf("hex->hsl: %v", err)
	}
	hex, err := ConvertColor(hsl, "hsl", "hex")
	if err != nil {
		t.Fatalf("hsl->hex: %v", err)
	}
	if hex != "#FF0000" {
		t.Errorf("roundtrip got %q, want #FF0000", hex)
	}
}

func TestConvertColorCMYKRoundtrip(t *testing.T) {
	// White: hex -> cmyk -> hex
	cmyk, err := ConvertColor("#FFFFFF", "hex", "cmyk")
	if err != nil {
		t.Fatalf("hex->cmyk: %v", err)
	}
	hex, err := ConvertColor(cmyk, "cmyk", "hex")
	if err != nil {
		t.Fatalf("cmyk->hex: %v", err)
	}
	if hex != "#FFFFFF" {
		t.Errorf("roundtrip got %q, want #FFFFFF", hex)
	}
}

func TestConvertColorHSVToFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
		toFmt string
		want  string
	}{
		// Pure red: hsv(0, 100%, 100%) -> rgb(255,0,0)
		{name: "red to hex", input: "hsv(0, 100%, 100%)", toFmt: "hex", want: "#FF0000"},
		{name: "red to rgb", input: "hsv(0, 100%, 100%)", toFmt: "rgb", want: "rgb(255,0,0)"},

		// Green: hsv(120, 100%, 50%) -> rgb(0,128,0) approximately
		{name: "green to hex", input: "hsv(120, 100%, 50%)", toFmt: "hex", want: "#008000"},
		{name: "green to rgb", input: "hsv(120, 100%, 50%)", toFmt: "rgb", want: "rgb(0,128,0)"},

		// Pure green: hsv(120, 100%, 100%) -> rgb(0,255,0)
		{name: "pure green to hex", input: "hsv(120, 100%, 100%)", toFmt: "hex", want: "#00FF00"},

		// Pure blue: hsv(240, 100%, 100%) -> rgb(0,0,255)
		{name: "blue to hex", input: "hsv(240, 100%, 100%)", toFmt: "hex", want: "#0000FF"},
		{name: "blue to rgb", input: "hsv(240, 100%, 100%)", toFmt: "rgb", want: "rgb(0,0,255)"},

		// Yellow: hsv(60, 100%, 100%) -> rgb(255,255,0)
		{name: "yellow to hex", input: "hsv(60, 100%, 100%)", toFmt: "hex", want: "#FFFF00"},

		// Cyan: hsv(180, 100%, 100%) -> rgb(0,255,255)
		{name: "cyan to hex", input: "hsv(180, 100%, 100%)", toFmt: "hex", want: "#00FFFF"},

		// Magenta: hsv(300, 100%, 100%) -> rgb(255,0,255)
		{name: "magenta to hex", input: "hsv(300, 100%, 100%)", toFmt: "hex", want: "#FF00FF"},

		// White: hsv(0, 0%, 100%) -> rgb(255,255,255)
		{name: "white to hex", input: "hsv(0, 0%, 100%)", toFmt: "hex", want: "#FFFFFF"},

		// Black: hsv(0, 0%, 0%) -> rgb(0,0,0)
		{name: "black to hex", input: "hsv(0, 0%, 0%)", toFmt: "hex", want: "#000000"},

		// Gray (zero saturation): hsv(0, 0%, 50%) -> rgb(128,128,128)
		{name: "gray to hex", input: "hsv(0, 0%, 50%)", toFmt: "hex", want: "#808080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertColor(tt.input, "hsv", tt.toFmt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ConvertColor(%q, hsv, %q) = %q, want %q", tt.input, tt.toFmt, got, tt.want)
			}
		})
	}
}

func TestConvertColorHSVAutoDetect(t *testing.T) {
	got, err := ConvertColor("hsv(0, 100%, 100%)", "auto", "hex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "#FF0000" {
		t.Errorf("got %q, want #FF0000", got)
	}
}

func TestConvertColorHSVRoundtrip(t *testing.T) {
	// Convert red: hex -> hsv -> hex
	hsv, err := ConvertColor("#FF0000", "hex", "hsv")
	if err != nil {
		t.Fatalf("hex->hsv: %v", err)
	}
	hex, err := ConvertColor(hsv, "hsv", "hex")
	if err != nil {
		t.Fatalf("hsv->hex: %v", err)
	}
	if hex != "#FF0000" {
		t.Errorf("roundtrip got %q, want #FF0000", hex)
	}
}

func TestConvertColorHSVInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "wrong part count", input: "hsv(0, 100%)"},
		{name: "invalid hue", input: "hsv(abc, 100%, 100%)"},
		{name: "invalid saturation", input: "hsv(0, xyz%, 100%)"},
		{name: "invalid value", input: "hsv(0, 100%, bad%)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertColor(tt.input, "hsv", "hex")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestConvertColorRGBFuncToHex(t *testing.T) {
	got, err := ConvertColor("rgb(0,255,0)", "rgb", "hex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "#00FF00" {
		t.Errorf("got %q, want #00FF00", got)
	}
}
