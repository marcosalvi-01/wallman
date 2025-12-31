package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	appConfig *Config
	cfgFile   string
)

type Config struct {
	WallpaperDirs []string `yaml:"wallpaper_directories"`
	TravelSubDirs bool     `yaml:"travel_sub_directories"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("TODO: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	expandedDirs := make([]string, len(config.WallpaperDirs))
	for i, dir := range config.WallpaperDirs {
		expandedDirs[i] = expandPath(dir)
	}

	config.WallpaperDirs = expandedDirs

	return &config, err
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	return appConfig
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	return os.ExpandEnv(path)
}
