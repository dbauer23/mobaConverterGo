package config

import (
	"fmt"
	"log"
	"moba-converter-go/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all possible parameters for the input file.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("configPath")

		ops, _, _, err := config.LoadConfigurations(configPath)
		if err != nil {
			log.Fatalf("Error loading configurations: %v", err)
		}

		for key, option := range ops {
			fmt.Printf("Option: %s\n", key)
			fmt.Printf("  Section: %s\n", option.Section)
			fmt.Printf("  Default Value: %s\n", option.Default)
			fmt.Printf("  Help: %s\n", option.Help)
			if len(option.Options) > 0 {
				fmt.Println("  Possible Values:")
				for optKey := range option.Options {
					fmt.Printf("    - %s\n", optKey)
				}
			}
			fmt.Println()
		}
	},
}
