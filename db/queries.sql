-- name: InsertWallpaperHistory :exec
INSERT INTO
    wallpaper_history (path, set_at)
VALUES
    (?, ?);

-- name: UpdateCurrentWallpaper :one
INSERT
    OR REPLACE INTO current_wallpaper (id, path, set_at)
VALUES
    (1, ?, ?)
RETURNING
    path,
    set_at;

-- name: GetCurrentWallpaper :one
SELECT
    path,
    set_at
FROM
    current_wallpaper
WHERE
    id = 1;

-- name: GetWallpaperHistory :many
SELECT
    id,
    path,
    set_at,
    unset_at
FROM
    wallpaper_history
WHERE
    (
        ? IS NULL
        OR set_at >= ?
    )
    AND (
        ? IS NULL
        OR id <= ?
    )
ORDER BY
    set_at DESC
LIMIT
    ?;

-- name: GetPreviousWallpaper :one
SELECT
    id,
    path,
    set_at
FROM
    wallpaper_history
WHERE
    unset_at IS NOT NULL
ORDER BY
    unset_at DESC
LIMIT
    1;

-- name: MarkWallpaperUnset :exec
UPDATE
    wallpaper_history
SET
    unset_at = ?
WHERE
    path = ?
    AND unset_at IS NULL;

-- name: GetRandomCycle :one
SELECT shuffled_wallpapers, current_index FROM random_cycle WHERE id = 1;

-- name: UpsertRandomCycle :exec
INSERT OR REPLACE INTO random_cycle (id, shuffled_wallpapers, current_index) VALUES (1, ?, ?);
