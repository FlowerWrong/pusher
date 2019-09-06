package pusher

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/FlowerWrong/pusher/api/forms"
	"github.com/FlowerWrong/pusher/log"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// Session is a middleman between the websocket connection and the hub.
type Session struct {
	hub           *Hub
	conn          *websocket.Conn
	client        string
	version       string
	protocol      int
	send          chan []byte // Buffered channel of outbound messages.
	mutex         sync.Mutex
	subscriptions map[string]bool
	pubSub        *redis.PubSub
	socketID      string
	closed        bool
}

func (s *Session) resetReadTimeout() {
	_ = s.conn.SetReadDeadline(time.Now().Add(readTimeout))
}

func (s *Session) resetPongTimeout() {
	_ = s.conn.SetReadDeadline(time.Now().Add(pongTimeout))
}

func (s *Session) write(message []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_ = s.conn.SetWriteDeadline(time.Now().Add(writeWait))
	w, err := s.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, _ = w.Write(message)
	return w.Close()
}

// Send data to client connection
func (s *Session) Send(msg []byte) {
	s.mutex.Lock()

	if s.send == nil {
		s.mutex.Unlock()
		return
	}

	select {
	case s.send <- msg:
	default:
		if s.send != nil {
			close(s.send)
		}
		s.send = nil
	}
	s.mutex.Unlock()
}

// Close websocket connection with the specified reason
func (s *Session) Close(code int) {
	s.mutex.Lock()
	if s.closed {
		s.mutex.Unlock()
		return
	}
	s.closed = true
	s.mutex.Unlock()

	_ = s.conn.SetWriteDeadline(time.Now().Add(time.Second))
	_ = s.conn.WriteMessage(websocket.TextMessage, ErrPack(code))
	_ = s.conn.Close()
	log.Infoln(s.socketID, "leaved")
}

func (s *Session) stop(code int) error {
	if s.pubSub != nil {
		_ = s.pubSub.Close()
	}

	for channel := range s.subscriptions {
		if IsPresenceChannel(channel) {
			err := PublishPresenceMemberRemovedEvent(s.socketID, channel)
			if err != nil {
				return err
			}
		}
	}

	err := DelUserAndSubscriptions(s.socketID, s.subscriptions)
	if err != nil {
		return err
	}

	var hookEvents []*HookEvent
	for channel := range s.subscriptions {
		if ChannelSubscriberCount(channel) == 0 {
			hookEvent := HookEvent{Name: "channel_vacated", Channel: channel}
			hookEvents = append(hookEvents, &hookEvent)
		}
	}
	go TriggerHook(hookEvents...)

	for channel := range s.subscriptions {
		delete(s.subscriptions, channel)
	}

	s.Close(code)
	return nil
}

