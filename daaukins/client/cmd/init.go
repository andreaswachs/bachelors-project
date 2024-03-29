/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/andreaswachs/bachelors-project/daaukins/client/config"
	"github.com/spf13/cobra"
)

var (
	force bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wasWritten, err := config.CreateConfigFile(force)
		if err != nil {
			panic(err)
		}

		if !wasWritten {
			cmd.Println("Config file already exists. Use --force to overwrite.")
		} else {
			cmd.Println("Config file created at", config.DknConfigFile())
			cmd.Println("Please edit the file and add your Daaukins server address and port")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force initialization even if config file already exists")
}
