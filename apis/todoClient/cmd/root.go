/*
Copyright © 2023 KaranbirS
Copyrights apply to this source code, Check LICENSE for more details.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoClient",
	Short: "A Todo API client",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// for config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todoClient.yaml)")
	rootCmd.PersistentFlags().String("api-root", "http://localhost:8080", "Todo API URL")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("TODO") // TODO_API_ROOT would match if present

	cobra.OnInitialize(initConfig)
	viper.BindPFlags(rootCmd.PersistentFlags().Lookup("api-root"))
}
func initConfig() {
	// if configfile is passed, use that
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// search for path inside home dir
		viper.AddConfigPath(home)
		viper.SetConfigName(".todoClient")
	}
	// read any environment variables if matches
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
