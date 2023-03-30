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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <host1>...<hostN>",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return deleteAction(os.Stdout, hostsFile, args)
	},
	Aliases:      []string{"d"},
	Args:         cobra.MinimumNArgs(1), //min 1 to delete
	SilenceUsage: true,
}

// deleteAction - deletes a host from file
func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	// load from file
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	// delete
	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Deleted host:", h)
	}

	// save
	return hl.Save(hostsFile)

}

func init() {
	hostsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
