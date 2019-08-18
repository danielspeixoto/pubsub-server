package main

import (
	"pubsub/src"
	"time"
)

func main() {
	src.StartServer(10001, time.Millisecond * 250)
}

