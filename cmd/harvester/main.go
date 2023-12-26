package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/yedamao/harvester/function"
	"github.com/yedamao/harvester/harvester"
)

var (
	cfg                  string
	cluster              string
	master               bool
	clusterShardingCount int

	redisURL string

	matchPattern string
	action       string
	parallel     int
)

func init() {
	flag.StringVar(&redisURL, "redisUrl", "", "eg: redis://<user>:<password>@<host>:<port>/<db_number>. separated by commas.")

	flag.StringVar(&matchPattern, "matchPattern", "*", "scan match pattern")
	flag.StringVar(&action, "action", "dump", "dump the key data")
	flag.IntVar(&parallel, "parallel", 1, "the number of worker to run parallel")
}

func main() {
	flag.Parse()

	options, err := parseRedisURL(redisURL)
	if err != nil {
		panic(err)
	}

	for _, op := range options {
		fmt.Fprintln(os.Stderr, op)
	}

	fmt.Fprintln(os.Stderr, "Would you like to launch the harvester? [Y/n]")

	var cmd string
	fmt.Scanf("%s", &cmd)
	if strings.ToUpper(cmd) != "Y" {
		fmt.Println("Mission Canceled")
		os.Exit(-1)
	}

	cfg := &harvester.Config{
		Prefix:   matchPattern,
		Opts:     options,
		Parallel: parallel,
	}

	fn := function.Dump
	if action == "dump-strings" {
		fn = function.DumpStrings
	}

	if action == "dump-sorted-sets" {
		fn = function.DumpSortedSets
	}

	if action == "dump-hashes" {
		fn = function.DumpHashes
	}

	h, err := harvester.New(cfg, fn)
	if err != nil {
		panic(err)
	}

	handleSignals(h.Stop)

	h.Run()
}

func handleSignals(stopFunction func()) {
	var callback sync.Once

	// On ^C or SIGTERM, gracefully stop the sniffer
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		log.Println("service", "Received sigterm/sigint, stopping")
		callback.Do(stopFunction)
	}()
}
