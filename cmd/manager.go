package cmd

import (
	"fmt"
	"runtime"

	"github.com/marcosalvi-01/wallman/db/sqlc"
	"github.com/marcosalvi-01/wallman/hyprpaper"
	"github.com/marcosalvi-01/wallman/macos"
)

type Manager interface {
	Next() error
	Previous() error
	Random(trueRandom bool) error
	Current() (string, error)
	History() ([]string, error)
	Set(path string) error
}

func GetManager(wallpaperDirs []string, travelSubdir bool, managerType string, queries *sqlc.Queries, dryRun bool) (Manager, error) {
	if managerType == "" || managerType == "auto" {
		if runtime.GOOS == "darwin" {
			managerType = "mac"
		} else {
			managerType = "hyprpaper"
		}
	}
	switch managerType {
	case "hyprpaper":
		return hyprpaper.New(wallpaperDirs, travelSubdir, queries, dryRun)
	case "mac":
		return macos.New(wallpaperDirs, travelSubdir, queries, dryRun)
	default:
		return nil, fmt.Errorf("unsupported manager type: %s", managerType)
	}
}
