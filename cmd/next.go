package cmd

import (
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Set next wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		config := GetConfig()
		managerType := manager
		if managerType == "" {
			managerType = config.Manager
		}
		man, err := GetManager(config.WallpaperDirs, config.TravelSubDirs, managerType, appQueries, dryRun)
		if err != nil {
			panic(err)
		}

		err = man.Next()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
