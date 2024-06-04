package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2/spotify"

	"github.com/otaleghani/spotify-widget/internal/database"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func IsRefreshTokenValid() (bool, error) {
	set := isRefreshTokenSet()
	if !set {
		return false, nil
	}
	auth, err := database.OpenAuthFile()
	if err != nil {
		return false, err
	}

	tokenURL := spotify.Endpoint.TokenURL

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", auth.RefreshToken)
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("could not create request: %v", err)
	}

	req.SetBasicAuth(auth.ClientId, auth.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not fetch new access token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return false, fmt.Errorf("could not decode error response: %v", err)
		}
		return false, fmt.Errorf("error: %s, description: %s", errorResponse.Error, errorResponse.ErrorDescription)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return false, fmt.Errorf("could not decode token response: %v", err)
	}

	fmt.Println(tokenResponse)
	// Here you'll need to save the accesstoken
	err = database.SaveToken(tokenResponse.AccessToken, auth.RefreshToken)
	if err != nil {
		return false, fmt.Errorf("database.SaveToken error: %v", err)
	}

	return true, nil
}

func isRefreshTokenSet() bool {
	auth, err := database.OpenAuthFile()
	if err != nil {
		return false
	}

	if auth.RefreshToken != "" {
		return true
	}
	return false
}
