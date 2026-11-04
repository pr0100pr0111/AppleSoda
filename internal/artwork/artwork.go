package artwork

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	lipgloss "charm.land/lipgloss/v2"
)

func Render(data []byte, cols, rows int, persistentID string) (string, error) {
	if len(data) == 0 || cols <= 0 || rows <= 0 {
		return "", nil
	}

	cacheDir, err := cacheDirectory()
	if err != nil {
		return "", err
	}

	cacheKey := fmt.Sprintf("%s_%d_%d", persistentID, cols, rows)
	cachePath := filepath.Join(cacheDir, cacheKey+".txt")

	if cached, err := os.ReadFile(cachePath); err == nil {
		return string(cached), nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resized := resizeImage(img, cols, rows*2)
	rendered := renderHalfBlock(resized)

	_ = os.MkdirAll(cacheDir, 0755)
	_ = os.WriteFile(cachePath, []byte(rendered), 0644)

	return rendered, nil
}

func cacheDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "applesoda", "artwork"), nil
}

func resizeImage(src image.Image, targetW, targetH int) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Max.X - srcBounds.Min.X
	srcH := srcBounds.Max.Y - srcBounds.Min.Y

	if srcW == 0 || srcH == 0 {
		return image.NewNRGBA(image.Rect(0, 0, targetW, targetH))
	}

	scaleX := float64(srcW) / float64(targetW)
	scaleY := float64(srcH) / float64(targetH)

	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	dstW := int(float64(srcW) / scale)
	dstH := int(float64(srcH) / scale)

	if dstW > targetW {
		dstW = targetW
	}
	if dstH > targetH {
		dstH = targetH
	}
	if dstW <= 0 {
		dstW = 1
	}
	if dstH <= 0 {
		dstH = 1
	}

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			srcX := int(float64(x)*scale) + srcBounds.Min.X
			srcY := int(float64(y)*scale) + srcBounds.Min.Y
			if srcX >= srcBounds.Max.X {
				srcX = srcBounds.Max.X - 1
			}
			if srcY >= srcBounds.Max.Y {
				srcY = srcBounds.Max.Y - 1
			}
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func renderHalfBlock(img image.Image) string {
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	rows := h / 2
	if rows == 0 {
		rows = 1
	}

	var out []byte
	for row := 0; row < rows; row++ {
		topY := bounds.Min.Y + row*2
		botY := topY + 1

		for x := bounds.Min.X; x < bounds.Min.X+w; x++ {
			topC := toRGBA(img.At(x, topY))
			var botC color.RGBA
			if botY < bounds.Max.Y {
				botC = toRGBA(img.At(x, botY))
			}

			fg := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", topC.R, topC.G, topC.B))
			bg := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", botC.R, botC.G, botC.B))

			cell := lipgloss.NewStyle().Foreground(fg).Background(bg).Render("▀")
			out = append(out, cell...)
		}
		if row < rows-1 {
			out = append(out, '\n')
		}
	}
	return string(out)
}

func toRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return color.RGBA{}
	}
	return color.RGBA{
		R: uint8(r * 255 / a),
		G: uint8(g * 255 / a),
		B: uint8(b * 255 / a),
		A: uint8(a >> 8),
	}
}
