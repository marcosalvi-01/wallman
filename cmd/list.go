package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"wallman/cmd/common"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available wallpapers",
	Long:  `Lists all wallpapers found in the configured directories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		config := GetConfig()

		walls, err := common.List(config.WallpaperDirs, config.TravelSubDirs)
		if err != nil {
			return fmt.Errorf("failed to list wallpapers: %w", err)
		}

		if jsonOutput {
			return json.NewEncoder(os.Stdout).Encode(walls)
		}

		for _, wall := range walls {
			fmt.Println(wall)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("json", false, "Output in JSON format")
}
