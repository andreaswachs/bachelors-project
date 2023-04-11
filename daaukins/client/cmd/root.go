/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/andreaswachs/bachelors-project/daaukins/client/api"
	"github.com/andreaswachs/bachelors-project/daaukins/client/config"
	"github.com/spf13/cobra"
)

var (
	FlagsForce bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dkn",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensures that for all commands except "config init" the config file is loaded
		// before the command is executed.
		// If the user inputted "config init", the config file is attempted to be created
		if strings.HasPrefix(cmd.Use, "config") && len(args) == 1 && args[0] == "init" {
			config.Initialize(FlagsForce)
			os.Exit(0)
		} else {
			if err := config.Load(); err != nil {
				cmd.Println("Error loading config file:", err)
				os.Exit(1)
			}

			// We don't want to wait for the connection to be established,
			// cobra can do some processing while we connect in the background
			if !strings.HasPrefix(cmd.Use, "config") {
				go api.Initialize(config.ServerAddress(), config.ServerPort())
			}
		}
	},
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
	rootCmd.PersistentFlags().BoolVar(&FlagsForce, "force", false, "Force initialization even if config file already exists")
}
