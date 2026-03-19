package graphic

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func createTestPNG(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestConvertImagePNGToJPEG(t *testing.T) {
	input := createTestPNG(10, 10, color.RGBA{255, 0, 0, 255})
	out, err := ConvertImage(input, "jpeg", 90, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// JPEG files start with 0xFF 0xD8
	if len(out) < 2 || out[0] != 0xFF || out[1] != 0xD8 {
		t.Error("output does not start with JPEG magic bytes 0xFF 0xD8")
	}
}

func TestConvertImagePNGToPNG(t *testing.T) {
	input := createTestPNG(10, 10, color.RGBA{0, 255, 0, 255})
	out, err := ConvertImage(input, "png", 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PNG magic: 0x89 P N G
	if len(out) < 4 || out[0] != 0x89 || out[1] != 'P' || out[2] != 'N' || out[3] != 'G' {
		t.Error("output does not start with PNG magic bytes")
	}
}

func TestConvertImagePNGToGIF(t *testing.T) {
	input := createTestPNG(10, 10, color.RGBA{0, 0, 255, 255})
	out, err := ConvertImage(input, "gif", 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) < 3 || string(out[:3]) != "GIF" {
		t.Error("output does not start with GIF header")
	}
}

func TestConvertImageResize(t *testing.T) {
	input := createTestPNG(100, 100, color.RGBA{0, 0, 255, 255})
	out, err := ConvertImage(input, "png", 0, 50, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("failed to decode resized PNG: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("resized to %dx%d, want 50x50", bounds.Dx(), bounds.Dy())
	}
}

func TestConvertImageInvalidFormat(t *testing.T) {
	input := createTestPNG(10, 10, color.White)
	_, err := ConvertImage(input, "bmp", 0, 0, 0)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestConvertImageCorruptInput(t *testing.T) {
	_, err := ConvertImage([]byte("not an image at all"), "png", 0, 0, 0)
	if err == nil {
		t.Fatal("expected error for corrupt input, got nil")
	}
}

func TestConvertImageEmptyInput(t *testing.T) {
	_, err := ConvertImage([]byte{}, "png", 0, 0, 0)
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestConvertImageJPEGQuality(t *testing.T) {
	// Use a more complex image to see quality difference
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 2), uint8(y * 2), uint8((x + y)), 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	input := buf.Bytes()

	lowQ, err := ConvertImage(input, "jpeg", 1, 0, 0)
	if err != nil {
		t.Fatalf("low quality: %v", err)
	}
	highQ, err := ConvertImage(input, "jpeg", 100, 0, 0)
	if err != nil {
		t.Fatalf("high quality: %v", err)
	}
	if len(lowQ) >= len(highQ) {
		t.Errorf("low quality (%d bytes) should be smaller than high quality (%d bytes)", len(lowQ), len(highQ))
	}
}

func TestConvertImageResizeNonSquare(t *testing.T) {
	input := createTestPNG(200, 100, color.RGBA{255, 128, 0, 255})
	out, err := ConvertImage(input, "png", 0, 50, 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 25 {
		t.Errorf("resized to %dx%d, want 50x25", bounds.Dx(), bounds.Dy())
	}
}

func TestConvertImageJPGAlias(t *testing.T) {
	input := createTestPNG(10, 10, color.RGBA{255, 0, 0, 255})
	out, err := ConvertImage(input, "jpg", 85, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) < 2 || out[0] != 0xFF || out[1] != 0xD8 {
		t.Error("jpg alias: output does not start with JPEG magic bytes")
	}
}
