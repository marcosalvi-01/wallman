package hyprpaper_test

import (
	"testing"

	"wallman/hyprpaper"
)

func TestNew(t *testing.T) {
	_, err := hyprpaper.New([]string{"/home/marco/gruvbox-wallpapers/wallpapers/"}, true, nil, nil, false)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
}
