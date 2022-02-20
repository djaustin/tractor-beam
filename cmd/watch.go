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
	"os"
	"os/signal"
	"syscall"

	"github.com/djaustin/tractor-beam/db"
	l "github.com/djaustin/tractor-beam/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:        "watch SPREADSHEET_PATH REDIS_ADDRESS",
	Args:       cobra.ExactArgs(2),
	Short:      "Watches a spreadsheet file for changes and synchronises the Redis database on change",
	Aliases:    []string{"w"},
	SuggestFor: []string{"monitor", "track"},
	Example:    fmt.Sprintf("%s watch ./Book1.xlsx 192.168.1.2:6379 -p password123", rootCmd.Use),
	Run: func(cmd *cobra.Command, args []string) {
		spreadsheetPath, redisAddress := args[0], args[1]
		// Run initial sync before watching
		l.Logger.Info("starting initial synchronisation of database from file")
		redisClient := redis.NewClient(&redis.Options{Addr: redisAddress})
		l.Logger.Debugf("new Redis client created for '%s'", redisAddress)

		updateCount, err := db.SyncDatabase(cmd.Context(),
			redisClient,
			viper.GetString("redis_prefix"),
			spreadsheetPath,
			viper.GetString("worksheet"),
			viper.GetString("key_column"),
			viper.GetString("value_column"),
		)
		if err != nil {
			l.Logger.Fatalf("unable to synchronise %s with data from file %s: %v", redisAddress, spreadsheetPath, err)
		}
		l.Logger.Info("initial synchronisation completed")
		l.Logger.Infof("%d values synchronised from %s to %s", updateCount, spreadsheetPath, redisAddress)

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			l.Logger.Fatalf("failed to create new file watcher: %v", err)
		}
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op != fsnotify.Write && event.Op != fsnotify.Create {
						continue
					}
					l.Logger.Debugf("detected %s on watched file '%s'", event.Op.String(), event.Name)
					l.Logger.Infof("starting synchronisation of data from file to %v", redisAddress)
					updateCount, err := db.SyncDatabase(cmd.Context(),
						redis.NewClient(&redis.Options{Addr: redisAddress}),
						viper.GetString("redis_prefix"),
						spreadsheetPath,
						viper.GetString("worksheet"),
						viper.GetString("key_column"),
						viper.GetString("value_column"),
					)
					if err != nil {
						l.Logger.Errorf("error synchronising database: %v", err)
					}
					l.Logger.Infof("%d values synchronised from '%s' to '%s'", updateCount, spreadsheetPath, redisAddress)

				case err := <-watcher.Errors:
					l.Logger.Warnf("error monitoring file for changes: %v", err.Error())
				}

			}
		}()

		err = watcher.Add(spreadsheetPath)
		cobra.CheckErr(err)
		l.Logger.Infof("watching for changes to %q", spreadsheetPath)

		var c = make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		<-c
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
