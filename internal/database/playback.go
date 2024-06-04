package database

import (
	"fmt"
	"os"
	"path/filepath"
)

func PlaybackFilepath() (string, error) {
	cacheDir, err := cacheDirectory()
	if err != nil {
		return "", err
	}
	playbackFilepath := filepath.Join(cacheDir, "playback.json")
	return playbackFilepath, nil
}

func UpdatePlayback(body []byte) error {
	filePath, err := PlaybackFilepath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, body, 0600); err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	return nil
}
