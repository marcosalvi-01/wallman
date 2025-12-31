package cmd

import (
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Set random wallpaper",
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

		err = man.Random()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}
