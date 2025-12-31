package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current config",
	Long:  `Displays the current configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		config := GetConfig()

		if jsonOutput {
			return json.NewEncoder(os.Stdout).Encode(config)
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("error marshaling config: %w", err)
		}

		fmt.Print(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().Bool("json", false, "Output in JSON format")
}
