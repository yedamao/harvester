package function

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// DumpStrings the key and value of string
func DumpStrings(key string, client *redis.Client) error {
	val, err := client.Get(context.TODO(), key).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Println(key, encode(val))
	return nil
}
