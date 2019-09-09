package pusher

import (
	"bufio"
	"net/http"
	"os"
	"time"

	"github.com/FlowerWrong/pusher/env"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// AppEnv is app env
var AppEnv string

const (
	// DEVELOPMENT env
	DEVELOPMENT = "development"
	// TEST env
	TEST = "test"
	// PRODUCTION env
	PRODUCTION = "production"
)
const (
	writeWait      = 10 * time.Second // Time allowed to write a message to the peer.
	maxMessageSize = 1024             // Maximum message size allowed from peer.
)

var (
	pongTimeout     time.Duration
	activityTimeout time.Duration
	readTimeout     time.Duration
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

// Setup ...
func Setup(file string) error {
	AppEnv = env.Get("APP_ENV", DEVELOPMENT)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(bufio.NewReader(f))

	pongTimeout = time.Duration(viper.GetInt64("PONG_TIMEOUT")) * time.Second
	activityTimeout = time.Duration(viper.GetInt64("ACTIVITY_TIMEOUT")) * time.Second
	readTimeout = activityTimeout / 9 * 10

	return nil
}
