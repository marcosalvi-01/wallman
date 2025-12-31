package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wallman config",
	Long:  `Creates a default config file if it doesn't exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var configPath string

		if cfgFile != "" {
			configPath = cfgFile
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("error getting home directory: %w", err)
			}
			configPath = filepath.Join(homeDir, ".config", "wallman.yaml")
		}

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Config file already exists at %s\n", configPath)
			return nil
		}

		// Create directory
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		// Create default config
		config := defaultConfig()
		data, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("error marshaling config: %w", err)
		}

		// Write config file
		if err := os.WriteFile(configPath, data, 0o600); err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}

		fmt.Printf("Created default config at %s\n", configPath)
		return nil
	},
}
