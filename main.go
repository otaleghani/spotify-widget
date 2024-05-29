package main

import (
	"flag"
	"fmt"

	"github.com/otaleghani/spotify-widget/internal/playback"
)

func main() {
	var id = flag.String("i", "", "help message for flag n")
	var secret = flag.String("s", "", "help message for flag n")
	flag.Parse()

	// If id and secret are specified, save them
	if *id != "" && *secret != "" {
		if err := playback.SaveClientId(*id, *secret); err != nil {
			fmt.Println(err)
			return
		}
	}

  valid, _ := playback.IsRefreshTokenValid() 

  if !valid {
	  if err := playback.GetOauth2(); err != nil {
	  	fmt.Println(err)
	  	return
	  }
  }

	go playback.RefreshPlayback()
	select {}
}
