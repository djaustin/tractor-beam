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
func SyncDatabase(ctx context.Context, client *redis.Client, redisPrefix, srcPath, worksheet, keyColumn, valColumn string) (int, error) {
	l.Logger.Infof("extracting key-value pairs from spreadsheet")
	pairs, err := spreadsheet.ExtractPairs(srcPath, worksheet, keyColumn, valColumn)
	if err != nil {
		return 0, fmt.Errorf("unable to extract data from spreadsheet: %w", err)
	}
	l.Logger.Infof("key-value pairs extracted from spreadsheet")
	count := 0
	l.Logger.Infof("starting insert of %d key-value pairs to database", len(pairs))
	for k, v := range pairs {
		redisKey := redisPrefix + k
		res := client.Set(ctx, redisKey, v, 0)
		if err := res.Err(); err != nil {
			return count, fmt.Errorf("unable to set data pair %s=%s: %w", redisKey, v, err)
		}
		count++
	}
	return count, nil
}
