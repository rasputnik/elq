package main

import (
	"flag"
	"github.com/buger/goterm"
	"github.com/go-redis/redis"
	"sort"
	"time"
)

func main() {

	var server, auth string
	var interval int
	flag.StringVar(&server, "h", "localhost:6379", "redis server and port")
	flag.StringVar(&auth, "a", "", "password")
	flag.IntVar(&interval, "i", 3, "refresh interval in seconds")

	flag.Parse()

	var redis_opts = redis.Options{
		Addr:     server,
		Password: auth,
		DB:       0,
	}

	client := redis.NewClient(&redis_opts)

	tic := time.Tick(time.Duration(interval) * time.Second)
	for _ = range tic {
		checkQueues(client)
	}

}

func checkQueues(r *redis.Client) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)

	_, err := r.Ping().Result()
	if err != nil {
		goterm.Println("ping failed - bad host?")
	}

	queues, err := r.Keys("*").Result()
	if err != nil {
		goterm.Println("no keys?")
	}
	sort.Strings(queues)

	goterm.Clear()
	for _, q := range queues {
		qlen, err := r.LLen(q).Result()
		if err != nil {
			// probably not a list, skip
			continue
		}
		goterm.Printf("%s : %d\n", q, qlen)
	}
	goterm.Flush()
}
