package database

import (
	"file/filepath"
	"io"
	"net/http"
)

func Setup(id, secret string) error {
	cacheDir, err := cacheDirectory()
	if err != nil {
		return err
	}
	authFilepath := filepath.Join(cacheDir, "auth.json")

	if _, err = os.Stat(authFilepath); os.IsNotExist(err) {
		auth := AuthData{
			ClientId:     "",
			ClientSecret: "",
		}
		if err = writeAuthFile(auth); err != nil {
			return err
		}
		return nil
	}

	auth, err := openAuthFile()
	if err != nil {
		return err
	}

	auth.ClientId = id
	auth.ClientSecret = secret

	if err = writeAuthFile(auth); err != nil {
		return err
	}

	fontBoldFilepath := filepath.Join(cacheDir, "assets", "Inter-Bold.ttf")
	fontRegularFilepath := filepath.Join(cacheDir, "assets", "Inter-Regular.ttf")

	if _, err = os.Stat(fontBoldFilepath); os.IsNotExist(err) {
		err = downloadAssets()
		if err != nil {
			return err
		}
	}

	return nil
}

func cacheDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %v", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "spotify-widget")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("could not create directory: %v", err)
	}

	return filePath, nil
}

func downloadAssets() error {
	fontBoldFilepath := filepath.Join(cacheDir, "assets", "Inter-Bold.ttf")
	fontRegularFilepath := filepath.Join(cacheDir, "assets", "Inter-Regular.ttf")
	fontBoldUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Bold.ttf"
	fontRegularUrl := "https://github.com/otaleghani/spotify-widget/raw/main/assets/fonts/Inter-Regular.ttf"

	fontBoldFile, err := os.Create(fontRegularFilepath)
	if err != nil {
		return err
	}
	defer out.Close()
	fontRegularFile, err := os.Create(fontBoldFilepath)
	if err != nil {
		return err
	}
	defer out.Close()

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
