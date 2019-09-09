package pusher

import (
	"net/http"
	"time"

	"github.com/FlowerWrong/pusher/db"
	"github.com/FlowerWrong/pusher/log"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, appKey, client, version, protocolStr string) {
	protocol, _ := Str2Int(protocolStr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	if appKey != viper.GetString("APP_KEY") {
		log.Error("Error app key", appKey)
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		_ = conn.WriteMessage(websocket.TextMessage, ErrPack(4001))
		_ = conn.Close()
		return
	}

	if !funk.Contains(SupportedProtocolVersions, protocol) {
		log.Error("Unsupported protocol version", protocol)
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
		_ = conn.WriteMessage(websocket.TextMessage, ErrPack(4007))
		_ = conn.Close()
		return
	}

	socketID := GenerateSocketID()
	session := &Session{hub: hub, conn: conn, client: client, version: version, protocol: protocol, send: make(chan []byte, maxMessageSize), subscriptions: make(map[string]bool), pubSub: new(redis.PubSub), socketID: socketID, closed: false}
	session.hub.register <- session

	session.pubSub = db.Redis().Subscribe()
	go session.subscribePump()
	go session.writePump()
	go session.readPump()

	session.send <- EstablishPack(socketID)
}
