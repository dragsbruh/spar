package utils

import (
	"context"
	"os"

	"github.com/dragsbruh/spar.git/internal/tokenutil"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func PrepareClient(ctx context.Context) *spotify.Client {
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	var auth = spotifyauth.New(spotifyauth.WithClientID(clientId), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithRedirectURL("http://127.0.0.1:8080/callback"))

	token := tokenutil.Load(tokenutil.TokenFile)
	if token == nil {
		token = tokenutil.GetNewToken(":8080", auth)
	} else {
		log.Println("Loaded token from", tokenutil.TokenFile)
	}

	return spotify.New(auth.Client(ctx, token))
}
