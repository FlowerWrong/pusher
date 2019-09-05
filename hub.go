package pusher

import "github.com/FlowerWrong/pusher/log"

// Hub maintains the set of active sessions and broadcasts messages to the sessions.
type Hub struct {
	// Registered sessions.
	sessions map[*Session]bool

	// Inbound messages from the sessions.
	broadcast chan []byte

	// Register requests from the sessions.
	register chan *Session

	// Unregister requests from sessions.
	unregister chan *Session
}

// NewHub builds new hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		sessions:   make(map[*Session]bool),
	}
}

// Run hub
func (h *Hub) Run() {
	for {
		select {
		case session := <-h.register:
			h.sessions[session] = true
			err := AddUser(session.socketID)
			if err != nil {
				log.Error(err)
			}
			log.Infoln(session.socketID, "joined")
		case session := <-h.unregister:
			h.cleanSession(session)
		case message := <-h.broadcast:
			for session := range h.sessions {
				session.Send(message)
			}
		}
	}
}

func (h *Hub) cleanSession(session *Session) {
	if _, ok := h.sessions[session]; ok {
		delete(h.sessions, session)
	}
}
