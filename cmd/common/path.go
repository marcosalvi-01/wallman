package common

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands ~ to home directory and environment variables in a path
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	return os.ExpandEnv(path)
}
