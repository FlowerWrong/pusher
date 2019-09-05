package pusher

import "errors"

// @doc https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol#error-codes
var (
	ErrCodes = map[int]string{
		4000: "Application only accepts SSL connections, reconnect using wss://",
		4001: "Application does not exist",
		4003: "Application disabled",
		4004: "Application is over connection quota",
		4005: "Path not found",
		4006: "Invalid version string format",
		4007: "Unsupported protocol version",
		4008: "No protocol version supplied",
		4009: "Connection is unauthorized",
		4100: "Over capacity",
		4200: "Generic reconnect immediately",
		4201: "Pong reply not received: ping was sent to the client, but no reply was received - see ping and pong messages",
		4202: "Closed after inactivity: Client has been inactive for a long time (currently 24 hours) and client does not support ping. Please upgrade to a newer WebSocket draft or implement version 5 or above of this protocol",
		4301: "Client event rejected due to rate limit",
	}

	ErrInvalidChannel                 = errors.New("channel's name are invalid")
	ErrPresenceSubscriberTooMuch      = errors.New("presence channel limit to 100 members maximum")
	ErrPresenceUserDataTooMuch        = errors.New("presence channel limit to 1KB user object")
	ErrPresenceUserIDTooLong          = errors.New("presence channel limit to 128 characters user id")
	ErrAPIReqEventDataTooLarge        = errors.New("request event data too large, limit to 10kb")
	ErrAPIReqEventChannelsSizeTooLong = errors.New("request event channels size too big, limit to 100")
)
