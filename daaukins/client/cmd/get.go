/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/andreaswachs/bachelors-project/daaukins/client/api"
	"github.com/andreaswachs/bachelors-project/daaukins/service"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [labs/identifier]",
	Short: "Get information about a lab or all labs",
	Long: `If you specify a lab identifier, you will get information about that lab.
	If you get "labs" as an argument, you will get a list of all labs.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "labs":
			labs(cmd, args)
		default:
			lab(cmd, args)
		}
	},
}

func lab(cmd *cobra.Command, args []string) {
	response, err := api.GetLab(args[0])
	if err != nil {
		cmd.Println("There were an error when trying to get the lab: ")
	}

	ppLabHeader(cmd)
	ppLab(cmd, response.Lab)
}

func labs(cmd *cobra.Command, args []string) {
	response, err := api.GetLabs()
	if err != nil {
		cmd.Println("There were an error when trying to get all labs: ")
	}

	ppLabHeader(cmd)
	for _, lab := range response.Labs {
		ppLab(cmd, lab)
	}
}

func ppLabHeader(cmd *cobra.Command) {
	cmd.Printf("Id\t\t\t\tName\t\tChallenges\tUsers\tServerId\n")
}

func ppLab(cmd *cobra.Command, lab *service.LabDescription) {
	cmd.Printf("%s\t%s\t%d\t\t%d\t%s\n",
		lab.Id,
		lab.Name,
		lab.NumChallenges,
		lab.NumUsers,
		lab.ServerId,
	)
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
