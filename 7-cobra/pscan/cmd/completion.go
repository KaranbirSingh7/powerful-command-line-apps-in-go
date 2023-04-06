/*
Copyright © 2023 KaranbirS
Copyrights apply to this source code, Check LICENSE for more details.
*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate bash completion for your command",
	Long: `To load your completions run
source <(pScan completion)

To load completion automatically on login, add this line to your .bashrc file:
$ ~/.bashrc
source <(pScan completion)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completionAction(os.Stdout)
	},
}

func completionAction(out io.Writer) error {
	return rootCmd.GenZshCompletion(out)

}

func init() {
	rootCmd.AddCommand(completionCmd)
}
