package api

import (
	"net/http"
	"strings"

	"github.com/FlowerWrong/pusher"
	"github.com/gin-gonic/gin"
)

// @doc https://pusher.com/docs/channels/library_auth_reference/rest-api#channels

// ChannelIndex ...
func ChannelIndex(c *gin.Context) {
	filterByPrefix := c.Query("filter_by_prefix")
	info := c.Query("info")

	var selectedChannels []string
	channels := pusher.Channels()
	for _, channel := range channels {
		if pusher.ChannelSubscriberCount(channel) > 0 {
			// occupied channels
			if filterByPrefix != "" && strings.HasPrefix(channel, filterByPrefix) {
				selectedChannels = append(selectedChannels, channel)
			} else {
				selectedChannels = append(selectedChannels, channel)
			}
		}
	}

	splitFn := func(c rune) bool {
		return c == ','
	}
	infoFields := strings.FieldsFunc(info, splitFn)
	data := make(map[string]map[string]interface{})
	for _, channel := range selectedChannels {
		if len(infoFields) > 0 {
			for _, infoField := range infoFields {
				if infoField == "user_count" && pusher.IsPresenceChannel(channel) {
					data[channel] = map[string]interface{}{"user_count": len(pusher.PresenceChannelUserIDs(channel))}
				}
			}
		} else {
			data[channel] = make(map[string]interface{})
		}
	}

	c.JSON(http.StatusOK, gin.H{"channels": data})
}

// ChannelShow ...
func ChannelShow(c *gin.Context) {
	channel := c.Param("channel_name")
	info := c.Query("info")

	splitFn := func(c rune) bool {
		return c == ','
	}
	infoFields := strings.FieldsFunc(info, splitFn)

	data := make(map[string]interface{})
	subscriptionsCount := pusher.ChannelSubscriberCount(channel)
	data["occupied"] = subscriptionsCount > 0
	for _, infoField := range infoFields {
		if infoField == "user_count" && pusher.IsPresenceChannel(channel) {
			data["user_count"] = len(pusher.PresenceChannelUserIDs(channel))
		}
		if infoField == "subscription_count" {
			data["subscription_count"] = subscriptionsCount
		}
	}

	c.JSON(http.StatusOK, data)
}

// ChannelUsers ...
func ChannelUsers(c *gin.Context) {
	channel := c.Param("channel_name")

	if !pusher.IsPresenceChannel(channel) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "presence only"})
		return
	}

	socketIDs := pusher.GetSubscriptions(channel)
	var data []map[string]string
	for _, socketID := range socketIDs {
		presenceChannelData, err := pusher.GetPresenceChannelData(socketID, channel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		data = append(data, map[string]string{"id": presenceChannelData.UserID})
	}

	c.JSON(http.StatusOK, gin.H{"users": data})
}
