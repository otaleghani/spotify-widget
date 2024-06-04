package image

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"

	"github.com/otaleghani/spotify-widget/internal/database"
)

func CreateImage(trackName, artistName string) error {
	// Load the image
	imageFilepath, err := database.CurrentImageFilepath()
	if err != nil {
		return err
	}

	imgFile, err := os.Open(filepath.Clean(imageFilepath))
	if err != nil {
		return err
	}
	defer imgFile.Close()
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return err
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

	fontBoldFilepath, err := database.FontBoldFilepath()
	if err != nil {
		return err
	}
	fontRegularFilepath, err := database.FontRegularFilepath()
	if err != nil {
		return err
	}

	center := 50.0
	// Draw the text on the right
	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace(fontBoldFilepath, 48) // Replace with your font file
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, trackName, W-250), 220, center+48)

	dc.SetRGB(0.5, 0.5, 0.5)
	err = dc.LoadFontFace(fontRegularFilepath, 24) // Replace with your font file
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, artistName, W-250), 220, center+72+24)

	// Save the result
	outputFilepath, err := database.WidgetImageFilepath()
	if err != nil {
		return err
	}
	err = dc.SavePNG(outputFilepath)
	if err != nil {
		return err
	}
	return nil
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
