package cmd

import (
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Set random wallpaper",
	Long:  `Sets a random wallpaper. By default, cycles through all wallpapers without repeating until all have been used, then reshuffles. Use --true-random for completely random selection from all wallpapers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		trueRandom, _ := cmd.Flags().GetBool("true-random")
		config := GetConfig()
		managerType := manager
		if managerType == "" {
			managerType = config.Manager
		}
		man, err := GetManager(config.WallpaperDirs, config.TravelSubDirs, managerType, appQueries, dryRun)
		if err != nil {
			return err
		}

		return man.Random(trueRandom)
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().Bool("true-random", false, "Pick completely random wallpaper from all available (disables cycling)")
}
