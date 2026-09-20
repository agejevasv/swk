package graphic

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/draw"
)

func ConvertImage(input []byte, toFormat string, quality int, width, height int) ([]byte, error) {
	// Check the header before decoding: a small file can declare enormous
	// dimensions and the decode would allocate all of it.
	cfg, _, err := image.DecodeConfig(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	if cfg.Width > MaxImageDim || cfg.Height > MaxImageDim {
		return nil, fmt.Errorf("image is %dx%d, larger than the %dx%d limit",
			cfg.Width, cfg.Height, MaxImageDim, MaxImageDim)
	}

	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	if width > 0 && height > 0 {
		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		img = dst
	}

	var buf bytes.Buffer

	switch strings.ToLower(strings.TrimSpace(toFormat)) {
	case "png":
		err = png.Encode(&buf, img)
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	case "gif":
		err = gif.Encode(&buf, img, nil)
	default:
		return nil, fmt.Errorf("unsupported output format %q: must be png, jpeg, or gif", toFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("encoding image: %w", err)
	}

	return buf.Bytes(), nil
}
