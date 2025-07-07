package listfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type ListfileItem struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
	ID   string `yaml:"id"`
}

type Listfile struct {
	OutDirectory  string         `yaml:"outdir"`
	TempDirectory string         `yaml:"tempdir"`
	MetaPath      string         `yaml:"metapath"`
	Items         []ListfileItem `yaml:"items"`
	Workers       int            `yaml:"workers"`
}

func Load(path string) (Listfile, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return Listfile{}, fmt.Errorf("reading listfile: %w", err)
	}

	var listfile Listfile
	if err := yaml.Unmarshal(source, &listfile); err != nil {
		return listfile, fmt.Errorf("unmarshalling yaml: %w", err)
	}

	if listfile.OutDirectory == "" {
		return listfile, fmt.Errorf("`outdir` not specified")
	}
	if listfile.TempDirectory == "" {
		return listfile, fmt.Errorf("`tempdir` not specified")
	}
	if listfile.MetaPath == "" {
		return listfile, fmt.Errorf("`metapath` not specified")
	}
	if len(listfile.Items) == 0 {
		return listfile, fmt.Errorf("`items` needs atleast one item")
	}
	if listfile.Workers < 1 {
		return listfile, fmt.Errorf("`workers` must be atleast 1")
	}

	for i, item := range listfile.Items {
		if item.Name == "" {
			return listfile, fmt.Errorf("`name` not specified for item %d (0 indexed)", i)
		}
		if item.Kind != "artist" && item.Kind != "playlist" && item.Kind != "track" {
			return listfile, fmt.Errorf("`kind` must be artist/playlist/track for item %d (0 indexed)", i)
		}
		if item.ID == "" {
			return listfile, fmt.Errorf("`id` not specified for item %d (0 indexed)", i)
		}
	}

	listfile.OutDirectory, err = ExpandAbs(listfile.OutDirectory)
	if err != nil {
		return listfile, fmt.Errorf("making outdir absolute: %w", err)
	}
	listfile.TempDirectory, err = ExpandAbs(listfile.TempDirectory)
	if err != nil {
		return listfile, fmt.Errorf("making tempdir absolute: %w", err)
	}
	listfile.MetaPath, err = ExpandAbs(listfile.MetaPath)
	if err != nil {
		return listfile, fmt.Errorf("making metapath absolute: %w", err)
	}

	if err := os.MkdirAll(listfile.TempDirectory, os.ModePerm); err != nil {
		return listfile, fmt.Errorf("making tempdir: %w", err)
	}
	if err := os.MkdirAll(listfile.OutDirectory, os.ModePerm); err != nil {
		return listfile, fmt.Errorf("making outdir: %w", err)
	}

	return listfile, nil
}

// 💪
func ExpandAbs(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		if path == "~" {
			path = home
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, path[2:])
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
