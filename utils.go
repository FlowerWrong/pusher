package pusher

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	channelValidationRegex = regexp.MustCompile("^[-a-zA-Z0-9_=@,.;]+$")
	validAppIDRegex        = regexp.MustCompile(`^[0-9]+$`)
	maxChannelNameSize     = 164
)

// GenerateSocketID generate a new random Hash
func GenerateSocketID() string {
	return fmt.Sprintf("%d.%d", rand.Intn(math.MaxInt32), rand.Intn(math.MaxInt32))
}

// IsPresenceChannel ...
func IsPresenceChannel(channel string) bool {
	return strings.HasPrefix(channel, "presence-")
}

// IsPrivateChannel ...
func IsPrivateChannel(channel string) bool {
	return strings.HasPrefix(channel, "private-")
}

// IsClientEvent ...
func IsClientEvent(event string) bool {
	return strings.HasPrefix(event, "client-")
}

// ValidChannel ...
func ValidChannel(channel string) bool {
	if len(channel) > maxChannelNameSize || !channelValidationRegex.MatchString(channel) {
		return false
	}
	return true
}

// ValidAppID ...
func ValidAppID(appID string) bool {
	if !validAppIDRegex.MatchString(appID) {
		return false
	}
	return true
}

// Str2Int64 string -> int64
func Str2Int64(a string) (int64, error) {
	b, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return 0, err
	}
	return b, nil
}

// Str2Int string -> int
func Str2Int(a string) (int, error) {
	b, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return b, nil
}
