package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "TODO",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var configPath string

		if cfgFile != "" {
			configPath = cfgFile
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
				return fmt.Errorf("TODO: %w", err)
			}

			possiblePaths := []string{
				filepath.Join(homeDir, ".config", "wallman.yaml"),
				filepath.Join(homeDir, ".config", "wallman.yml"),
				filepath.Join(homeDir, ".config", "wallman", "wallman.yaml"),
				filepath.Join(homeDir, ".config", "wallman", "wallman.yml"),
				filepath.Join(homeDir, ".config", "wallman", "config.yaml"),
				filepath.Join(homeDir, ".config", "wallman", "config.yml"),
			}

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					configPath = path
					break
				}
			}
		}

		if configPath == "" {
			fmt.Fprintf(os.Stderr, "No config file found. Please create one at ~/.wallman.yaml\n")
			return fmt.Errorf("TODO")
		}

		return nil
	},
}
