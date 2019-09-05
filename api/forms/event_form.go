package forms

// EventForm ...
type EventForm struct {
	Name     string   `form:"name" json:"name" binding:"required"`
	Data     string   `form:"data" json:"data" binding:"required"` // limit to 10kb
	Channels []string `form:"channels" json:"channels,omitempty"`  // limit to 100 channel
	Channel  string   `form:"channel" json:"channel,omitempty"`
	SocketID string   `form:"socket_id" json:"socket_id,omitempty"` // excludes the event from being sent to
	UserID   string   `json:"user_id,omitempty"`                    // optional, present only if this is a `client event` on a `presence channel`
}

// BatchEventForm ...
type BatchEventForm struct {
	Batch []EventForm `form:"batch" json:"batch" binding:"required"`
}
