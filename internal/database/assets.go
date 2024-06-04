package database

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadAssets() error {
	fontBoldFilepath, err := FontBoldFilepath()
	if err != nil {
		return err
	}
	fontRegularFilepath, err := FontRegularFilepath()
	if err != nil {
		return err
	}

	fontBoldUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Bold.ttf"
	fontRegularUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Regular.ttf"

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

	_, err = io.Copy(fontBoldFile, getFontBold.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(fontRegularFile, getFontRegular.Body)
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
