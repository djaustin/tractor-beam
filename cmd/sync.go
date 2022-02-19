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

	"github.com/djaustin/tractor-beam/db"
	l "github.com/djaustin/tractor-beam/logger"
	"github.com/go-redis/redis/v8"
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
		spreadsheetPath, redisAddress := args[0], args[1]
		l.Logger.Info("starting synchronisation of database from file")
		count, err := db.SyncDatabase(cmd.Context(),
			redis.NewClient(&redis.Options{Addr: redisAddress}),
			viper.GetString("redis_prefix"),
			spreadsheetPath,
			viper.GetString("worksheet"),
			viper.GetString("key_column"),
			viper.GetString("value_column"),
		)
		if err != nil {
			l.Logger.Fatalf("unable to synchronise database from file: %v", err)
		}
		l.Logger.Infof("%d values synchronised from '%s' to '%s'", count, spreadsheetPath, redisAddress)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
