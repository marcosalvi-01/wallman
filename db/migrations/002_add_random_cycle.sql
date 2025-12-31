-- +goose Up
CREATE TABLE random_cycle (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    shuffled_wallpapers TEXT NOT NULL,
    current_index INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE random_cycle;