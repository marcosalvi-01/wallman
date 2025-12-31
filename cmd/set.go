package cmd

import (
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <path>",
	Short: "Set specific wallpaper",
	Long:  `Sets a specific wallpaper file as active, even if it's not in the configured directories. The path must be to a valid JPEG or PNG file.`,
	Args:  cobra.ExactArgs(1),
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

		return man.Set(args[0])
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
