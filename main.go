package main

import (
  "flag"
  "fmt"
  "github.com/go-redis/redis"
)

func main() {

  var server string
  flag.StringVar(&server, "h", "localhost:6379", "redis server and port")

  flag.Parse()

  var redis_opts = redis.Options{
    Addr: server,
    Password: "",
    DB: 0,
  } 

  client := redis.NewClient(&redis_opts)

  pong, err := client.Ping().Result()
  fmt.Println(pong, err)

}
