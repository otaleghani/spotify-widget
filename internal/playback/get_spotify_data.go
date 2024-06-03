package playback

import (
  "encoding/json"
  "os"
  "fmt"
  "path/filepath"
)

type Data struct {
    Item Item `json:"item"`
}

type Item struct {
    Album Album `json:"album"`
    Name string `json:"name"`
}

type Album struct {
    Artists []Artist `json:"artists"`
    Images []Image `json:"images"`
    Name string `json:"name"`
}

type Image struct {
  Url string `json:"url"`
}

type Artist struct {
    Name string `json:"name"`
}



func parseSpotifyData() (string, string, error) {
  // 1. loads data in memory
  homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("could not get home directory: %v", err)
	}
	filePath := filepath.Join(homeDir, ".cache", "spotify-widget", "playback.json")
  rawFile, err := os.ReadFile(filePath)
  if err != nil {
    return "", "", fmt.Errorf("could not decode file")
  }
  var data Data
  err = json.Unmarshal(rawFile, &data)
  if err != nil {
    return "", "", fmt.Errorf("error decoding json")
  }

  artistName := data.Item.Album.Artists[0].Name
  trackName := data.Item.Name
  imageUrl := data.Item.Album.Images[0].Url

  fmt.Println(artistName, trackName, imageUrl)
  return trackName, artistName, nil
}
