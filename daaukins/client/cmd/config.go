/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/andreaswachs/bachelors-project/daaukins/client/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [action]",
	Short: "Configure the dkn tool",
	Long: `Interact with the configuration
	
Action:
  init		Initialize the config file
  show		Show the current config
  path		Show the path to the config file`,
	Annotations: map[string]string{
		"config": "true",
	},
	ValidArgs: []string{"init", "show", "path"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			switch args[0] {
			case "show":
				cmd.Println("Config file at", config.DknConfigFile())
				cmd.Printf("Server address: %s\n", config.ServerAddress())
				cmd.Printf("Server port: %s\n", config.ServerPort())
			case "path":
				cmd.Println(config.DknConfigFile())
			case "init":
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
