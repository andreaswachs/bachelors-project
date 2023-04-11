/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/client/api"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove single labs by id, or labs either by server or all",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			if serverId == "unset" && !all {
				cmd.Println("You need to specify a server id or use the --all flag")
				cmd.Help()
				os.Exit(1)
			}

			if all {
				serverId = ""
			}

			response, err := api.RemoveLabs(serverId)
			if err != nil {
				cmd.Println("Error removing all labs:", err)
				os.Exit(1)
			}

			cmd.Println("All labs removed successfully:", response.Ok)
			return
		case 1:
			response, err := api.RemoveLab(args[0])
			if err != nil {
				cmd.Println("Error removing lab:", err)
				os.Exit(1)
			}

			cmd.Println("Lab removed successfully:", response.Ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&all, "all", "a", false, "Remove all labs")
	removeCmd.Flags().StringVarP(&serverId, "serverid", "s", "unset", "Remove all labs from a specific server")
}
