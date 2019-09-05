package pusher

import (
	"github.com/FlowerWrong/pusher/db"
)

// channel set: pusher:channels -> [channel_name string]
// channel's subscriptions set: pusher:{channel_name}:subscriptions -> [socket_id string]
// user set: pusher:users -> [socket_id string]
// user info: pusher:{socket_id}:info -> hash{user_id: string...}

const (
	PusherChannelSetKey = "pusher:channels"
	PusherUserSetKey    = "pusher:users"
)

func subscriptionKey(channel string) string {
	return "pusher:" + channel + ":subscriptions"
}

func userInfoKey(socketID string) string {
	return "pusher:" + socketID + ":info"
}

// AddChannel to redis set
func AddChannel(channel string) error {
	if !db.Redis().SIsMember(PusherChannelSetKey, channel).Val() {
		return db.Redis().SAdd(PusherChannelSetKey, channel).Err()
	}
	return nil
}

// Channels get channels set
func Channels() []string {
	return db.Redis().SMembers(PusherChannelSetKey).Val()
}

// AddUser to redis set
func AddUser(socketID string) error {
	if !db.Redis().SIsMember(PusherUserSetKey, socketID).Val() {
		return db.Redis().SAdd(PusherUserSetKey, socketID).Err()
	}
	return nil
}

// DelUserAndSubscriptions del user info hash and del user id from user set
func DelUserAndSubscriptions(socketID string, subscriptions map[string]bool) error {
	pl := db.Redis().TxPipeline() // transaction pipleline
	for channel, _ := range subscriptions {
		pl.SRem(subscriptionKey(channel), socketID)
	}
	pl.SRem(PusherUserSetKey, socketID)
	pl.Del(userInfoKey(socketID))
	_, err := pl.Exec()
	return err
}

// AddSubscription add user to channel subscriptions set
func AddSubscription(channel, socketID string) error {
	key := subscriptionKey(channel)
	if !db.Redis().SIsMember(key, socketID).Val() {
		return db.Redis().SAdd(key, socketID).Err()
	}
	return nil
}

// GetSubscriptions ...
func GetSubscriptions(channel string) []string {
	return db.Redis().SMembers(subscriptionKey(channel)).Val()
}

// ChannelSubscriberCount ...
func ChannelSubscriberCount(channel string) int64 {
	return db.Redis().SCard(subscriptionKey(channel)).Val()
}

// AddChannelAndSubscription ...
func AddChannelAndSubscription(channel, socketID string) error {
	pl := db.Redis().TxPipeline() // transaction pipleline
	pl.SAdd(PusherChannelSetKey, channel)
	key := subscriptionKey(channel)
	pl.SAdd(key, socketID)
	_, err := pl.Exec()
	return err
}

// DelSubscription del user from channel subscriptions set
func DelSubscription(channel, socketID string) error {
	return db.Redis().SRem(subscriptionKey(channel), socketID).Err()
}

// SetUserInfo set user info to hash
func SetUserInfo(socketID string, info map[string]interface{}) error {
	return db.Redis().HMSet(userInfoKey(socketID), info).Err()
}

// SetUserInfoField set user field to hash
func SetUserInfoField(socketID string, k string, v interface{}) error {
	return db.Redis().HSet(userInfoKey(socketID), k, v).Err()
}

// GetUserInfo ...
func GetUserInfo(socketID string) map[string]string {
	return db.Redis().HGetAll(userInfoKey(socketID)).Val()
}

// DelUserInfoFields ...
func DelUserInfoFields(socketID string, fields ...string) error {
	return db.Redis().HDel(userInfoKey(socketID), fields...).Err()
}

// UserInfoChannelDataKey ...
func UserInfoChannelDataKey(channel string) string {
	return channel + "-channel_data"
}

// UserInfoAuthKey ...
func UserInfoAuthKey(channel string) string {
	return channel + "-auth"
}
