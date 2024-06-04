package image

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func createRoundedImage(img image.Image, w, h int, r float64) image.Image {
	resizedImg := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), r)
	dc.Clip()
	dc.DrawImageAnchored(resizedImg, w/2, h/2, 0.5, 0.5)
	return dc.Image()
}

func truncateText(dc *gg.Context, text string, maxWidth float64) string {
	width, _ := dc.MeasureString(text)
	if width <= maxWidth {
		return text
	}
	for len(text) > 0 {
		text = text[:len(text)-1]
		width, _ = dc.MeasureString(text + "...")
		if width <= maxWidth {
			return text + "..."
		}
	}
	return "..." // Fallback in case the text is too short to be truncated
}
