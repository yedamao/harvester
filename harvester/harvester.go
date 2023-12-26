package harvester

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// JustDoItFunc key process func
type JustDoItFunc func(key string, client *redis.Client) error

// Harvester concurrent scan redis instance
// execute JustDoIt func on every founded key
type Harvester struct {
	dbs []*redis.Client

	prefix   string
	justDoIt JustDoItFunc

	// pipeline pass key between scan loop and worker routine
	pipeline chan *obj4Pipe

	parallel int

	done              chan struct{}
	wg                sync.WaitGroup
	scanLoopWaitGroup sync.WaitGroup
}

// New Harvester
func New(cfg *Config, fn JustDoItFunc) (*Harvester, error) {

	var dbs []*redis.Client

	for _, option := range cfg.Opts {
		rdb := redis.NewClient(&redis.Options{
			Addr:     option.Addr,
			Password: option.Password, // no password set
			DB:       option.DB,       // use default DB
		})

		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			return nil, err
		}

		dbs = append(dbs, rdb)
	}

	h := &Harvester{
		dbs:      dbs,
		prefix:   cfg.Prefix,
		parallel: cfg.Parallel,
		justDoIt: fn,
		pipeline: make(chan *obj4Pipe),
		done:     make(chan struct{}),
	}

	return h, nil
}

// Run Harvester
func (h *Harvester) Run() {

	// spawn worker routine
	for i := 0; i < h.parallel; i++ {
		h.wg.Add(1)
		go h.worker()
	}

	for _, db := range h.dbs {
		h.scanLoopWaitGroup.Add(1)
		go h.scanLoop(db)
	}

	h.scanLoopWaitGroup.Wait()
	logrus.Info("scan loop worker. all done")

	close(h.pipeline)

	h.wg.Wait()
	logrus.Info("justDoIt worker. all done")

	logrus.Info("harvester stopped")
}

func (h *Harvester) scanLoop(db *redis.Client) {
	defer h.scanLoopWaitGroup.Done()

	ctx := context.Background()
	opt := db.Options()

	var cursor uint64 = 0
	for {

		select {
		case <-h.done:
			return
		default:
		}

		keys, cur, err := db.Scan(ctx, cursor, h.prefix, 100).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			h.pipeline <- &obj4Pipe{db, key}
		}

		if cur == 0 {
			// Scan end
			logrus.Infof("instance: %s/%d", opt.Addr, opt.DB)
			logrus.Infof("Scan end. cursor %d\n", cur)
			return
		}
		cursor = cur
	}
}

type obj4Pipe struct {
	db  *redis.Client
	key string
}

func (h *Harvester) worker() {
	defer h.wg.Done()

	for {
		obj, ok := <-h.pipeline
		if !ok {
			return
		}

		if err := h.justDoIt(obj.key, obj.db); err != nil {
			logrus.Errorf("do it err: %v", err)
		}
	}
}

// Stop a running Harvester
func (h *Harvester) Stop() {
	close(h.done)
}
