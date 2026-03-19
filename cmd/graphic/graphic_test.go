package graphic

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pflag "github.com/spf13/pflag"
)

func resetFlagChanged(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.String(), err
}

func executeCommandBinary(args ...string) ([]byte, error) {
	buf := new(bytes.Buffer)
	Cmd.SetOut(buf)
	Cmd.SetErr(buf)
	Cmd.SetArgs(args)
	err := Cmd.Execute()
	return buf.Bytes(), err
}

func createTestPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test PNG: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode test PNG: %v", err)
	}
	return path
}

// --- Color tests ---

func TestColorHexToAll(t *testing.T) {
	t.Cleanup(func() {
		colorFrom = "auto"
		colorTo = "all"
		resetFlagChanged(colorCmd.Flags())
	})

	out, err := executeCommand("color", "#ff0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Default --to is "all", should contain multiple formats
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "rgb") && !strings.Contains(lower, "hsl") {
		t.Fatalf("expected multiple color formats in output, got: %s", out)
	}
}

func TestColorRGBToHex(t *testing.T) {
	t.Cleanup(func() {
		colorFrom = "auto"
		colorTo = "all"
		resetFlagChanged(colorCmd.Flags())
	})

	out, err := executeCommand("color", "--from", "rgb", "--to", "hex", "rgb(255, 0, 0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lower := strings.ToLower(strings.TrimSpace(out))
	if !strings.Contains(lower, "ff0000") && !strings.Contains(lower, "#ff0000") {
		t.Fatalf("expected hex color output containing 'ff0000', got: %s", out)
	}
}

func TestColorAutoDetect(t *testing.T) {
	t.Cleanup(func() {
		colorFrom = "auto"
		colorTo = "all"
		resetFlagChanged(colorCmd.Flags())
	})

	out, err := executeCommand("color", "#00ff00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty output for auto-detected color")
	}
}

// --- Image tests ---

func TestImageConvertToJPEG(t *testing.T) {
	t.Cleanup(func() {
		imageToFormat = ""
		imageQuality = 85
		imageResize = ""
		imageInput = ""
		imageOutput = ""
		resetFlagChanged(imageCmd.Flags())
	})

	inputPath := createTestPNG(t)
	outputPath := filepath.Join(t.TempDir(), "out.jpg")

	_, err := executeCommand("image", "--to", "jpeg", "--input", inputPath, "--output", outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	// JPEG files start with FF D8 FF
	if len(data) < 3 || data[0] != 0xFF || data[1] != 0xD8 || data[2] != 0xFF {
		t.Fatal("output file does not have JPEG header bytes")
	}
}

func TestImageResize(t *testing.T) {
	t.Cleanup(func() {
		imageToFormat = ""
		imageQuality = 85
		imageResize = ""
		imageInput = ""
		imageOutput = ""
		resetFlagChanged(imageCmd.Flags())
	})

	inputPath := createTestPNG(t)
	outputPath := filepath.Join(t.TempDir(), "resized.png")

	_, err := executeCommand("image", "--to", "png", "--input", inputPath, "--output", outputPath, "--resize", "5x5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("failed to open output: %v", err)
	}
	defer f.Close()

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatalf("failed to decode PNG config: %v", err)
	}
	if cfg.Width != 5 || cfg.Height != 5 {
		t.Fatalf("expected 5x5 image, got %dx%d", cfg.Width, cfg.Height)
	}
}

func TestImageInputOutputWithTempFiles(t *testing.T) {
	t.Cleanup(func() {
		imageToFormat = ""
		imageQuality = 85
		imageResize = ""
		imageInput = ""
		imageOutput = ""
		resetFlagChanged(imageCmd.Flags())
	})

	inputPath := createTestPNG(t)
	outputPath := filepath.Join(t.TempDir(), "converted.png")

	_, err := executeCommand("image", "--to", "png", "--input", inputPath, "--output", outputPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("expected output file to exist")
	}
}
