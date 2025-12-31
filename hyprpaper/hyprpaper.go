package hyprpaper

import (
	"fmt"
	"io/fs"
	"math/rand/v2"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	configDir  = filepath.Join(os.Getenv("HOME"), ".local", "share", "wallman")
	imageRegex = regexp.MustCompile(`^.*\.(jpeg|png)$`)
)

type Hyprpaper struct {
	configDir     string
	wallpaperDirs []string
}

func New(wallpaperDirs []string, travelSubdirs bool) (*Hyprpaper, error) {
	walls := make([]string, 0)

	for _, wallpaperDir := range wallpaperDirs {
		if travelSubdirs {
			filepath.WalkDir(wallpaperDir, func(dirPath string, d fs.DirEntry, err error) error {
				if err != nil {
					// TODO show a warning
					return nil
				}
				if d.IsDir() || !isImg(d.Name()) {
					return nil
				}
				walls = append(walls, path.Join(dirPath, d.Name()))
				return nil
			})
		} else {
			entries, err := os.ReadDir(wallpaperDir)
			if err != nil {
				return nil, fmt.Errorf("failed to read wallpaper directory: %w", err)
			}

			for _, entry := range entries {
				if entry.IsDir() || !isImg(entry.Name()) {
					continue
				}
				walls = append(walls, path.Join(wallpaperDir, entry.Name()))
			}
		}
	}

	if _, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(configDir, 0o700)
			if err != nil {
				return nil, fmt.Errorf("failed to create config directory: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to access config directory: %w", err)
		}
	}

	randWalls := make([]string, len(walls))
	copy(randWalls, walls)

	rand.Shuffle(len(randWalls), func(i, j int) {
		randWalls[i], randWalls[j] = randWalls[j], randWalls[i]
	})

	randomFile := path.Join(configDir, "random")

	if _, err := os.Stat(randomFile); err != nil {
		if os.IsNotExist(err) {
			join := strings.Join(randWalls, "\n") + "\n"

			err := os.WriteFile(randomFile, []byte(join), 0o600)
			if err != nil {
				return nil, fmt.Errorf("failed to write random wallpaper file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to access config directory: %w", err)
		}
	}

	return &Hyprpaper{
		configDir:     configDir,
		wallpaperDirs: wallpaperDirs,
	}, nil
}

func (h *Hyprpaper) Next() error {
	return nil
}

func (h *Hyprpaper) Previous() error {
	return nil
}

func (h *Hyprpaper) Random() error {
	return nil
}

func (h *Hyprpaper) Current() (string, error) {
	return "", nil
}

func (h *Hyprpaper) History() ([]string, error) {
	return nil, nil
}

func isImg(fileName string) bool {
	return imageRegex.Match([]byte(fileName))
}
