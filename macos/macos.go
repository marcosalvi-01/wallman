// Package macos provides wallpaper management for macOS using osascript.
package macos

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/marcosalvi-01/wallman/cmd/common"
	"github.com/marcosalvi-01/wallman/db"
	"github.com/marcosalvi-01/wallman/db/sqlc"
)

var configDir = filepath.Join(os.Getenv("HOME"), ".local", "share", "wallman")

type MacOS struct {
	configDir     string
	wallpaperDirs []string
	wallpapers    []string
	queries       *sqlc.Queries
	dryRun        bool
}

func New(wallpaperDirs []string, travelSubdirs bool, queries *sqlc.Queries, dryRun bool) (*MacOS, error) {
	walls := make([]string, 0)

	for _, wallpaperDir := range wallpaperDirs {
		if travelSubdirs {
			err := filepath.WalkDir(wallpaperDir, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return fmt.Errorf("failed to access %s: %v", path, err)
				}
				if d.IsDir() || !common.IsImage(d.Name()) {
					return nil
				}
				walls = append(walls, path)
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to walk directory %s: %w", wallpaperDir, err)
			}
		} else {
			entries, err := os.ReadDir(wallpaperDir)
			if err != nil {
				return nil, fmt.Errorf("failed to read wallpaper directory: %w", err)
			}

			for _, entry := range entries {
				if entry.IsDir() || !common.IsImage(entry.Name()) {
					continue
				}
				walls = append(walls, filepath.Join(wallpaperDir, entry.Name()))
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

	common.ShuffleSlice(randWalls)

	randomFile := filepath.Join(configDir, "random")

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

	return &MacOS{
		configDir:     configDir,
		wallpaperDirs: wallpaperDirs,
		wallpapers:    walls,
		queries:       queries,
		dryRun:        dryRun,
	}, nil
}

func (m *MacOS) Next() error {
	if len(m.wallpapers) == 0 {
		return fmt.Errorf("no wallpapers available")
	}

	current, err := db.GetCurrentWallpaperPath()
	index := -1
	if err == nil {
		for i, w := range m.wallpapers {
			if w == current {
				index = i
				break
			}
		}
	}

	nextIndex := 0
	if index != -1 {
		nextIndex = (index + 1) % len(m.wallpapers)
	}

	path := m.wallpapers[nextIndex]

	err = db.SetWallpaper(path)
	if err != nil {
		return err
	}

	if !m.dryRun {
		err := setWallpaper(path)
		if err != nil {
			return fmt.Errorf("failed to set next wallpaper: %w", err)
		}
	}

	return nil
}

func (m *MacOS) Previous() error {
	path, err := db.GetPreviousWallpaper()
	if err != nil {
		return err
	}

	err = db.SetWallpaper(path)
	if err != nil {
		return err
	}

	if !m.dryRun {
		err := setWallpaper(path)
		if err != nil {
			return fmt.Errorf("failed to set previous wallpaper: %w", err)
		}
	}

	return nil
}

func (m *MacOS) Random(trueRandom bool) error {
	if len(m.wallpapers) == 0 {
		return fmt.Errorf("no wallpapers available")
	}

	if trueRandom {
		bigInt := big.NewInt(int64(len(m.wallpapers)))
		randInt, randErr := crand.Int(crand.Reader, bigInt)
		if randErr != nil {
			return randErr
		}
		randomIndex := int(randInt.Int64())
		path := m.wallpapers[randomIndex]

		err := db.SetWallpaper(path)
		if err != nil {
			return err
		}

		if !m.dryRun {
			err := setWallpaper(path)
			if err != nil {
				return fmt.Errorf("failed to set random wallpaper: %w", err)
			}
		}

		return nil
	}

	shuffled, index, err := db.GetRandomCycle()
	if err != nil {
		shuffled = make([]string, len(m.wallpapers))
		copy(shuffled, m.wallpapers)
		common.ShuffleSlice(shuffled)
		index = 0
		err = db.UpsertRandomCycle(shuffled, index)
		if err != nil {
			return fmt.Errorf("failed to initialize random cycle: %w", err)
		}
	}

	valid := len(shuffled) == len(m.wallpapers)
	if valid {
		for _, s := range shuffled {
			if !slices.Contains(m.wallpapers, s) {
				valid = false
				break
			}
		}
	}
	if !valid {
		shuffled = make([]string, len(m.wallpapers))
		copy(shuffled, m.wallpapers)
		common.ShuffleSlice(shuffled)
		index = 0
		err = db.UpsertRandomCycle(shuffled, index)
		if err != nil {
			return fmt.Errorf("failed to reset random cycle: %w", err)
		}
	}

	path := shuffled[index]
	err = db.SetWallpaper(path)
	if err != nil {
		return err
	}

	if !m.dryRun {
		err := setWallpaper(path)
		if err != nil {
			return fmt.Errorf("failed to set random wallpaper: %w", err)
		}
	}

	index++
	if index >= len(shuffled) {
		common.ShuffleSlice(shuffled)
		index = 0
	}
	err = db.UpsertRandomCycle(shuffled, index)
	if err != nil {
		return fmt.Errorf("failed to update random cycle: %w", err)
	}

	return nil
}

func (m *MacOS) Current() (string, error) {
	return db.GetCurrentWallpaperPath()
}

func (m *MacOS) History() ([]string, error) {
	history, err := db.GetWallpaperHistory(100)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(history))
	for i, h := range history {
		paths[i] = h.Path
	}
	return paths, nil
}

func (m *MacOS) Set(path string) error {
	path = common.ExpandPath(path)

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access wallpaper file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file")
	}

	if !common.IsImage(filepath.Base(path)) {
		return fmt.Errorf("unsupported image format (only JPEG, PNG, BMP, WEBP are supported)")
	}

	err = db.SetWallpaper(path)
	if err != nil {
		return fmt.Errorf("failed to set wallpaper in database: %w", err)
	}

	if !m.dryRun {
		err := setWallpaper(path)
		if err != nil {
			return fmt.Errorf("failed to set wallpaper: %w", err)
		}
	}

	return nil
}

func setWallpaper(path string) error {
	cmd := fmt.Sprintf(`tell application "System Events" to set picture of every desktop to POSIX file "%s"`, path)
	err := exec.Command("osascript", "-e", cmd).Run()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}
	return nil
}
