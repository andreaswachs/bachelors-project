/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/andreaswachs/bachelors-project/daaukins/client/api"
	"github.com/spf13/cobra"
)

var (
	serverId string
	all      bool = true
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [labs/servers/frontends/identifier]",
	Short: "Get information about a lab or all labs",
	Long: `If you specify a lab identifier, you will get information about that lab.
	Any other argument will get you information about all <resource>.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "labs":
			labs(cmd, args)
		case "servers":
			servers(cmd, args)
		case "frontends":
			frontends(cmd, args)
		default:
			lab(cmd, args)
		}
	},
}

func frontends(cmd *cobra.Command, args []string) {
	response, err := api.GetFrontends()
	if err != nil {
		cmd.Printf("There were an error when trying to get all frontends:\n%s\n", err)
		return
	}

	printer.AddHeader("HOST", "PORT", "SERVER", "LINK")
	for _, frontend := range response.Frontends {
		printer.AddRow(frontend.Host,
			frontend.Port,
			frontend.ServerId,
			fmt.Sprintf("https://%s:%s/vnc.html", frontend.Host, frontend.Port))
	}
}

func lab(cmd *cobra.Command, args []string) {
	response, err := api.GetLab(args[0])
	if err != nil {
		cmd.Printf("There were an error when trying to get the lab:\n%s\n", err)
		return
	}

	ppLabHeader(cmd)
	printer.AddRow(response.Lab.Id, response.Lab.Name, response.Lab.NumChallenges, response.Lab.NumUsers, response.Lab.ServerId)
}

func labs(cmd *cobra.Command, args []string) {
	response, err := api.GetLabs(serverId)
	if err != nil {
		cmd.Println("There were an error when trying to get all labs: ")
		return
	}

	ppLabHeader(cmd)
	for _, lab := range response.Labs {
		printer.AddRow(lab.Id, lab.Name, lab.NumChallenges, lab.NumUsers, lab.ServerId)
	}
}

func servers(cmd *cobra.Command, args []string) {
	response, err := api.GetServers()
	if err != nil {
		cmd.Println("There were an error when trying to get all servers: ")
		return
	}

	printer.AddHeader("ID", "NAME", "MODE", "LABS", "CONNECTED")
	for _, server := range response.Servers {
		printer.AddRow(server.Id, server.Name, server.Mode, server.NumLabs, server.Connected)
	}
}

func ppLabHeader(cmd *cobra.Command) {
	printer.AddHeader("ID", "NAME", "CHALLENGES", "USERS", "SERVERID")
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&serverId, "id", "i", "", "specifies a server ID for the given resource")
}
