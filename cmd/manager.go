package cmd

import (
	"fmt"

	"wallman/db/sqlc"
	"wallman/hyprpaper"
)

type Manager interface {
	Next() error
	Previous() error
	Random() error
	Current() (string, error)
	History() ([]string, error)
}

func GetManager(wallpaperDirs []string, travelSubdir bool, managerType string, queries *sqlc.Queries, dryRun bool) (Manager, error) {
	switch managerType {
	case "", "auto", "hyprpaper":
		return hyprpaper.New(wallpaperDirs, travelSubdir, queries, dryRun)
	default:
		return nil, fmt.Errorf("unsupported manager type: %s", managerType)
	}
}
