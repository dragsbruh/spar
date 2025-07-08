package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dragsbruh/spar/internal/api"
	"github.com/dragsbruh/spar/internal/downloader"
	"github.com/dragsbruh/spar/internal/listfile"
	"github.com/dragsbruh/spar/internal/misc"
	"github.com/dragsbruh/spar/internal/server"
	"github.com/dragsbruh/spar/internal/tokenutil"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	var (
		localOnly    bool
		listfilePath string
	)

	cmd := cli.Command{
		Name:  "spar",
		Usage: "spotify archiver",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "sets log level to debug",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "sync",
				Usage:   "sync a spar.yml listfile",
				Aliases: []string{"s"},

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "local",
						Usage:       "does not interact with spotify api (uses metadata file)",
						Aliases:     []string{"l"},
						Value:       false,
						Destination: &localOnly,
					},
					&cli.StringFlag{
						Name:        "listfile",
						Usage:       "listfile to use instead of spar.yml",
						Aliases:     []string{"file", "lf"},
						Value:       "spar.yml",
						Destination: &listfilePath,
					},
				},

				Action: func(ctx context.Context, c *cli.Command) error {
					conf, err := listfile.Load(listfilePath)
					if err != nil {
						return fmt.Errorf("error loading listfile: %v", err)
					}

					var tracks []spotify.FullTrack
					if !localOnly {
						log.Infof("preparing spotify api client")
						client, err := PrepareClient(ctx, 8080)
						if err != nil {
							return fmt.Errorf("error preparing client: %v", err)
						}

						log.Info("indexing items from api")
						for i, item := range conf.Items {
							itemTracks, err := api.GetItem(ctx, client, item, 200*time.Millisecond)
							if err != nil {
								return fmt.Errorf("getting item: %w", err)
							}
							tracks = append(tracks, itemTracks...)
							log.Infof("(%d/%d) got %d tracks (%d total) for `%s` (`%s`)", i+1, len(conf.Items), len(itemTracks), len(tracks), item.ID, item.Name)
						}

						log.Infof("saving tracks metadata to %s", conf.MetaPath)
						if err := misc.SaveTracks(tracks, conf.MetaPath); err != nil {
							return fmt.Errorf("saving tracks metadata: %w", err)
						}
					} else {
						log.Infof("loading tracks metadata from %s", conf.MetaPath)
						tracks, err = misc.LoadTracks(conf.MetaPath)
						if err != nil {
							return fmt.Errorf("loading tracks metadata: %w", err)
						}
					}

					log.Info("downloading tracks")
					if err := downloader.DownloadTracks(tracks, conf.TempDirectory, conf.OutDirectory, conf.Workers); err != nil {
						return fmt.Errorf("downloading tracks: %w", err)
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func PrepareClient(ctx context.Context, port int) (*spotify.Client, error) {
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientId == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID is not set")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_SECRET is not set")
	}

	redirectUrl := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	auth := spotifyauth.New(spotifyauth.WithClientID(clientId), spotifyauth.WithClientSecret(clientSecret), spotifyauth.WithRedirectURL(redirectUrl))

	log.Infof("attempting to load token")
	token, err := tokenutil.Load()
	if err != nil {
		return nil, fmt.Errorf("loading token: %w", err)
	}

	if token == nil || !token.Valid() {
		log.Warnf("token non-existent/invalid, reauthenticating")
		token, err = server.StartWaitForToken(auth, port)
		if err != nil {
			return nil, fmt.Errorf("token server: %w", err)
		}

		log.Infof("saving token")
		if err := tokenutil.Save(token); err != nil {
			return nil, fmt.Errorf("saving token: %w", err)
		}
	}

	log.Infof("creating client")
	return spotify.New(auth.Client(ctx, token)), nil
}
