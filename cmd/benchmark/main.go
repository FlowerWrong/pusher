package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/pusher/pusher-http-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	client := pusher.Client{
		AppID:   "854355",
		Key:     "5c107909cb0b804d6e21",
		Secret:  "1ff025dccc3dcba9ec82",
		Cluster: "us3",
		Secure:  true,
	}

	data := map[string]string{"message": "hello world"}
	client.Trigger("presence-my-channel", "new-message", data)
}
