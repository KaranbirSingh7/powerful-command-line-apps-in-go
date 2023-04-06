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
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List hosts in hosts file",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		return listAction(os.Stdout, hostsFile, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func listAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	// iteratate over each row and print it out
	for _, h := range hl.Hosts {
		if _, err := fmt.Fprintln(out, h); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
