package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AuthData struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

func SaveToken(access, refresh string) error {
	auth, err := openAuthFile()
	if err != nil {
		return err
	}

	auth.AccessToken = access
	auth.RefreshToken = refresh

	if err = writeAuthFile(auth); err != nil {
		return err
	}
	return nil
}

func openAuthFile() (AuthData, error) {
	filePath, err := getAuthPath()
	if err != nil {
		return AuthData{}, err
	}

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return AuthData{}, err
	}

	auth := AuthData{}
	err = json.Unmarshal(rawData, &auth)
	if err != nil {
		return AuthData{}, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return auth, nil
}

func writeAuthFile(auth AuthData) error {
	filePath, err := getAuthPath()
	if err != nil {
		return err
	}

	encodedData, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("could not encode data: %v", err)
	}

	if err := ioutil.WriteFile(filePath, encodedData, 0644); err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}
	fmt.Println("saved this: ", auth)

	return nil
}
