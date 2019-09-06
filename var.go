package pusher

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
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

// InitVar ...
func InitVar() {
	pongTimeout = time.Duration(viper.GetInt64("pusher_pong_timeout")) * time.Second
	activityTimeout = time.Duration(viper.GetInt64("pusher_activity_timeout")) * time.Second
	readTimeout = activityTimeout / 9 * 10
}
