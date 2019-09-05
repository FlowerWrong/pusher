package pusher

import (
	"encoding/json"

	"github.com/FlowerWrong/pusher/api/forms"
	"github.com/thoas/go-funk"
)

var (
	PresenceMaximumMembersCount int64 = 100
	PresenceMaximumUserDataSize       = 1000
	PresenceMaximumUserIDLen          = 128
)

// PresenceChannelData ...
type PresenceChannelData struct {
	UserID   string            `json:"user_id"` // user id, not socket id
	UserInfo map[string]string `json:"user_info,omitempty"`
}

// MemberRemovedData ...
type MemberRemovedData struct {
	UserID string `json:"user_id"`
}

// PresenceData ...
type PresenceData struct {
	IDs   []string                     `json:"ids"`
	Hash  map[string]map[string]string `json:"hash"`
	Count int                          `json:"count"`
}

// GetPresenceChannelData ...
func GetPresenceChannelData(socketID, channel string) (*PresenceChannelData, error) {
	user := GetUserInfo(socketID)
	var presenceChannelData PresenceChannelData
	err := json.Unmarshal([]byte(user[UserInfoChannelDataKey(channel)]), &presenceChannelData)
	if err != nil {
		return nil, err
	}
	return &presenceChannelData, nil
}

// GetPresenceChannelDataByUser ...
func GetPresenceChannelDataByUser(user map[string]string, channel string) (*PresenceChannelData, error) {
	var presenceChannelData PresenceChannelData
	err := json.Unmarshal([]byte(user[UserInfoChannelDataKey(channel)]), &presenceChannelData)
	if err != nil {
		return nil, err
	}
	return &presenceChannelData, nil
}

// PresenceChannelUserIDs ...
func PresenceChannelUserIDs(channel string) []string {
	subs := GetSubscriptions(channel)
	var userIDs []string
	for _, socketID := range subs {
		data, err := GetPresenceChannelDataByUser(GetUserInfo(socketID), channel)
		if err != nil {
			continue
		}
		if funk.Contains(userIDs, data.UserID) {
			continue
		}
		userIDs = append(userIDs, data.UserID)
	}
	return userIDs
}

// PublishPresenceMemberAddedEvent ...
func PublishPresenceMemberAddedEvent(socketID, channel string) error {
	user := GetUserInfo(socketID)
	// member_added event
	eventForm := forms.EventForm{Name: PusherPresenceMemberAdded, Channel: channel, Data: user[UserInfoChannelDataKey(channel)], SocketID: socketID}
	err := PublishEventForm(eventForm)
	if err != nil {
		return err
	}

	// member_added hook
	presenceChannelData, err := GetPresenceChannelDataByUser(user, channel)
	if err != nil {
		return err
	}
	hookEvent := HookEvent{Name: "member_added", Channel: channel, UserID: presenceChannelData.UserID}
	go TriggerHook(&hookEvent)
	return nil
}

// PublishPresenceMemberRemovedEvent ...
func PublishPresenceMemberRemovedEvent(socketID, channel string) error {
	presenceChannelData, err := GetPresenceChannelData(socketID, channel)
	if err != nil {
		return err
	}

	// member_removed hook
	hookEvent := HookEvent{Name: "member_removed", Channel: channel, UserID: presenceChannelData.UserID}
	go TriggerHook(&hookEvent)

	memberRemovedData := MemberRemovedData{UserID: presenceChannelData.UserID}
	b, err := json.Marshal(memberRemovedData)
	if err != nil {
		return err
	}

	// member_removed event
	eventForm := forms.EventForm{Name: PusherPresenceMemberRemoved, Channel: channel, Data: string(b[:]), SocketID: socketID}
	err = PublishEventForm(eventForm)
	if err != nil {
		return err
	}
	return nil
}
