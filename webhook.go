package pusher

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

type webhook struct {
	TimeMs int64       `json:"time_ms"`
	Events []HookEvent `json:"events"`
}

// HookEvent ...
type HookEvent struct {
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Event    string `json:"event,omitempty"`
	Data     string `json:"data,omitempty"`
	SocketID string `json:"socket_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// TriggerHook @doc https://pusher.com/docs/channels/server_api/webhooks
func TriggerHook(events ...*HookEvent) {
	if !viper.GetBool("pusher_webhook_enabled") {
		return
	}
	timeMs := time.Now().UnixNano() / 1e6
	hook := webhook{TimeMs: timeMs}

	for _, event := range events {
		hook.Events = append(hook.Events, *event)
	}

	b, err := json.Marshal(hook)
	if err != nil {
		panic(err)
	}

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Pusher-Key", viper.GetString("pusher_app_key")).
		SetHeader("X-Pusher-Signature", HmacSignature(string(b[:]), viper.GetString("pusher_app_secret"))).
		SetBody(b).
		Post(viper.GetString("pusher_webhook_url"))
	if !resp.IsSuccess() {
		// TODO retry
	}
}
