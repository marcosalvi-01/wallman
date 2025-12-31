package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"tilde expansion", "~/test", filepath.Join(home, "test")},
		{"env expansion", "$HOME/test", filepath.Join(home, "test")},
		{"no expansion", "/absolute/path", "/absolute/path"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPath(tt.input)
			if got != tt.want {
				t.Errorf("expandPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tempDir := t.TempDir()
	nonExistent := filepath.Join(tempDir, "nonexistent")
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"valid dirs", &Config{WallpaperDirs: []string{tempDir}}, false},
		{"non-existent dir", &Config{WallpaperDirs: []string{nonExistent}}, true},
		{"empty dirs", &Config{WallpaperDirs: []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindConfigPath(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config")
	os.MkdirAll(configDir, 0o755)
	configFile := filepath.Join(configDir, "wallman.yaml")
	os.WriteFile(configFile, []byte("test"), 0o644)

	customFile := filepath.Join(tempDir, "custom.yaml")

	tests := []struct {
		name    string
		cfgFile string
		homeDir string
		want    string
	}{
		{"custom path", customFile, tempDir, customFile},
		{"found file", "", tempDir, configFile},
		{"no file", "", tempDir, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "no file" {
				os.Remove(configFile)
			} else {
				os.WriteFile(configFile, []byte("test"), 0o644)
			}
			got := findConfigPath(tt.cfgFile, tt.homeDir)
			if got != tt.want {
				t.Errorf("findConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfgFile   string
		hasConfig bool
		want      *Config
	}{
		{"valid config", "", true, nil}, // will set want dynamically
		{"no config", "", false, defaultConfig()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			t.Setenv("HOME", tempDir)

			configDir := filepath.Join(tempDir, ".config")
			os.MkdirAll(configDir, 0o755)
			configFile := filepath.Join(configDir, "wallman.yaml")

			if tt.hasConfig {
				validYaml := `wallpaper_directories:
  - "` + tempDir + `"
travel_sub_directories: true`
				os.WriteFile(configFile, []byte(validYaml), 0o644)
				tt.want = &Config{WallpaperDirs: []string{tempDir}, TravelSubDirs: true}
			} else {
				os.Remove(configFile)
			}

			oldCfgFile := cfgFile
			oldAppConfig := appConfig
			cfgFile = tt.cfgFile
			appConfig = nil

			initConfig()

			cfgFile = oldCfgFile
			defer func() { appConfig = oldAppConfig }()

			if !reflect.DeepEqual(appConfig, tt.want) {
				t.Errorf("initConfig() set appConfig = %v, want %v", appConfig, tt.want)
			}
		})
	}
}
