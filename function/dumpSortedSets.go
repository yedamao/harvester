package function

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// DumpSortedSets the key and value of string
func DumpSortedSets(key string, client *redis.Client) error {
	var offset, limit int64 = 0, 50
	for {
		items, err := client.ZRangeWithScores(context.TODO(), key, offset, offset+limit).Result()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ZRangeWithScores error: %v\n", err)
			break
		}

		// 遍历结束
		if len(items) == 0 {
			break
		}

		for _, item := range items {
			member, _ := item.Member.(string)
			fmt.Println(key, int64(item.Score), encode(member))
		}

		offset += int64(len(items))
	}
	return nil
}
