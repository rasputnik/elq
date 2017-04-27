package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"sort"
)

func main() {

	var server, auth string
	flag.StringVar(&server, "h", "localhost:6379", "redis server and port")
	flag.StringVar(&auth, "a", "", "password")

	flag.Parse()

	var redis_opts = redis.Options{
		Addr:     server,
		Password: auth,
		DB:       0,
	}

	client := redis.NewClient(&redis_opts)

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("ping failed - bad host? bombing out.")
	}

	queues, err := client.Keys("*").Result()
	if err != nil {
		fmt.Println("no keys? bombing out.")
	}
	sort.Strings(queues)
	for _, q := range queues {
		qlen, err := client.LLen(q).Result()
		if err != nil {
			// probably not a list, skip
			continue
		}
		fmt.Printf("%s : %d\n", q, qlen)
	}

}
