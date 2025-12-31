package common

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func List(dirs []string, subdirs bool) ([]string, error) {
	walls := make([]string, 0)

	for _, wallpaperDir := range dirs {
		if subdirs {
			filepath.WalkDir(wallpaperDir, func(dirPath string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() || !isImg(d.Name()) {
					return nil
				}
				walls = append(walls, dirPath)
				return nil
			})
		} else {
			entries, err := os.ReadDir(wallpaperDir)
			if err != nil {
				return nil, fmt.Errorf("failed to read wallpaper directory %s: %w", wallpaperDir, err)
			}

			for _, entry := range entries {
				if entry.IsDir() || !isImg(entry.Name()) {
					continue
				}
				walls = append(walls, path.Join(wallpaperDir, entry.Name()))
			}
		}
	}

	return walls, nil
}

var imageRegex = regexp.MustCompile(`^.*\.(jpeg|png)$`)

func isImg(fileName string) bool {
	return imageRegex.Match([]byte(fileName))
}
