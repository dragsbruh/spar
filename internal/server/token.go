package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// i dont like printing stuff to stdout that is NOT logging here
func StartWaitForToken(auth *spotifyauth.Authenticator, port int) (*oauth2.Token, error) {
	errChan := make(chan error)
	tokChan := make(chan *oauth2.Token)

	state := generateState(4)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(r.Context(), state, r)
		if err != nil {
			log.Errorf("couldnt get token: %v", err)
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			return
		}

		if st := r.FormValue("state"); st != state {
			log.Errorf("state mismatch, got %s", st)
			http.NotFound(w, r)
			return
		}

		log.Info("user authenticated")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authenticated"))

		tokChan <- tok
	})

	go func() {
		log.Info("starting http server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server closed: %w", err)
		}
	}()

	log.Infof("opening `%s` in browser", auth.AuthURL(state))
	if err := browser.OpenURL(auth.AuthURL(state)); err != nil {
		log.Errorf("error opening browser: %v", err)
		log.Warn("please open the above url to authenticate")
	}

	select {
	case token := <-tokChan:
		return token, nil
	case err := <-errChan:
		return nil, err
	}
}

func generateState(size int) string {
	state := ""
	for range size {
		state += strconv.Itoa(rand.Int())
	}
	return state
}
