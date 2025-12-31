# Wallman

Frontend for different wallpaper managers & stuff

# How

```go
type Manager interface {
	Next() error
	Previous() error
	Random() error
	Current() (string, error)
	History() ([]string, error)
	Set(path string) error
}
```

Implement the interface to interact with the cli commands.

Each implementation should have the commands:

| Command    | Description                                                                                                                                            |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Next`     | Set the next wallpaper as the active. "Next" is defined by the alphabetical order in the directories                                                   |
| `Previous` | It will access the `History` of set wallpapers for this Manager and set the active to the last one                                                     |
| `Random`   | It will take a random wallpaper from the possible wallpapers and set it as active                                                                      |
| `Current`  | Return the currently set wallpaper for this Manager (its path i think?)                                                                                |
| `History`  | Will return the history of the wallpapers set for this Manager                                                                                         |
| `Set`      | Set a specific wallpaper as active, even if it's not in the configured directories                                                                     |

Additionally, there is a global `list` command that lists all available wallpapers from the configured directories.
