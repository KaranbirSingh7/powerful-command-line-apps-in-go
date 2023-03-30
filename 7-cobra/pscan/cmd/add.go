/*
Copyright © 2023 KaranbirS
Copyrights apply to this source code, Check LICENSE for more details.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/karanbirsingh7/pscan/scan"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1>...<hostN>",
	Short:        "Add new host(s) to the list",
	Args:         cobra.MinimumNArgs(1), //TIP: this is neat feature
	Aliases:      []string{"a"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		fmt.Println(hostsFile)
		if err != nil {
			return err
		}
		return addActions(os.Stdout, hostsFile, args)
	},
}

func addActions(output io.Writer, hostsFile string, args []string) error {
	hl := scan.HostsList{}

	// load file
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	// add to list
	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}
		fmt.Fprintln(output, "Added host:", h)
	}

	// save file
	return hl.Save(hostsFile)

}

func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
