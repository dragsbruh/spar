package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/dragsbruh/spar.git/internal/music"
	"github.com/dragsbruh/spar.git/internal/utils"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	list, err := music.GetLocalList("list.csv")
	if err != nil {
		log.Fatalf("Error loading list.csv: %v", err)
	}
	log.Infof("Loaded %d items from list.csv", len(list))

	acquiredTracks := []spotify.FullTrack{}

	log.Info("Preparing client")
	client := utils.PrepareClient(ctx)

	for _, item := range list {
		if strings.HasPrefix(item.Name, "#") {
			continue
		}
		tracks := music.GetItemTracks(ctx, client, item, 250*time.Millisecond)
		acquiredTracks = append(acquiredTracks, tracks...)

		log.Infof("Acquired %d tracks for item `%s` (id: `%s`, total acquired: %d)", len(tracks), item.Name, item.Id, len(acquiredTracks))
	}

	data, err := json.MarshalIndent(acquiredTracks, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	if err := os.WriteFile("data.json", data, os.ModePerm); err != nil {
		log.Fatalf("Error writing JSON: %v", err)
	}

	log.Infof("Completed and saved to data.json")
}
