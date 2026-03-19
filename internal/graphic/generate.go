package graphic

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
)

func randomPalette() []color.RGBA {
	count := 8 + rand.Intn(4)
	colors := make([]color.RGBA, count)
	step := 360.0 / float64(count)
	offset := rand.Float64() * 360
	for i := range colors {
		h := math.Mod(offset+float64(i)*step+(rand.Float64()-0.5)*step*0.5, 360)
		s := 0.55 + rand.Float64()*0.45
		v := 0.55 + rand.Float64()*0.45
		r, g, b := hsvToRGB(h, s, v)
		colors[i] = color.RGBA{r, g, b, 255}
	}
	// Shuffle so adjacent shapes don't get adjacent hues
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})
	return colors
}

func GenerateImage(width, height int, style string) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	pal := randomPalette()
	fillGradientBackground(img, width, height, pal)

	switch style {
	case "circles":
		drawCircles(img, width, height, pal)
	case "squares":
		drawSquares(img, width, height, pal)
	case "lines":
		drawLines(img, width, height, pal)
	case "mixed", "random", "":
		drawCircles(img, width, height, pal)
		drawSquares(img, width, height, pal)
		drawLines(img, width, height, pal)
	default:
		return nil, fmt.Errorf("unknown style %q (available: circles, squares, lines, mixed)", style)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func randomAlpha() uint8 {
	return uint8(100 + rand.Intn(156)) // 100-255
}

func pickColor(palette []color.RGBA) color.RGBA {
	c := palette[rand.Intn(len(palette))]
	c.A = randomAlpha()
	return c
}

func drawCircles(img *image.RGBA, w, h int, pal []color.RGBA) {
	count := 10 + rand.Intn(20)
	for i := 0; i < count; i++ {
		cx := rand.Intn(w)
		cy := rand.Intn(h)
		r := 10 + rand.Intn(min(w, h)/4)
		c := pickColor(pal)
		filled := rand.Intn(2) == 0
		if filled {
			fillCircle(img, cx, cy, r, c)
		} else {
			strokeCircle(img, cx, cy, r, 2+rand.Intn(3), c)
		}
	}
}

func drawSquares(img *image.RGBA, w, h int, pal []color.RGBA) {
	count := 8 + rand.Intn(15)
	for i := 0; i < count; i++ {
		x := rand.Intn(w)
		y := rand.Intn(h)
		size := 10 + rand.Intn(min(w, h)/4)
		c := pickColor(pal)
		filled := rand.Intn(2) == 0
		if filled {
			fillRect(img, x, y, size, size, c)
		} else {
			strokeRect(img, x, y, size, size, 2+rand.Intn(3), c)
		}
	}
}

func drawLines(img *image.RGBA, w, h int, pal []color.RGBA) {
	count := 10 + rand.Intn(15)
	for i := 0; i < count; i++ {
		x1 := rand.Intn(w)
		y1 := rand.Intn(h)
		x2 := rand.Intn(w)
		y2 := rand.Intn(h)
		thickness := 1 + rand.Intn(4)
		c := pickColor(pal)
		drawLine(img, x1, y1, x2, y2, thickness, c)
	}
}

func blendPixel(img *image.RGBA, x, y int, c color.RGBA) {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return
	}
	if c.A == 255 {
		img.Set(x, y, c)
		return
	}
	bg := img.RGBAAt(x, y)
	a := float64(c.A) / 255.0
	img.Set(x, y, color.RGBA{
		R: uint8(float64(c.R)*a + float64(bg.R)*(1-a)),
		G: uint8(float64(c.G)*a + float64(bg.G)*(1-a)),
		B: uint8(float64(c.B)*a + float64(bg.B)*(1-a)),
		A: 255,
	})
}

func fillCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			if dx*dx+dy*dy <= float64(r*r) {
				blendPixel(img, x, y, c)
			}
		}
	}
}

func strokeCircle(img *image.RGBA, cx, cy, r, thickness int, c color.RGBA) {
	rOuter := float64(r)
	rInner := float64(r - thickness)
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			dist := dx*dx + dy*dy
			if dist <= rOuter*rOuter && dist >= rInner*rInner {
				blendPixel(img, x, y, c)
			}
		}
	}
}

func fillRect(img *image.RGBA, x0, y0, w, h int, c color.RGBA) {
	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			blendPixel(img, x, y, c)
		}
	}
}

func strokeRect(img *image.RGBA, x0, y0, w, h, thickness int, c color.RGBA) {
	fillRect(img, x0, y0, w, thickness, c)
	fillRect(img, x0, y0+h-thickness, w, thickness, c)
	fillRect(img, x0, y0, thickness, h, c)
	fillRect(img, x0+w-thickness, y0, thickness, h, c)
}

func drawLine(img *image.RGBA, x1, y1, x2, y2, thickness int, c color.RGBA) {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	steps := int(math.Max(dx, dy))
	if steps == 0 {
		return
	}
	xInc := float64(x2-x1) / float64(steps)
	yInc := float64(y2-y1) / float64(steps)

	half := thickness / 2
	x, y := float64(x1), float64(y1)
	for i := 0; i <= steps; i++ {
		for ty := -half; ty <= half; ty++ {
			for tx := -half; tx <= half; tx++ {
				blendPixel(img, int(x)+tx, int(y)+ty, c)
			}
		}
		x += xInc
		y += yInc
	}
}

func fillGradientBackground(img *image.RGBA, w, h int, pal []color.RGBA) {
	seed := pal[rand.Intn(len(pal))]
	c1 := darken(seed, 0.12)
	c2 := darken(seed, 0.06)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Diagonal interpolation: 0 at top-left, 1 at bottom-right
			t := (float64(x)/float64(w) + float64(y)/float64(h)) / 2.0
			r := lerp(float64(c1.R), float64(c2.R), t)
			g := lerp(float64(c1.G), float64(c2.G), t)
			b := lerp(float64(c1.B), float64(c2.B), t)
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
}

func darken(c color.RGBA, factor float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: 255,
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
