package graphic

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

func ConvertImage(input []byte, toFormat string, quality int, width, height int) ([]byte, error) {
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

	switch toFormat {
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
