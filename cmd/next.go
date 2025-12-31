package cmd

import (
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Set next wallpaper",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := GetConfig()
		managerType := manager
		if managerType == "" {
			managerType = config.Manager
		}
		man, err := GetManager(config.WallpaperDirs, config.TravelSubDirs, managerType, appQueries, dryRun)
		if err != nil {
			return err
		}

		return man.Next()
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
