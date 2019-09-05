package pusher

import (
	"encoding/json"
	"errors"

	"github.com/FlowerWrong/pusher/api/forms"
	"github.com/FlowerWrong/pusher/log"
)

// ParsePayload ...
func (s *Session) ParsePayload(payloadBytes []byte) error {
	var requestEvent RequestEvent
	err := json.Unmarshal(payloadBytes, &requestEvent)
	if err != nil {
		return err
	}

	switch requestEvent.Event {
	case PusherSubscribe:
		var subscribe Subscribe
		err = json.Unmarshal(payloadBytes, &subscribe)
		if err != nil {
			return err
		}

		channel := subscribe.Data.Channel

		if !ValidChannel(channel) {
			return ErrInvalidChannel
		}

		if s.subscriptions[channel] {
			return errors.New("have subscribed")
		}
		if IsPresenceChannel(channel) {
			if ChannelSubscriberCount(channel) >= PresenceMaximumMembersCount {
				return ErrPresenceSubscriberTooMuch
			}

			userData := subscribe.Data.ChannelData
			if len(userData) > PresenceMaximumUserDataSize {
				return ErrPresenceUserDataTooMuch
			}

			var presenceChannelData PresenceChannelData
			err := json.Unmarshal([]byte(userData), &presenceChannelData)
			if err != nil {
				return err
			}
			if len(presenceChannelData.UserID) > PresenceMaximumUserIDLen {
				return ErrPresenceUserIDTooLong
			}

			err = SetUserInfo(s.socketID, map[string]interface{}{
				UserInfoChannelDataKey(channel): userData,
				UserInfoAuthKey(channel):        subscribe.Data.Auth,
			})
			if err != nil {
				return err
			}
		}
		if IsPrivateChannel(channel) {
			err = SetUserInfoField(s.socketID, UserInfoAuthKey(channel), subscribe.Data.Auth)
			if err != nil {
				return err
			}
		}

		err = s.subscribe(channel)
		if err != nil {
			return err
		}
	case PusherUnsubscribe:
		var unsubscribe Unsubscribe
		err = json.Unmarshal(payloadBytes, &unsubscribe)
		if err != nil {
			return err
		}

		channel := unsubscribe.Data.Channel
		if !ValidChannel(channel) {
			return ErrInvalidChannel
		}

		if !s.subscriptions[channel] {
			return errors.New("not subscribed")
		}
		err = s.unsubscribe(channel)
		if err != nil {
			return err
		}
	case PusherPing:
		s.Send(PingPongPack(PusherPong))
	case PusherPong:
		s.resetReadTimeout()
	default:
		if IsClientEvent(requestEvent.Event) {
			// @doc https://pusher.com/docs/channels/server_api/excluding-event-recipients

			// TODO rate limit https://pusher.com/docs/channels/using_channels/events#rate-limit-your-events

			var clientChannelEvent ClientChannelEvent
			err := json.Unmarshal(payloadBytes, &clientChannelEvent)
			if err != nil {
				return err
			}

			channel := clientChannelEvent.Channel

			if !s.subscriptions[channel] {
				return errors.New("you must subscribed to this channel")
			}

			if !IsPresenceChannel(channel) && !IsPrivateChannel(channel) {
				return errors.New("client event only support private or presence channel")
			}

			var eventForm forms.EventForm
			eventForm.Channel = channel
			eventForm.Name = clientChannelEvent.Event
			eventForm.SocketID = s.socketID // exclude
			dataBytes, err := clientChannelEvent.Data.MarshalJSON()
			if err != nil {
				return err
			}
			eventForm.Data = string(dataBytes[:])

			hookEvent := HookEvent{Name: "client_event", Channel: channel, Event: clientChannelEvent.Event, Data: string(dataBytes[:]), SocketID: s.socketID}

			if IsPresenceChannel(channel) {
				presenceChannelData, err := GetPresenceChannelData(s.socketID, channel)
				if err != nil {
					return err
				}
				hookEvent.UserID = presenceChannelData.UserID

				eventForm.UserID = presenceChannelData.UserID
			}

			err = PublishEventForm(eventForm)
			if err != nil {
				return err
			}

			go TriggerHook(&hookEvent)
		} else {
			log.Error("Unsupported event name", requestEvent.Event)
		}
	}

	return nil
}
