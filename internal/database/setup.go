package database

import (
	"fmt"
	"os"
	"path/filepath"
)

func Setup(id, secret string) error {
	authFilepath, err := getAuthPath()
	if err != nil {
		return err
	}

	if _, err = os.Stat(authFilepath); os.IsNotExist(err) {
		auth := AuthData{
			ClientId:     id,
			ClientSecret: secret,
		}
		if err = writeAuthFile(auth); err != nil {
			return err
		}
		return nil
	} else {
		auth, err := OpenAuthFile()
		if err != nil {
			return err
		}
		auth.ClientId = id
		auth.ClientSecret = secret
		if err = writeAuthFile(auth); err != nil {
			return err
		}
	}

	fontBoldFilepath, err := FontBoldFilepath()
	if err != nil {
		return err
	}
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
