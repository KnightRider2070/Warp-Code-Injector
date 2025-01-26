package cmd

import (
	"fmt"
	"wci/config"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays application information",
	Long:  "Displays detailed information about the application, including name, version, version number, and author.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Application Information:")
		fmt.Println("-------------------------")
		fmt.Printf("Name          : %s\n", config.AppName)
		fmt.Printf("Version       : %s\n", config.Version)
		fmt.Printf("Version Number: %s\n", config.VersionNumber)
		fmt.Println("Author        : KnightRider2070")
		fmt.Println("GitHub        : https://github.com/KnightRider2070/Warp-Code-Injector")
		fmt.Println("-------------------------")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
