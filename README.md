# spotify-widget
[![Go Report Card](https://goreportcard.com/badge/github.com/otaleghani/spotify-widget)](https://goreportcard.com/report/github.com/otaleghani/spotify-widget)
![Ci tests](https://github.com/otaleghani/spotify-widget/actions/workflows/tests.yml/badge.svg)

A simple spotify widget for you Github profile.

<img src="https://spotify_image.talesign.com/image?some=thing" width="420">

## Overview

`spotify-widget` logs in your account via the Spotify APIs and retrives your playback information every 10 seconds. If it finds a playing song it display a lable that says "CURRENTLY PLAYING", else it displays "LAST LISTENED TO".

This program saves the Spotify OAuth2 tokens locally under `~/HOME/.cache/spotify-widget/auth.json` and every time the access token expires it gets a new one by using the refresh token. If I'm not mistaken the refresh token from Spotify does not expire. So you'll need to rotate it manually.

## Installation

You will need Go at least 1.2x. Afterwards you can go install spotify-widget with:

``` bash
go install github.com/otaleghani/sbes
```

## How to use

### You'll need

- Go v1.2+
- A Spotify Premium account
- A new app from [Spotify for Developers](https://developer.spotify.com/)
- A running server to serve the output image and to login

N.B.: This program uses port :9000 for the login and port :9001 to serve the image. At the moment is not possible to change the port.

### Steps

1. Create a new [Spotify app](https://developer.spotify.com/dashboard)
2. Specify under "Redirect URIs" your server domain with "/callback" path. So you'll need to add http://your-domain.com/callback
3. You can now save the app and retrive it's "Client ID" and "Client Secret"
4. Install onto your server spotify-widget (`go install spotify-widget`)
5. Run command `spotify-widget -i "YOUR_SPOTIFY_CLIENT_ID" -s "YOUR_SPOTIFY_CLIENT_SECRET -d "http://yourdomain.com`. This will prompt the program to visit the port :9000 to complete the login.
6. Visit the address and login with your Spotify account
7. Now the server will start on port :9001 to serve your image at `:9001/image`!

## About the server

As I previously mentioned, you'll need to expose two ports, :9000 for the Spotify login and :9001 for serving the image. Spotify does not care too much about the address. It accepts `http`, it accepts `:9000/callback`. You can even use your IPv4 address directly, it would not budge. 

If you want to have this little program with your domain you could reverse proxy it.

I personally use [Caddy](https://caddyserver.com/) to reverse proxy the :9000 and :9001 ports to my domain. You could do something like this:

``` Caddyfile
login.yourdomain.com {
    reverse_proxy localhost:9000
}
image.yourdomain.com {
    reverse_proxy localhost:9001
}
```

## Future development

- [ ] Top songs
- [ ] More widget styles
- [ ] Add link to song
