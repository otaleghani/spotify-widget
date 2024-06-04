package database

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Downloads the needed assets
func downloadAssets() error {
	fontBoldFilepath, err := FontBoldFilepath()
	if err != nil {
		return err
	}
	fontRegularFilepath, err := FontRegularFilepath()
	if err != nil {
		return err
	}
	lltFilepath, err := LltFilepath()
	if err != nil {
		return err
	}
	cltFilepath, err := CltFilepath()
	if err != nil {
		return err
	}

	fontBoldUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Bold.ttf"
	fontRegularUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Regular.ttf"
	lltImageUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/currently-listening-to.png"
	cltImageUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/last-listened-to.png"

	fontBoldFile, err := os.Create(filepath.Clean(fontBoldFilepath))
	if err != nil {
		return err
	}
	defer fontBoldFile.Close()
	fontRegularFile, err := os.Create(filepath.Clean(fontRegularFilepath))
	if err != nil {
		return err
	}
	defer fontRegularFile.Close()
	lltFile, err := os.Create(filepath.Clean(lltFilepath))
	if err != nil {
		return err
	}
	defer lltFile.Close()
	cltFile, err := os.Create(filepath.Clean(cltFilepath))
	if err != nil {
		return err
	}
	defer cltFile.Close()

	// Get the data
	getFontBold, err := http.Get(fontBoldUrl)
	if err != nil {
		return err
	}
	defer getFontBold.Body.Close()
	getFontRegular, err := http.Get(fontRegularUrl)
	if err != nil {
		return err
	}
	defer getFontRegular.Body.Close()
	getLltImage, err := http.Get(lltImageUrl)
	if err != nil {
		return err
	}
	defer getLltImage.Body.Close()
	getCltImage, err := http.Get(cltImageUrl)
	if err != nil {
		return err
	}
	defer getCltImage.Body.Close()

	_, err = io.Copy(fontBoldFile, getFontBold.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(fontRegularFile, getFontRegular.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(lltFile, getLltImage.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(cltFile, getCltImage.Body)
	if err != nil {
		return err
	}

	return nil
}

func FontBoldFilepath() (string, error) {
	assetsDir, err := assetsDirectory()
	if err != nil {
		return "", err
	}
	fontBoldFilepath := filepath.Join(assetsDir, "Inter-Bold.ttf")
	return fontBoldFilepath, nil
}

func FontRegularFilepath() (string, error) {
	assetsDir, err := assetsDirectory()
	if err != nil {
		return "", err
	}
	fontRegularFilepath := filepath.Join(assetsDir, "Inter-Regular.ttf")
	return fontRegularFilepath, nil
}

func LltFilepath() (string, error) {
	assetsDir, err := assetsDirectory()
	if err != nil {
		return "", err
	}
	lltFilepath := filepath.Join(assetsDir, "last-listened-to.png")
	return lltFilepath, nil
}

func CltFilepath() (string, error) {
	assetsDir, err := assetsDirectory()
	if err != nil {
		return "", err
	}
	cltFilepath := filepath.Join(assetsDir, "currently-listening-to.png")
	return cltFilepath, nil
}

func assetsDirectory() (string, error) {
	cacheDir, err := cacheDirectory()
	if err != nil {
		return "", err
	}
	assetsDir := filepath.Join(cacheDir, "assets")
	if err := os.MkdirAll(assetsDir, 0750); err != nil {
		return "", fmt.Errorf("could not create directory: %v", err)
	}
	return assetsDir, nil
}
