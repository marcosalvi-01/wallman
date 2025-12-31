package hyprpaper_test

import (
	"testing"

	"wallman/hyprpaper"
)

func TestNew(t *testing.T) {
	out, err := hyprpaper.New([]string{"/home/marco/gruvbox-wallpapers/wallpapers/"}, true)
	if err != nil {
		panic(err)
	}
}
