package hyprpaper

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/marcosalvi-01/wallman/db"
	"github.com/marcosalvi-01/wallman/db/sqlc"
)

var (
	configDir  = filepath.Join(os.Getenv("HOME"), ".local", "share", "wallman")
	imageRegex = regexp.MustCompile(`^.*\.(jpeg|png)$`)
)

type Hyprpaper struct {
	configDir     string
	wallpaperDirs []string
	wallpapers    []string
	queries       *sqlc.Queries
	dryRun        bool
}

func New(wallpaperDirs []string, travelSubdirs bool, queries *sqlc.Queries, dryRun bool) (*Hyprpaper, error) {
	walls := make([]string, 0)

	for _, wallpaperDir := range wallpaperDirs {
		if travelSubdirs {
			if err := filepath.WalkDir(wallpaperDir, func(dirPath string, d fs.DirEntry, err error) error {
				if err != nil {
					log.Printf("warning: failed to access %s: %v", dirPath, err)
					return nil
				}
				if d.IsDir() || !isImg(d.Name()) {
					return nil
				}
				walls = append(walls, dirPath)
				return nil
			}); err != nil {
				return nil, fmt.Errorf("failed to walk directory %s: %w", wallpaperDir, err)
			}
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
		wallpapers:    walls,
		queries:       queries,
		dryRun:        dryRun,
	}, nil
}

func (h *Hyprpaper) Next() error {
	if len(h.wallpapers) == 0 {
		return fmt.Errorf("no wallpapers available")
	}

	current, err := db.GetCurrentWallpaperPath()
	index := -1
	if err == nil {
		for i, w := range h.wallpapers {
			if w == current {
				index = i
				break
			}
		}
	}

	nextIndex := 0
	if index != -1 {
		nextIndex = (index + 1) % len(h.wallpapers)
	}

	path := h.wallpapers[nextIndex]

	err = db.SetWallpaper(path)
	if err != nil {
		return err
	}

	if !h.dryRun {
		err := setWallpaperToAllMonitors(path, "cover")
		if err != nil {
			return fmt.Errorf("failed to set next wallpaper on all monitors: %w", err)
		}
	}

	return nil
}

func (h *Hyprpaper) Previous() error {
	path, err := db.GetPreviousWallpaper()
	if err != nil {
		return err
	}

	err = db.SetWallpaper(path)
	if err != nil {
		return err
	}

	if !h.dryRun {
		err := setWallpaperToAllMonitors(path, "cover")
		if err != nil {
			return fmt.Errorf("failed to set previous wallpaper on all monitors: %w", err)
		}
	}

	return nil
}

func (h *Hyprpaper) Random(trueRandom bool) error {
	if len(h.wallpapers) == 0 {
		return fmt.Errorf("no wallpapers available")
	}

	if trueRandom {
		// Old behavior: pick completely random from all wallpapers
		bigInt := big.NewInt(int64(len(h.wallpapers)))
		randInt, randErr := crand.Int(crand.Reader, bigInt)
		if randErr != nil {
			return randErr
		}
		randomIndex := int(randInt.Int64())
		path := h.wallpapers[randomIndex]

		err := db.SetWallpaper(path)
		if err != nil {
			return err
		}

		if !h.dryRun {
			err := setWallpaperToAllMonitors(path, "cover")
			if err != nil {
				return fmt.Errorf("failed to set random wallpaper on all monitors: %w", err)
			}
		}

		return nil
	}

	// Cycle behavior
	shuffled, index, err := db.GetRandomCycle()
	if err != nil {
		// If no cycle, initialize it
		shuffled = make([]string, len(h.wallpapers))
		copy(shuffled, h.wallpapers)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		index = 0
		err = db.UpsertRandomCycle(shuffled, index)
		if err != nil {
			return fmt.Errorf("failed to initialize random cycle: %w", err)
		}
	}

	// Verify the current shuffled matches h.wallpapers (handle changes)
	valid := len(shuffled) == len(h.wallpapers)
	if valid {
		for _, s := range shuffled {
			found := false
			for _, w := range h.wallpapers {
				if s == w {
					found = true
					break
				}
			}
			if !found {
				valid = false
				break
			}
		}
	}
	if !valid {
		// Reset cycle
		shuffled = make([]string, len(h.wallpapers))
		copy(shuffled, h.wallpapers)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
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

	if !h.dryRun {
		err := setWallpaperToAllMonitors(path, "cover")
		if err != nil {
			return fmt.Errorf("failed to set random wallpaper on all monitors: %w", err)
		}
	}

	// Advance index
	index++
	if index >= len(shuffled) {
		// Cycle complete, reshuffle for next
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		index = 0
	}
	err = db.UpsertRandomCycle(shuffled, index)
	if err != nil {
		return fmt.Errorf("failed to update random cycle: %w", err)
	}

	return nil
}

func (h *Hyprpaper) Current() (string, error) {
	return db.GetCurrentWallpaperPath()
}

func (h *Hyprpaper) History() ([]string, error) {
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

func isImg(fileName string) bool {
	return imageRegex.Match([]byte(fileName))
}

func setWallpaperToAllMonitors(path, fit string) error {
	monitors, err := listMonitors()
	if err != nil {
		return fmt.Errorf("failed to list monitors: %w", err)
	}
	for _, monitor := range monitors {
		err := setWallpaper(path, monitor, fit)
		if err != nil {
			return fmt.Errorf("failed to set wallpaper on monitor %s: %w", monitor, err)
		}
	}
	return nil
}

func setWallpaper(path, monitor, fit string) error {
	args := fmt.Sprintf("%s,%s,%s", monitor, path, fit)
	err := exec.Command("hyprctl", "hyprpaper", "wallpaper", args).Run()
	if err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}
	return nil
}

func listMonitors() ([]string, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get monitors: %v", err)
	}

	var data []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("failed to parse monitors JSON: %v", err)
	}

	var names []string
	for _, m := range data {
		names = append(names, m.Name)
	}
	return names, nil
}
