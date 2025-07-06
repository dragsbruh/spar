package music

import (
	"encoding/csv"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ListItemType int

const (
	ListItemPlaylist ListItemType = iota
	ListItemArtist
	ListItemTrack
)

type ListItem struct {
	Name string
	Kind ListItemType
	Id   string
}

func GetLocalList(path string) ([]ListItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	items := []ListItem{}

	for _, record := range records {
		if len(record) != 3 {
			log.Fatal("Expected only 3 columns for row, got", len(record))
		}
		item := ListItem{
			Name: record[0],
			Id:   record[2],
		}
		itemType := strings.ToLower(record[1])
		switch itemType {
		case "playlist":
			item.Kind = ListItemPlaylist
		case "artist":
			item.Kind = ListItemArtist
		case "track":
			item.Kind = ListItemTrack
		default:
			log.Fatal("expected playlist or artist, got", itemType, "for item type")
		}

		items = append(items, item)
	}

	return items, nil
}
