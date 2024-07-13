package config

import (
	"fmt"
	"log"
	"moba-converter-go/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print some basic meta information about the loaded config.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("configPath")

		_, _, meta, err := config.LoadConfigurations(configPath)
		if err != nil {
			log.Fatalf("Error loading configurations: %v", err)
		}

		fmt.Println("Version of Value-Database")
		if version, ok := meta["version"]; ok {
			fmt.Printf("Version: %s\n", version)
		} else {
			fmt.Println("Version information not found.")
		}

		if changedWhen, ok := meta["changed_when"]; ok {
			fmt.Printf("Changed When: %s\n", changedWhen)
		} else {
			fmt.Print("Changed When information not found.")
		}
	},
}
