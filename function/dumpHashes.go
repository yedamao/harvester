package function

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// DumpHashes the key and value of string
func DumpHashes(key string, client *redis.Client) error {
	kvs, err := client.HGetAll(context.TODO(), key).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for field, value := range kvs {
		fmt.Println(key, field, encode(value))
	}
	return nil
}
