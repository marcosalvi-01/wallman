package cmd

import (
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version = "dev"

var (
	manager string
	dryRun  bool
)

func getVersion() string {
	if Version != "" && Version != "dev" {
		return Version
	}

	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wallman",
	Short: "Wallpaper manager",
	Long:  `A command-line wallpaper manager that supports cycling through wallpapers, setting random wallpapers, viewing current wallpaper and history across desktop environments like Hyprland.`,

	// If no subcommand is provided, show help.
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
	Version: getVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVar(&manager, "manager", "", "force specific manager type")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "simulate actions without changing system")

	rootCmd.AddCommand(initCmd)

	cobra.OnInitialize(initConfig)
}
