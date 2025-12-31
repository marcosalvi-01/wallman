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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func defaultConfig() *Config {
	return &Config{
		WallpaperDirs: []string{},
		TravelSubDirs: false,
	}
}

func initConfig() {
	var config *Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		config = defaultConfig()
		appConfig = config
		return
	}

	configPath := findConfigPath(cfgFile, homeDir)
	if configPath != "" {
		loaded, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config from %s: %v\nUsing defaults.\n", configPath, err)
			config = defaultConfig()
		} else {
			config = loaded
		}
	} else {
		config = defaultConfig()
	}

	// Expand paths
	expandedDirs := make([]string, len(config.WallpaperDirs))
	for i, dir := range config.WallpaperDirs {
		expandedDirs[i] = expandPath(dir)
	}
	config.WallpaperDirs = expandedDirs

	// Validate
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Config validation failed: %v\n", err)
		os.Exit(1)
	}

	appConfig = config
}

func findConfigPath(cfgFile, homeDir string) string {
	if cfgFile != "" {
		return cfgFile
	}

	possiblePaths := []string{
		filepath.Join(homeDir, ".config", "wallman.yaml"),
		filepath.Join(homeDir, ".config", "wallman.yml"),
		filepath.Join(homeDir, ".config", "wallman", "wallman.yaml"),
		filepath.Join(homeDir, ".config", "wallman", "wallman.yml"),
		filepath.Join(homeDir, ".config", "wallman", "config.yaml"),
		filepath.Join(homeDir, ".config", "wallman", "config.yml"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func validateConfig(config *Config) error {
	for _, dir := range config.WallpaperDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("wallpaper directory does not exist: %s", dir)
		}
	}
	return nil
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
