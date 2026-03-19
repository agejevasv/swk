package graphic

import (
	"bytes"
	"image/png"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	for _, style := range []string{"circles", "squares", "lines", "mixed", "random", ""} {
		t.Run(style, func(t *testing.T) {
			data, err := GenerateImage(200, 150, style)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("invalid PNG: %v", err)
			}
			b := img.Bounds()
			if b.Dx() != 200 || b.Dy() != 150 {
				t.Errorf("dimensions = %dx%d, want 200x150", b.Dx(), b.Dy())
			}
		})
	}
}

func TestGenerateImageNotBlank(t *testing.T) {
	data, err := GenerateImage(100, 100, "mixed")
	if err != nil {
		t.Fatal(err)
	}
	img, _ := png.Decode(bytes.NewReader(data))
	// Count non-background pixels
	nonBG := 0
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r>>8 != 30 || g>>8 != 30 || b>>8 != 30 {
				nonBG++
			}
		}
	}
	if nonBG < 100 {
		t.Errorf("image looks blank: only %d non-background pixels", nonBG)
	}
}

func TestGenerateImageInvalidStyle(t *testing.T) {
	_, err := GenerateImage(100, 100, "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}
