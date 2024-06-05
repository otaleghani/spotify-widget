package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	"github.com/otaleghani/spotify-widget/internal/database"
)

var (
	oauth2Config = &oauth2.Config{
		RedirectURL: "",
		Scopes:      []string{"user-read-currently-playing"},
		Endpoint:    spotify.Endpoint,
	}
	tokenChan = make(chan *oauth2.Token)
)

// generates a random string of the specified length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomString(16)
	if err != nil {
		http.Error(w, "Could not generate state", http.StatusInternalServerError)
		return
	}

	globalState = state
	url := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

var globalState string

func handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != globalState {
		http.Error(w, "State is not valid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Could not get token", http.StatusInternalServerError)
		return
	}

	tokenChan <- token
}

func GetOauth2(domain string) error {
	auth, err := database.OpenAuthFile()
	if err != nil {
		return err
	}

	oauth2Config.RedirectURL = domain + "/callback"
	oauth2Config.ClientID = auth.ClientId
	oauth2Config.ClientSecret = auth.ClientSecret

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/callback", handleCallback)

	srv := &http.Server{
		Addr:         ":9000",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("Server is starting at", src.Addr)
		log.Fatal(srv.ListenAndServe())
	}()

	token := <-tokenChan

	err = database.SaveToken(token.AccessToken, token.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}
