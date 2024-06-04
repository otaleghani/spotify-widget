package image

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/otaleghani/spotify-widget/internal/database"
)

func CurrentlyListeningTo(trackName, artistName string) error {
	// Load cover image
	coverImageFilepath, err := database.CurrentImageFilepath()
	if err != nil {
		return err
	}
	coverImgFile, err := os.Open(filepath.Clean(coverImageFilepath))
	if err != nil {
		return err
	}
	defer coverImgFile.Close()
	coverImg, err := jpeg.Decode(coverImgFile)
	if err != nil {
		return err
	}

	// Load background image
	cltFilepath, err := database.CltFilepath()
	if err != nil {
		return err
	}

	bgImageFile, err := os.Open(filepath.Clean(cltFilepath))
	if err != nil {
		return err
	}
	defer bgImageFile.Close()
	bgImg, err := png.Decode(bgImageFile)
	if err != nil {
		return err
	}

	W := 1170
	H := 260
	dc := gg.NewContext(W, H)
	dc.DrawImage(bgImg, 0, 0)

	const squareSize = 196
	radius := float64(squareSize) / 16
	roundedImage := createRoundedImage(coverImg, squareSize, squareSize, radius)
	//roundedImage := createRoundedImage(coverImg, squareSize, squareSize, radius)
	dc.DrawImage(roundedImage, 16*2, 16*2)

	fontBoldFilepath, err := database.FontBoldFilepath()
	if err != nil {
		return err
	}
	fontRegularFilepath, err := database.FontRegularFilepath()
	if err != nil {
		return err
	}

	center := 90.0
	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace(fontBoldFilepath, 40)
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, trackName, float64(W-300)), 267, center+48)

	dc.SetRGB(0.5, 0.5, 0.5)
	err = dc.LoadFontFace(fontRegularFilepath, 40)
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, artistName, float64(W-300)), 267, center+72+24)

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

func LastListenedTo(trackName, artistName string) error {
	// Load cover image
	coverImageFilepath, err := database.CurrentImageFilepath()
	if err != nil {
		return err
	}
	coverImgFile, err := os.Open(filepath.Clean(coverImageFilepath))
	if err != nil {
		return err
	}
	defer coverImgFile.Close()
	coverImg, err := jpeg.Decode(coverImgFile)
	if err != nil {
		return err
	}

	lltFilepath, err := database.LltFilepath()
	if err != nil {
		return err
	}

	// Load background image
	bgImageFile, err := os.Open(filepath.Clean(lltFilepath))
	if err != nil {
		return err
	}
	defer bgImageFile.Close()
	bgImg, err := png.Decode(bgImageFile)
	if err != nil {
		return err
	}

	W := 1170
	H := 260
	dc := gg.NewContext(W, H)
	dc.DrawImage(bgImg, 0, 0)

	const squareSize = 196
	radius := float64(squareSize) / 16
	roundedImage := createRoundedImage(ResizeAndCropToSquare(coverImg, squareSize), squareSize, squareSize, radius)
	dc.DrawImage(roundedImage, 16*2, 16*2)

	fontBoldFilepath, err := database.FontBoldFilepath()
	if err != nil {
		return err
	}
	fontRegularFilepath, err := database.FontRegularFilepath()
	if err != nil {
		return err
	}

	center := 90.0
	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace(fontBoldFilepath, 40)
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, trackName, float64(W-300)), 267, center+48)

	dc.SetRGB(0.5, 0.5, 0.5)
	err = dc.LoadFontFace(fontRegularFilepath, 40)
	if err != nil {
		return err
	}
	dc.DrawString(truncateText(dc, artistName, float64(W-300)), 267, center+72+24)

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

// ResizeAndCropToSquare resizes and crops an image to the given square size
func ResizeAndCropToSquare(im image.Image, size int) image.Image {
	// Get the minimum dimension to make a square
	minDim := im.Bounds().Dx()
	if im.Bounds().Dy() < minDim {
		minDim = im.Bounds().Dy()
	}

	// Center crop the image to a square
	rect := image.Rect(
		(im.Bounds().Dx()-minDim)/2,
		(im.Bounds().Dy()-minDim)/2,
		(im.Bounds().Dx()+minDim)/2,
		(im.Bounds().Dy()+minDim)/2,
	)
	cropped := image.NewRGBA(rect)
	draw.Draw(cropped, rect, im, rect.Min, draw.Src)

	// Resize the cropped image to the desired square size
	dc := gg.NewContext(size, size)
	dc.DrawImage(cropped, 0, 0)
	return dc.Image()
}
