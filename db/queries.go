package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"wallman/db/sqlc"
)

// SetWallpaper sets the current wallpaper and updates history.
func SetWallpaper(path string) error {
	q, err := Get()
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}

	ctx := context.Background()
	now := time.Now()

	// Mark previous as unset
	current, err := q.GetCurrentWallpaper(ctx)
	if err == nil {
		err = q.MarkWallpaperUnset(ctx, sqlc.MarkWallpaperUnsetParams{
			UnsetAt: &now,
			Path:    current.Path,
		})
		if err != nil {
			return fmt.Errorf("failed to mark previous unset: %w", err)
		}
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to get current: %w", err)
	}

	// Insert new history
	err = q.InsertWallpaperHistory(ctx, sqlc.InsertWallpaperHistoryParams{
		Path:  path,
		SetAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	// Update current
	_, err = q.UpdateCurrentWallpaper(ctx, sqlc.UpdateCurrentWallpaperParams{
		Path:  path,
		SetAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to update current: %w", err)
	}

	return nil
}

// GetCurrentWallpaperPath returns the current wallpaper path.
func GetCurrentWallpaperPath() (string, error) {
	q, err := Get()
	if err != nil {
		return "", fmt.Errorf("error getting db connection: %w", err)
	}

	current, err := q.GetCurrentWallpaper(context.Background())
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no current wallpaper set")
	}
	if err != nil {
		return "", fmt.Errorf("error getting current wallpaper: %w", err)
	}
	return current.Path, nil
}

// GetWallpaperHistory returns the wallpaper history.
func GetWallpaperHistory(limit int) ([]sqlc.WallpaperHistory, error) {
	q, err := Get()
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}

	history, err := q.GetWallpaperHistory(context.Background(), sqlc.GetWallpaperHistoryParams{
		Column1: nil,
		SetAt:   time.Time{},
		Column3: nil,
		ID:      0,
		Limit:   int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting wallpaper history: %w", err)
	}
	return history, nil
}

// GetPreviousWallpaper returns the previous wallpaper path.
func GetPreviousWallpaper() (string, error) {
	q, err := Get()
	if err != nil {
		return "", fmt.Errorf("error getting db connection: %w", err)
	}

	prev, err := q.GetPreviousWallpaper(context.Background())
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no previous wallpaper")
	}
	if err != nil {
		return "", fmt.Errorf("error getting previous wallpaper: %w", err)
	}
	return prev.Path, nil
}

// GetRandomCycle returns the current random cycle state.
func GetRandomCycle() (shuffled []string, index int, err error) {
	q, err := Get()
	if err != nil {
		return nil, 0, fmt.Errorf("error getting db connection: %w", err)
	}

	cycle, err := q.GetRandomCycle(context.Background())
	if err == sql.ErrNoRows {
		return nil, 0, fmt.Errorf("no random cycle set")
	}
	if err != nil {
		return nil, 0, fmt.Errorf("error getting random cycle: %w", err)
	}

	err = json.Unmarshal([]byte(cycle.ShuffledWallpapers), &shuffled)
	if err != nil {
		return nil, 0, fmt.Errorf("error unmarshaling shuffled wallpapers: %w", err)
	}
	return shuffled, int(cycle.CurrentIndex), nil
}

// UpsertRandomCycle updates the random cycle state.
func UpsertRandomCycle(shuffled []string, index int) error {
	q, err := Get()
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}

	data, err := json.Marshal(shuffled)
	if err != nil {
		return fmt.Errorf("error marshaling shuffled wallpapers: %w", err)
	}

	err = q.UpsertRandomCycle(context.Background(), sqlc.UpsertRandomCycleParams{
		ShuffledWallpapers: string(data),
		CurrentIndex:       int64(index),
	})
	if err != nil {
		return fmt.Errorf("error upserting random cycle: %w", err)
	}
	return nil
}
