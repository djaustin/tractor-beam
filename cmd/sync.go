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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:        "sync SPREADSHEET_PATH REDIS_ADDRESS",
	Args:       cobra.ExactArgs(2),
	Short:      "Synchronises data from a spreadsheet into a Redis database",
	Aliases:    []string{"s"},
	SuggestFor: []string{"transfer", "upload"},
	Example:    fmt.Sprintf("%s sync ./Book1.xlsx 192.168.1.2:6379 -p password123", rootCmd.Use),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("syncing Redis database %s with spreadsheet source %s\n", args[1], args[0])
		fmt.Printf("extracting keys from %s and values from %s\n", viper.GetString("key_column"), viper.GetString("value_column"))
		fmt.Printf("using prefx %s", viper.GetString("redis_prefix"))
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringP("password", "p", "", "password used to access Redis")
	syncCmd.Flags().StringP("keycol", "k", "key", "the header of the spreadsheet column containing keys")
	syncCmd.Flags().StringP("valcol", "v", "value", "the header of the spreadsheet column containing values")
	syncCmd.Flags().String("prefix", "", "prefix attached to all keys inserted into Redis")
	viper.BindPFlag("redis_password", syncCmd.Flags().Lookup("password"))
	viper.BindPFlag("redis_address", syncCmd.Flags().Lookup("target"))
	viper.BindPFlag("key_column", syncCmd.Flags().Lookup("keycol"))
	viper.BindPFlag("value_column", syncCmd.Flags().Lookup("valcol"))
	viper.BindPFlag("redis_prefix", syncCmd.Flags().Lookup("prefix"))
}
