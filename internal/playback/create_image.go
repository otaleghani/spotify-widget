package playback

import (
	"image"
	"image/color"
	// "image/draw"
	"image/jpeg"
	"os"

	"github.com/fogleman/gg"
  "github.com/nfnt/resize"
)

func CreateImage(trackName, artistName string) {
	// Load the image
	imgFile, err := os.Open("input.jpg") // Replace with your image file
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	// Create a new image with a white background
	const W = 800
	const H = 220
	dc := gg.NewContext(W, H)
	dc.SetColor(color.Black)
	dc.Clear()

	// Create the rounded image
	const squareSize = 200
	radius := float64(squareSize) / 24
	roundedImage := createRoundedImage(img, squareSize, squareSize, radius)

	// Draw the rounded image on the left
	dc.DrawImage(roundedImage, 10, 10)


  center := 50.0
	// Draw the text on the right
	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace("Inter-Bold.ttf", 48) // Replace with your font file
	dc.DrawString(truncateText(dc, trackName, W-250), 220, center + 48)

	dc.SetRGB(0.5, 0.5, 0.5)
	dc.LoadFontFace("Inter-Regular.ttf", 24) // Replace with your font file
	dc.DrawString(truncateText(dc, artistName, W-250), 220, center + 72 + 24)

	// Save the result
	dc.SavePNG("output.png")
}

func createRoundedImage(img image.Image, w, h int, r float64) image.Image {
  resizedImg := resize.Thumbnail(uint(w), uint(h), img, resize.Lanczos3)

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
