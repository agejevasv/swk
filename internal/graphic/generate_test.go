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

func TestGenerateImage_RejectsOutOfRangeDimensions(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
	}{
		{"zero", 0, 0},
		{"negative", -10, -10},
		{"below drawing minimum", 3, 3},
		{"one side too small", 100, 2},
		{"absurdly large", 1 << 20, 1 << 20},
		{"one side too large", 100, MaxImageDim + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GenerateImage(tt.width, tt.height, "mixed"); err == nil {
				t.Errorf("GenerateImage(%d, %d) should have returned an error", tt.width, tt.height)
			}
		})
	}
}

func TestGenerateImage_AcceptsSmallestAllowedSize(t *testing.T) {
	got, err := GenerateImage(MinImageDim, MinImageDim, "mixed")
	if err != nil {
		t.Fatalf("GenerateImage: %v", err)
	}
	if len(got) == 0 {
		t.Error("expected PNG data")
	}
}
