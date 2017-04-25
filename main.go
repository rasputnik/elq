package main

import (
  "flag"
  "fmt"
  "github.com/go-redis/redis"
)

func main() {

  var server, auth string
  flag.StringVar(&server, "h", "localhost:6379", "redis server and port")
  flag.StringVar(&auth, "a", "", "password")

  flag.Parse()

  var redis_opts = redis.Options{
    Addr: server,
    Password: auth,
    DB: 0,
  } 

  client := redis.NewClient(&redis_opts)

  pong, err := client.Ping().Result()
  fmt.Println(pong, err)

}
