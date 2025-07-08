package misc

import (
	"path/filepath"
	"strings"
)

var audioExts = map[string]bool{
	".mp3":  true,
	".m4a":  true,
	".ogg":  true,
	".opus": true,
	".flac": true,
	".wav":  true,
	".aac":  true,
	".alac": true,
	".webm": true,
}

func GetAudioFiles(base string) (string, error) {
	matches, err := filepath.Glob(base + ".*")
	if err != nil {
		return "", err
	}

	for _, path := range matches {
		ext := strings.ToLower(filepath.Ext(path))
		if audioExts[ext] {
			return path, nil
		}
	}

	return "", nil
}
