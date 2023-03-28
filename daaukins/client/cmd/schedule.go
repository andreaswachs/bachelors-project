/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andreaswachs/bachelors-project/daaukins/client/api"
	"github.com/spf13/cobra"
)

var (
	file string
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a lab",
	Long:  `Schedule a lab on a Daaukins server (or one of its followers)`,
	Run: func(cmd *cobra.Command, args []string) {
		if file == "" {
			cmd.Help()
			os.Exit(1)
		}

		labConfig := getLabsFileContent()

		response, err := api.ScheduleLab(labConfig)
		if err != nil {
			fmt.Println("Error scheduling lab:", err)
			os.Exit(1)
		}

		fmt.Println("Lab scheduled successfully with id:", response.Id)
	},
}

func getLabsFileContent() string {
	if file == "-" {
		reader := bufio.NewReader(os.Stdin)
		labSb := strings.Builder{}
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			_, err = labSb.WriteString(line)
			if err != nil {
				fmt.Println("Error writing to string builder")
				os.Exit(1)
			}
		}

		return labSb.String()
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist:", file)
		os.Exit(1)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", file)
		os.Exit(1)
	}

	return string(data)
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scheduleCmd.Flags().StringVar(&file, "file", "", "Specify the lab configuration file")
}
