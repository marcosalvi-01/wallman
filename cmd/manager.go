package cmd

import "wallman/hyprpaper"

type Manager interface {
	Next() error
	Previous() error
	Random() error
	Current() (string, error)
	History() ([]string, error)
}

func GetManager(wallpaperDirs []string, travelSubdir bool) (Manager, error) {
	// TODO switch between the types of managers

	return hyprpaper.New(wallpaperDirs, travelSubdir)
}
