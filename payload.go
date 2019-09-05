package pusher

import (
	"encoding/json"

	"github.com/spf13/viper"
)

const (
	PusherError                 = "pusher:error"
	PusherSubscribe             = "pusher:subscribe"
	PusherUnsubscribe           = "pusher:unsubscribe"
	PusherPing                  = "pusher:ping"
	PusherPong                  = "pusher:pong"
	PusherSubscriptionSucceeded = "pusher_internal:subscription_succeeded"
	PusherPresenceMemberAdded   = "pusher_internal:member_added"
	PusherPresenceMemberRemoved = "pusher_internal:member_removed"
)

// RequestEvent ...
type RequestEvent struct {
	Event string `json:"event"`
}

// Payload ...
type Payload struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// PayloadPack ...
func PayloadPack(event, data string) []byte {
	payload := Payload{Event: event, Data: data}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return payloadJSON
}

// EstablishData Channels -> Client
type EstablishData struct {
	SocketID        string `json:"socket_id"`
	ActivityTimeout int    `json:"activity_timeout"` // Protocol 7 and above
}

// EstablishPack Channels -> Client
func EstablishPack(socketID string) []byte {
	data := EstablishData{SocketID: socketID, ActivityTimeout: viper.GetInt("pusher_activity_timeout")}
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return PayloadPack("pusher:connection_established", string(b[:]))
}

// ErrData Channels -> Client
type ErrData struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"` // optional
}

// ErrPack Channels -> Client
func ErrPack(code int) []byte {
	data := ErrData{Message: ErrCodes[code], Code: code}
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return PayloadPack(PusherError, string(b[:]))
}

// Subscribe ...
type Subscribe struct {
	Event string        `json:"event"`
	Data  SubscribeData `json:"data"`
}

// SubscribeData ...
type SubscribeData struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth,omitempty"`         // optional
	ChannelData string `json:"channel_data,omitempty"` // optional PresenceChannelData
}

// PingPong ...
type PingPong struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// PingPongPack ...
func PingPongPack(event string) []byte {
	data := PingPong{Event: event, Data: map[string]string{}}
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

// Unsubscribe ...
type Unsubscribe struct {
	Event string          `json:"event"`
	Data  UnsubscribeData `json:"data"`
}

// UnsubscribeData ...
type UnsubscribeData struct {
	Channel string `json:"channel"`
}

// SubscriptionSucceeded CHANNELS -> CLIENT
// @doc https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol#presence-channel-events
type SubscriptionSucceeded struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// SubscriptionSucceededPack ...
func SubscriptionSucceededPack(channel, data string) []byte {
	d := SubscriptionSucceeded{Event: PusherSubscriptionSucceeded, Channel: channel, Data: data}
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return b
}

// ClientChannelEvent ...
type ClientChannelEvent struct {
	Event   string          `json:"event"`
	Data    json.RawMessage `json:"data"`
	Channel string          `json:"channel"`
}
