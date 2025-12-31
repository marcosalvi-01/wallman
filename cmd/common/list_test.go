package common_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"wallman/cmd/common"
)

func createFiles(dir string, files []string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	for _, file := range files {
		if _, err := os.Create(filepath.Join(dir, file)); err != nil {
			return err
		}
	}
	return nil
}

func expectedPaths(tempDir, subdir string, files []string) []string {
	var result []string
	for _, file := range files {
		result = append(result, filepath.Join(tempDir, subdir, file))
	}
	return result
}

func TestList(t *testing.T) {
	tempDir := t.TempDir()
	tests := []struct {
		name string

		setup func(dirs []string) error

		dirs    []string
		subdirs bool
		want    []string
		wantErr bool
	}{
		{
			name: "empty dir",
			setup: func(dirs []string) error {
				return createFiles(dirs[0], []string{})
			},
			dirs:    []string{filepath.Join(tempDir, "empty")},
			subdirs: false,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "no dirs",
			setup:   nil,
			dirs:    []string{},
			subdirs: false,
			want:    []string{},
			wantErr: false,
		},
		{
			name: "one file",
			setup: func(dirs []string) error {
				return createFiles(dirs[0], []string{"test.png"})
			},
			dirs:    []string{filepath.Join(tempDir, "onefile")},
			subdirs: false,
			want:    expectedPaths(tempDir, "onefile", []string{"test.png"}),
			wantErr: false,
		},
		{
			name: "multiple files",
			setup: func(dirs []string) error {
				return createFiles(dirs[0], []string{"image1.png", "image2.jpeg", "image3.png"})
			},
			dirs:    []string{filepath.Join(tempDir, "multifile")},
			subdirs: false,
			want:    expectedPaths(tempDir, "multifile", []string{"image1.png", "image2.jpeg", "image3.png"}),
			wantErr: false,
		},
		{
			name: "non-image files ignored",
			setup: func(dirs []string) error {
				return createFiles(dirs[0], []string{"image.png", "doc.txt", "script.go", "photo.jpeg"})
			},
			dirs:    []string{filepath.Join(tempDir, "mixed")},
			subdirs: false,
			want:    expectedPaths(tempDir, "mixed", []string{"image.png", "photo.jpeg"}),
			wantErr: false,
		},
		{
			name: "subdirectories enabled",
			setup: func(dirs []string) error {
				if err := createFiles(dirs[0], []string{"root.png"}); err != nil {
					return err
				}
				return createFiles(filepath.Join(dirs[0], "sub"), []string{"sub.jpeg"})
			},
			dirs:    []string{filepath.Join(tempDir, "nested")},
			subdirs: true,
			want:    append(expectedPaths(tempDir, "nested", []string{"root.png"}), expectedPaths(tempDir, "nested/sub", []string{"sub.jpeg"})...),
			wantErr: false,
		},
		{
			name: "multiple directories",
			setup: func(dirs []string) error {
				if err := createFiles(dirs[0], []string{"a.png"}); err != nil {
					return err
				}
				return createFiles(dirs[1], []string{"b.jpeg"})
			},
			dirs:    []string{filepath.Join(tempDir, "dir1"), filepath.Join(tempDir, "dir2")},
			subdirs: false,
			want:    append(expectedPaths(tempDir, "dir1", []string{"a.png"}), expectedPaths(tempDir, "dir2", []string{"b.jpeg"})...),
			wantErr: false,
		},
		{
			name:    "directory read error",
			setup:   nil,
			dirs:    []string{filepath.Join(tempDir, "nonexistent")},
			subdirs: false,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup(tt.dirs)
				if err != nil {
					t.Errorf("Setup failed unexpectedly: %v", tt)
				}
			}
			got, gotErr := common.List(tt.dirs, tt.subdirs)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("List() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("List() succeeded unexpectedly")
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
