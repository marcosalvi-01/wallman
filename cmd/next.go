package cmd

import (
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "next wallpaper TODO",
	Run: func(cmd *cobra.Command, args []string) {
		config := GetConfig()
		man, err := GetManager(config.WallpaperDirs, config.TravelSubDirs)
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
