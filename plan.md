# Wallman Implementation Plan

This plan outlines the features and tasks to complete the Wallman CLI tool for managing wallpapers across different systems.

## 1. CLI Structure and Commands

- [ ] Implement root command with global flags (`--config`, `--verbose`, `--manager`)
- [ ] Add `wallman next` subcommand (calls Manager.Next())
- [ ] Add `wallman previous` subcommand (calls Manager.Previous())
- [ ] Add `wallman random` subcommand (calls Manager.Random())
- [x] Add `wallman list` global command (scans config dirs, supports `--json` for structured output)
- [ ] Add `wallman current` subcommand (calls Manager.Current())
- [ ] Add `wallman history` subcommand (calls Manager.History(), with options like `--all`, `--since`)
- [ ] Add `wallman config` subcommand for validating/editing config
- [ ] Add `wallman init` subcommand to create default config

## 2. Manager Interface and Backends

- [ ] Update Manager interface to accept DB connection for persistence
- [ ] Implement full hyprpaper backend with actual wallpaper setting (e.g., hyprctl commands)
- [ ] Add factory logic to select manager based on config or auto-detection
- [ ] Stub backends for other systems (matugen, macOS) with placeholder implementations

## 3. Persistence with SQLite

- [ ] Create SQLite schema in `~/.local/share/wallman/wallman.db`
  - [ ] `current_wallpaper` table (id, path, set_at)
  - [ ] `wallpaper_history` table (id, path, set_at, unset_at) with indexes
  - [ ] `schema_version` table for migrations
- [ ] Implement DB initialization on first run
- [ ] Update Manager methods to read/write to DB
  - [ ] Next/Random: Insert history, update current, set unset_at
  - [ ] Previous: Query and set from history
  - [ ] Current: Select from current table
  - [ ] History: Query history with filters
- [ ] Add DB vacuuming and corruption recovery

## 4. Configuration Handling

- [ ] Enhance config loading with better error messages and defaults
- [ ] Add config validation (paths exist, readable)
- [ ] Support multiple config file locations with priority

## 5. Error Handling and Edge Cases

- [ ] Handle missing config (suggest `wallman init`)
- [ ] Handle no wallpapers (error on empty list)
- [ ] Handle empty history (error on previous/history if none)
- [ ] Validate paths and skip inaccessible ones
- [ ] Wrap Manager calls with proper error handling
- [ ] Add `--dry-run` flag for simulation
- [ ] Implement file locking for concurrent DB access

## 6. Testing and Quality

- [ ] Write unit tests for Manager implementations
- [ ] Write integration tests for CLI commands
- [ ] Test edge cases (missing tools, invalid configs)
- [ ] Remove debug prints and add proper logging
- [ ] Update README with usage examples

## 7. Extensibility and Polish

- [ ] Make output structured (JSON with `--json` flag)
- [ ] Add progress indicators for long operations
- [ ] Add version info and help documentation
