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

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan --ports/-p=8080,443,3000",
	Short: "Run a port scan on hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		// too lazy to do this
		_, err = cmd.Flags().GetString("ports-range")
		if err != nil {
			return err
		}

		// should we scan default ports or ports-range?
		// scan ports-range alongside default ports to be on safe side

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostsList{}
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports)
	return printResults(out, results)
}

func printResults(out io.Writer, results []scan.Results) error {
	message := ""

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			// no ports to add for printing, continue looping
			continue
		}

		message += fmt.Sprintln()

		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}

		message += fmt.Sprintln()
	}

	_, err := fmt.Fprint(out, message)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)
	// --ports=80,443
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")
	// --ports-range=1-1024
	// scanCmd.Flags().StringP("ports-range", "pr", "", "port range to scan. ex: 80-443")
}
