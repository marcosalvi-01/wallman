package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current wallpaper",
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

		path, err := man.Current()
		if err != nil {
			panic(err)
		}

		fmt.Println(path)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
