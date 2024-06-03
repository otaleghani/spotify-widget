package playback

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

  "github.com/otaleghani/spotify-widget/internal/database"
)

// SaveResponseToFile fetches the content from localhost:8080 using the
// provided access token and saves it to
// $HOME/.cache/spotify-widget/playback.json

func saveResponseToFile(accessToken string) error {
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return fmt.Errorf("could not get home directory: %v", err)
	// }

	// cacheDir := filepath.Join(homeDir, ".cache", "spotify-widget")
	// if err := os.MkdirAll(cacheDir, 0755); err != nil {
	// 	return fmt.Errorf("could not create directory: %v", err)
	// }

	req, err := http.NewRequest(
		"GET",
		"https://api.spotify.com/v1/me/player/currently-playing",
		nil,
	)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("401")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %v", err)
	}

  err = database.UpdatePlayback(body)
  if err != nil {
    return err
  }

	// filePath := filepath.Join(cacheDir, "playback.json")
	// if err := ioutil.WriteFile(filePath, body, 0644); err != nil {
	// 	return fmt.Errorf("could not write file: %v", err)
	// }

	return nil
}

func RefreshPlayback() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	auth, err := database.openAuthFile()
	if err != nil {
		fmt.Println("couldn't open auth file")
		return
	}

	for {
		select {
		case <-ticker.C:
			err := saveResponseToFile(auth.AccessToken)
			if err != nil {
				if err.Error() == "401" {
					log.Println("Access token invalid, getting a new one...")
          accessToken, err := oauth2.GetNewAccessToken(auth.ClientId, auth.ClientSecret, auth.RefreshToken)
          if err != nil {
            log.Printf("Error oauth2.GetNewAccessToken: %v", err)
            break
          }
          err = database.SaveToken(accessToken, auth.RefreshToken)
          if err != nil {
            log.Printf("Error oauth2.SaveToken: %v", err)
            break
          }
          continue
				}
				log.Printf("Error saving response to file: %v", err)
			} else {
				log.Println("Response saved to $HOME/.cache/spotify-widget/playback.json")
			}
			trackName, artistName, err := parseSpotifyData()
			CreateImage(trackName, artistName)
			if err != nil {
				log.Println("Error: ", err)
			}
		}
	}
}
