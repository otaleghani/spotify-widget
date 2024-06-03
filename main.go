package main

import (
	"flag"
	"fmt"

	"github.com/otaleghani/spotify-widget/internal/playback"
)

func main() {
	var id = flag.String("i", "", "the client id of your spotify account")
	var secret = flag.String("s", "", "the secret id of your spotify account")
	var domain = flag.String("d", "", "the domain this software is currently sitting on")
	flag.Parse()

	// If id and secret are specified, save them
	if *id != "" && *secret != "" {
		if err := database.Setup(*id, *secret); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Check if oauthfile is present

	valid, _ := playback.IsRefreshTokenValid()

	if !valid {
		if err := playback.GetOauth2(*domain); err != nil {
			fmt.Println(err)
			return
		}
	}

	go playback.RefreshPlayback()
	select {}
}
