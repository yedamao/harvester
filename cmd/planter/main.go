package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

var (
	redisURL string

	sourceFile string
	parallel   int
)

func init() {
	flag.StringVar(&sourceFile, "sourceFile", "", "source file to be imported to redis")
	flag.StringVar(&redisURL, "redisUrl", "", "redis://<user>:<password>@<host>:<port>/<db_number>")
}

func main() {
	flag.Parse()

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, opt)

	cli := redis.NewClient(opt)

	if err := cli.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "Would you like to launch the planter? [Y/n]")
	var cmd string
	fmt.Scanf("%s", &cmd)
	if strings.ToUpper(cmd) != "Y" {
		fmt.Println("Mission Canceled")
		os.Exit(-1)
	}

	fd, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		segments := strings.Split(line, " ")
		if len(segments) != 2 {
			panic(line)
		}

		if err := cli.Set(context.TODO(), segments[0], segments[1], 0).Err(); err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
