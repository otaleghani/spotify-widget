package database

import (
  "file/filepath"
  "ioutil"
)

//
func PlaybackFilepath() (string, error) {
  cacheDir, err := cacheDirectory()
  if err != nil {
    return "", err
  }

	playbackFilepath := filepath.Join(homeDir, "playback.json")
  
  return playbackFilepath, nil
}

func UpdatePlayback(body []byte) error {
  filePath, err := PlaybackFilepath()
  if err != nil {
    return err
  }

	if err := ioutil.WriteFile(filePath, body, 0644); err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

  return nil
}
