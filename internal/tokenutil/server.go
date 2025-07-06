package tokenutil

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const state = "sparrr_rawr"

func GetNewToken(addr string, auth *spotifyauth.Authenticator) *oauth2.Token {
	tokenChan := make(chan *oauth2.Token)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}

		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authenticated"))

		if err := Save(TokenFile, tok); err != nil {
			log.Fatal(err)
		}
	})

	authUrl := auth.AuthURL(state)
	log.Println("Opening browser for authentication")
	if err := browser.OpenURL(authUrl); err != nil {
		log.Fatalf("Failed to open auth url in browser: %v", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	token := <-tokenChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	return token
}
