package db

import (
	"context"
	"fmt"

	l "github.com/djaustin/tractor-beam/logger"
	"github.com/djaustin/tractor-beam/spreadsheet"
	"github.com/go-redis/redis/v8"
)

// SyncDatabase extracts key-value pairs from a spreadsheet and updates a Redis database with them.
// A count of updated pairs is returned on success.
func SyncDatabase(ctx context.Context, client *redis.Client, redisPrefix, srcPath, worksheet, keyColumn string, valColumns ...string) (int, error) {
	l.Logger.Debug("extracting key-value pairs from spreadsheet")
	maps, err := spreadsheet.ExtractMaps(srcPath, worksheet, keyColumn, valColumns...)
	if err != nil {
		return 0, fmt.Errorf("unable to extract data from spreadsheet: %w", err)
	}
	l.Logger.Debug("key-value pairs extracted from spreadsheet")
	count := 0
	l.Logger.Debug("starting insert of %d key-value pairs to database", len(maps))
	for k, v := range maps {
		redisKey := redisPrefix + k
		res := client.HSet(ctx, redisKey, v)
		if err := res.Err(); err != nil {
			return count, fmt.Errorf("unable to set data pair %s=%s: %w", redisKey, v, err)
		}
		count++
	}
	return count, nil
}
