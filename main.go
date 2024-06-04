package main

import (
	"flag"
	"log"

	"github.com/otaleghani/spotify-widget/internal/database"
	"github.com/otaleghani/spotify-widget/internal/oauth2"
	"github.com/otaleghani/spotify-widget/internal/playback"
	"github.com/otaleghani/spotify-widget/internal/server"
)

func main() {
	var id = flag.String("i", "", "the client id of your spotify account")
	var secret = flag.String("s", "", "the secret id of your spotify account")
	var domain = flag.String("d", "", "the domain this software is currently sitting on")
	flag.Parse()

	// If id and secret are specified, save them
	if *id != "" && *secret != "" {
		if err := database.Setup(*id, *secret); err != nil {
			log.Println("Error: ", err)
			return
		}
	}

	// Checks if RefreshToken is valid
	valid, _ := oauth2.IsRefreshTokenValid()
	if !valid {
		if err := oauth2.GetOauth2(*domain); err != nil {
			log.Println("Error: ", err)
			return
		}
	}

	// Serves the image and starts the playback refresher
	go server.Serve()
	go playback.RefreshPlayback()
	select {}
}
