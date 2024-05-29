package playback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// getNewAccessToken refreshes the access token using the refresh token
func getNewAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
	tokenURL := spotify.Endpoint.TokenURL

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("could not create request: %v", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not fetch new access token: %v", err)
	}
	defer resp.Body.Close()

	var token oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("could not decode response: %v", err)
	}

	return token.AccessToken, nil
}
