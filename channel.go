package pusher

import (
	"encoding/json"

	"github.com/FlowerWrong/pusher/api/forms"
	"github.com/FlowerWrong/pusher/db"
)

// ChannelEvent CHANNELS -> CLIENT
type ChannelEvent struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
	UserID  string `json:"user_id,omitempty"` // optional, present only if this is a `client event` on a `presence channel`
}

// PublishEventForm ...
func PublishEventForm(eventForm forms.EventForm) error {
	b, err := json.Marshal(eventForm)
	err = db.Redis().Publish(eventForm.Channel, b).Err()
	if err != nil {
		return err
	}
	return nil
}
