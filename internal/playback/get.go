package playback

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/otaleghani/spotify-widget/internal/database"
	"github.com/otaleghani/spotify-widget/internal/image"
	"github.com/otaleghani/spotify-widget/internal/oauth2"
)

func saveResponseToFile(accessToken string) error {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %v", err)
	}

	err = database.UpdatePlayback(body)
	if err != nil {
		return err
	}

	return nil
}

func RefreshPlayback() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	auth, err := database.OpenAuthFile()
	if err != nil {
		fmt.Println("couldn't open auth file")
		return
	}

	for {
		<-ticker.C
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
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		if trackName == "" {
			continue
		}

		err = image.CreateImage(trackName, artistName)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
	}

}
