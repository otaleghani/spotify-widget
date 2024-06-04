package database

import (
	"fmt"
	"os"
	"path/filepath"
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

	auth, err := OpenAuthFile()
	if err != nil {
		return err
	}

	auth.ClientId = id
	auth.ClientSecret = secret

	if err = writeAuthFile(auth); err != nil {
		return err
	}

	fontBoldFilepath := filepath.Join(cacheDir, "assets", "Inter-Bold.ttf")
	// fontRegularFilepath := filepath.Join(cacheDir, "assets", "Inter-Regular.ttf")

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
	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return "", fmt.Errorf("could not create directory: %v", err)
	}

	return cacheDir, nil
}
