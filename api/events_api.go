package api

import (
	"net/http"

	"github.com/FlowerWrong/pusher"
	"github.com/FlowerWrong/pusher/api/forms"
	"github.com/gin-gonic/gin"
)

// @doc https://pusher.com/docs/channels/library_auth_reference/rest-api

var maxTriggerableChannels = 100

// EventTrigger ...
func EventTrigger(c *gin.Context) {
	var eventForm forms.EventForm
	if err := c.ShouldBindJSON(&eventForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// limit to 10kb
	if len(eventForm.Data) > 10*1000 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": pusher.ErrAPIReqEventDataTooLarge.Error()})
		return
	}

	if len(eventForm.Channel) == 0 && len(eventForm.Channels) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing channel params"})
		return
	}

	if len(eventForm.Channel) > 0 && len(eventForm.Channels) == 0 {
		eventForm.Channels = append(eventForm.Channels, eventForm.Channel)
	}

	// limit to 100
	if len(eventForm.Channels) > maxTriggerableChannels {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": pusher.ErrAPIReqEventChannelsSizeTooLong.Error()})
		return
	}

	// publish to redis
	for _, channel := range eventForm.Channels {
		if !pusher.ValidChannel(channel) {
			c.JSON(http.StatusBadRequest, gin.H{"error": pusher.ErrInvalidChannel.Error()})
			return
		}
		err := publish(channel, eventForm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

// BatchEventTrigger ...
func BatchEventTrigger(c *gin.Context) {
	var batchEventForm forms.BatchEventForm
	if err := c.ShouldBindJSON(&batchEventForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, eventForm := range batchEventForm.Batch {
		// limit to 10kb
		if len(eventForm.Data) > 10*1000 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": pusher.ErrAPIReqEventDataTooLarge.Error()})
			return
		}

		if !pusher.ValidChannel(eventForm.Channel) {
			c.JSON(http.StatusBadRequest, gin.H{"error": pusher.ErrInvalidChannel.Error()})
			return
		}

		err := publish(eventForm.Channel, eventForm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func publish(channel string, eventForm forms.EventForm) error {
	event := forms.EventForm{Name: eventForm.Name, Data: eventForm.Data, Channel: channel, SocketID: eventForm.SocketID}
	return pusher.PublishEventForm(event)
}
