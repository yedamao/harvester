package function

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-redis/redis/v8"
)

// Dump the key data
func Dump(key string, client *redis.Client) error {
	fmt.Println(key, client.ObjectIdleTime(context.TODO(), key).Val())
	return nil
}

func encode(raw string) string {
	return url.QueryEscape(raw)
}
