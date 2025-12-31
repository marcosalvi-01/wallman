package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show wallpaper history",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		config := GetConfig()
		managerType := manager
		if managerType == "" {
			managerType = config.Manager
		}
		man, err := GetManager(config.WallpaperDirs, config.TravelSubDirs, managerType, appQueries, dryRun)
		if err != nil {
			return err
		}

		paths, err := man.History()
		if err != nil {
			return err
		}

		if jsonOutput {
			return json.NewEncoder(os.Stdout).Encode(paths)
		}

		for _, path := range paths {
			fmt.Println(path)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().Bool("json", false, "Output in JSON format")
}
