package main

import (
	"flag"
	"github.com/buger/goterm"
	"github.com/go-redis/redis"
	"sort"
	"time"
)

// key = queue name,
// val = length of queue
type QState map[string]int64

// return sorted key list
func (qs QState) Keys() []string {
	var keys []string
	for k, _ := range qs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

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

	var queueCounts, lastCounts QState

	tic := time.Tick(time.Duration(interval) * time.Second)
	for _ = range tic {
		queueCounts = checkQueues(client)
		displayCounts(queueCounts, lastCounts)
		lastCounts = queueCounts
	}

}

func displayCounts(q QState, oldq QState) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	for _, k := range q.Keys() {
		// check previous value
		oldv := oldq[k]
		goterm.Printf("%s : %d (delta: %d)\n", k, q[k], (q[k] - oldv))
	}
	goterm.Flush()
}

func checkQueues(r *redis.Client) QState {

	_, err := r.Ping().Result()
	if err != nil {
		goterm.Println("ping failed - bad host?")
	}

	queues, err := r.Keys("*").Result()
	if err != nil {
		goterm.Println("no keys?")
	}
	var results QState
	results = make(QState)

	for _, q := range queues {
		qlen, err := r.LLen(q).Result()
		if err != nil {
			// probably not a list, skip
			continue
		}
		results[q] = qlen
	}
	return results
}
