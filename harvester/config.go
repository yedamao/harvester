package harvester

import (
	"github.com/go-redis/redis/v8"
)

// Config for Harvester
type Config struct {
	Prefix   string
	Opts     []*redis.Options
	Parallel int
}
