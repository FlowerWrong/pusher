package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Verify signature
// @doc https://pusher.com/docs/channels/library_auth_reference/rest-api#authentication
func Verify(r *http.Request) (bool, error) {
	query := r.URL.Query()

	if !checkVersion(query.Get("auth_version")) {
		return false, errors.New("invalid version")
	}

	timestamp, _ := Str2Int64(query.Get("auth_timestamp"))
	if !checkTimestamp(timestamp, viper.GetInt64("TIMESTAMP_GRACE")) {
		return false, errors.New("invalid timestamp")
	}

	signature := query.Get("auth_signature")
	query.Del("auth_signature")
	queryString := prepareQueryString(query)
	stringToSign := strings.Join([]string{strings.ToUpper(r.Method), r.URL.Path, queryString}, "\n")
	return HmacSignature(stringToSign, viper.GetString("APP_SECRET")) == signature, nil
}

// HmacSignature ...
func HmacSignature(toSign, secret string) string {
	return hex.EncodeToString(hmacBytes([]byte(toSign), []byte(secret)))
}

func checkVersion(version string) bool {
	return version == "1.0"
}

func checkTimestamp(timestamp, grace int64) bool {
	if (time.Now().Unix() - timestamp) >= grace {
		return false
	}
	return true
}

func checkBodyMD5(toSign []byte, md5Str string) bool {
	return md5Hex(toSign) == md5Str
}

func prepareQueryString(params url.Values) string {
	var keys []string
	for key := range params {
		keys = append(keys, strings.ToLower(key))
	}

	sort.Strings(keys)
	var pieces []string
	for _, key := range keys {
		pieces = append(pieces, key+"="+params.Get(key))
	}

	return strings.Join(pieces, "&")
}

func hmacBytes(toSign, secret []byte) []byte {
	_authSignature := hmac.New(sha256.New, secret)
	_authSignature.Write(toSign)
	return _authSignature.Sum(nil)
}

func md5Hex(toSign []byte) string {
	h := md5.New()
	h.Write(toSign)
	return hex.EncodeToString(h.Sum(nil))
}
