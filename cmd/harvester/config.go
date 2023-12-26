package main

import (
	"strings"

	"github.com/go-redis/redis/v8"
)

func parseRedisURL(str string) ([]*redis.Options, error) {
	var opts []*redis.Options
	for _, url := range strings.Split(redisURL, ",") {
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}

	return opts, nil
}
