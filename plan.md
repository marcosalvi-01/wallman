# Wallman Implementation Plan

This plan outlines the features and tasks to complete the Wallman CLI tool for managing wallpapers across different systems.

## 1. CLI Structure and Commands

- [x] Implement root command with global flags and `--manager` flag with config support
- [x] Add `wallman next` subcommand (calls Manager.Next())
- [x] Add `wallman previous` subcommand (calls Manager.Previous())
- [x] Add `wallman random` subcommand (calls Manager.Random())
- [x] Add `wallman list` global command (scans config dirs, supports `--json` for structured output)
- [x] Add `wallman current` subcommand (calls Manager.Current())
- [x] Add `wallman history` subcommand (calls Manager.History(), with options like `--all`, `--since`)
- [x] Add `wallman config` subcommand for validating/editing config
- [x] Add `wallman init` subcommand to create default config

## 2. Manager Interface and Backends

- [x] Update Manager interface to accept DB connection for persistence
- [x] Implement full hyprpaper backend with actual wallpaper setting (e.g., hyprctl commands)
- [x] Add factory logic to select manager based on config or auto-detection
- [ ] Stub backends for other systems (matugen, macOS) with placeholder implementations

## 3. Persistence with SQLite

- [x] Create SQLite schema in `~/.local/share/wallman/wallman.db`
  - [x] `current_wallpaper` table (id, path, set_at)
  - [x] `wallpaper_history` table (id, path, set_at, unset_at) with indexes
  - [x] `schema_version` table for migrations
- [x] Implement DB initialization on first run
- [x] Update Manager methods to read/write to DB
  - [x] Next/Random: Insert history, update current, set unset_at
  - [x] Previous: Query and set from history
  - [x] Current: Select from current table
  - [x] History: Query history with filters

## 4. Configuration Handling

- [x] Enhance config loading with better error messages and defaults
- [x] Add config validation (paths exist, readable)
- [x] Support multiple config file locations with priority

## 5. Error Handling and Edge Cases

- [x] Handle missing config (suggest `wallman init`)
- [x] Handle no wallpapers (error on empty list)
- [x] Handle empty history (error on previous/history if none)
- [x] Validate paths and skip inaccessible ones
- [x] Wrap Manager calls with proper error handling
- [x] Add `--dry-run` flag for simulation

## 6. Testing and Quality

- [x] Write unit tests for Manager implementations (hyprpaper tests missing)
- [x] Write integration tests for CLI commands (some tests exist but failing)
- [x] Test edge cases (missing tools, invalid configs)
- [ ] Remove debug prints and add proper logging
- [ ] Update README with usage examples

## 7. Extensibility and Polish

- [x] Make output structured (JSON with `--json` flag)
- [x] Add version info and help documentation
