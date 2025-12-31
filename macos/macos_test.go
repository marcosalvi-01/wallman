package macos_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcosalvi-01/wallman/macos"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	wallpaperDir := filepath.Join(tempDir, "wallpapers")

	err := os.MkdirAll(wallpaperDir, 0o750)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create some test image files
	files := []string{"test.png", "image.jpeg", "test.jpg", "image.bmp", "test.webp", "notimage.txt"}
	for _, file := range files {
		path := filepath.Join(wallpaperDir, file)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
		err = f.Close()
		if err != nil {
			t.Fatalf("Failed to close file %s: %v", file, err)
		}
	}

	// Test New
	m, err := macos.New([]string{wallpaperDir}, false, nil, true) // dryRun true, queries nil for simplicity
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if m == nil {
		t.Fatal("New() returned nil")
	}
}
