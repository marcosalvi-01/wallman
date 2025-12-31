-- +goose Up
CREATE TABLE current_wallpaper (
    id INTEGER PRIMARY KEY,
    path TEXT NOT NULL,
    set_at DATETIME NOT NULL
);

CREATE TABLE wallpaper_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    set_at DATETIME NOT NULL,
    unset_at DATETIME
);

CREATE INDEX idx_wallpaper_history_set_at ON wallpaper_history (set_at);
CREATE INDEX idx_wallpaper_history_unset_at ON wallpaper_history (unset_at);

-- +goose Down
DROP INDEX idx_wallpaper_history_unset_at;
DROP INDEX idx_wallpaper_history_set_at;
DROP TABLE wallpaper_history;
DROP TABLE current_wallpaper;