func (s *Session) subscribe(channel string) error {
	err := s.pubSub.Subscribe(channel)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) unsubscribe(channel string) error {
	err := s.pubSub.Unsubscribe(channel)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) unsubscribeSuccessCallback(channel string) error {
	delete(s.subscriptions, channel)
	err := DelSubscription(channel, s.socketID)
	if err != nil {
		return err
	}

	if ChannelSubscriberCount(channel) == 0 {
		hookEvent := HookEvent{Name: "channel_vacated", Channel: channel}
		go TriggerHook(&hookEvent)
	}

	if IsPresenceChannel(channel) {
		err = PublishPresenceMemberRemovedEvent(s.socketID, channel)
		if err != nil {
			return err
		}

		err = DelUserInfoFields(s.socketID, UserInfoAuthKey(channel), UserInfoChannelDataKey(channel))
		if err != nil {
			return err
		}
	}

	if IsPrivateChannel(channel) {
		err = DelUserInfoFields(s.socketID, UserInfoAuthKey(channel))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) subscribeSuccessCallback(channel string) error {
	s.subscriptions[channel] = true
	err := AddChannelAndSubscription(channel, s.socketID)
	if err != nil {
		return err
	}

	if ChannelSubscriberCount(channel) == 1 {
		hookEvent := HookEvent{Name: "channel_occupied", Channel: channel}
		go TriggerHook(&hookEvent)
	}

	if IsPresenceChannel(channel) {
		var presenceData PresenceData
		presenceData.Hash = make(map[string]map[string]string)
		subscriptions := GetSubscriptions(channel)

		for _, socketID := range subscriptions {
			presenceChannelData, err := GetPresenceChannelData(socketID, channel)
			if err != nil {
				return err
			}
			presenceData.IDs = append(presenceData.IDs, presenceChannelData.UserID)
			presenceData.Hash[presenceChannelData.UserID] = presenceChannelData.UserInfo
		}
		presenceData.Count = len(subscriptions)

		b, err := json.Marshal(map[string]PresenceData{"presence": presenceData})
		if err != nil {
			return err
		}
		s.Send(SubscriptionSucceededPack(channel, string(b[:])))

		err = PublishPresenceMemberAddedEvent(s.socketID, channel)
		if err != nil {
			return err
		}
	} else {
		s.Send(SubscriptionSucceededPack(channel, "{}"))
	}
	return nil
}

func (s *Session) subscribePump() {
	defer func() {
		s.hub.unregister <- s
		_ = s.stop(4100)
	}()
	for {
		msgI, err := s.pubSub.Receive()
		if err != nil {
			log.Error(err)
			return
		}
		switch msg := msgI.(type) {
		case *redis.Subscription:
			channel := msg.Channel
			log.Infoln(s.socketID, msg.Kind, "success to", channel)

			switch msg.Kind {
			case "subscribe":
				if s.subscriptions[channel] {
					continue
				}
				err = s.subscribeSuccessCallback(channel)
				if err != nil {
					log.Error(err)
					return
				}
			case "unsubscribe":
				if !s.subscriptions[channel] {
					continue
				}
				err = s.unsubscribeSuccessCallback(channel)
				if err != nil {
					log.Error(err)
					return
				}
			default:
				log.Error("unhandled subscription kind", msg.Kind)
			}
		case *redis.Message:
			log.Infoln(s.socketID, "received", msg.Payload, "from", msg.Channel)

			var eventForm forms.EventForm
			err := json.Unmarshal([]byte(msg.Payload), &eventForm)
			if err != nil {
				log.Error(err)
				return
			}

			// discard message
			if eventForm.SocketID == s.socketID {
				continue
			}

			channelEvent := ChannelEvent{Event: eventForm.Name, Channel: msg.Channel, Data: eventForm.Data}
			if IsPresenceChannel(channelEvent.Channel) && eventForm.UserID != "" {
				channelEvent.UserID = eventForm.UserID
			}

			b, err := json.Marshal(channelEvent)
			if err != nil {
				log.Error(err)
				return
			}
			s.Send(b)
		default:
			log.Error("unhandled", msg)
		}
	}
}

func (s *Session) readPump() {
	defer func() {
		s.hub.unregister <- s
	}()
	s.conn.SetReadLimit(maxMessageSize)
	s.resetReadTimeout()
	for {
		_, msgRaw, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Error(err)
				_ = s.stop(4100)
			} else {
				_ = s.stop(4202)
			}
			break
		}
		err = s.ParsePayload(msgRaw)
		if err != nil {
			log.Error(err)
			_ = s.stop(4200)
		}
	}
}

func (s *Session) writePump() {
	ticker := time.NewTicker(activityTimeout)
	defer func() {
		ticker.Stop()
		s.hub.unregister <- s
		_ = s.stop(4100)
	}()
	for {
		select {
		case message, ok := <-s.send:
			if !ok {
				return
			}
			err := s.write(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			s.Send(PingPongPack(PusherPing))
			s.resetPongTimeout()
		}
	}
}
