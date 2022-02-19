/*
Copyright © 2022 Dan Austin <dan.austin@hey.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tractor-beam",
	Short: "A CLI for synchronising a Redis key-value store with an Excel spreadsheet",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tractor-beam.yaml)")
	rootCmd.PersistentFlags().StringP("password", "p", "", "password used to access Redis")
	rootCmd.PersistentFlags().StringP("keycol", "k", "key", "the header of the spreadsheet column containing keys")
	rootCmd.PersistentFlags().StringP("valcol", "v", "value", "the header of the spreadsheet column containing values")
	rootCmd.PersistentFlags().StringP("sheet", "s", "Sheet1", "the name of the worksheet containing data for sync")
	rootCmd.PersistentFlags().String("prefix", "", "prefix attached to all keys inserted into Redis")

	viper.BindPFlag("redis_password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("key_column", rootCmd.PersistentFlags().Lookup("keycol"))
	viper.BindPFlag("value_column", rootCmd.PersistentFlags().Lookup("valcol"))
	viper.BindPFlag("worksheet", rootCmd.PersistentFlags().Lookup("sheet"))
	viper.BindPFlag("redis_prefix", rootCmd.PersistentFlags().Lookup("prefix"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tractor-beam" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tractor-beam")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
