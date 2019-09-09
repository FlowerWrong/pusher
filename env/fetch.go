package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Lookup ...
func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Get ...
func Get(key string, fallback ...string) string {
	var fb string
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	return s
}

// GetString ...
func GetString(key string, fallback ...string) string {
	return Get(key, fallback...)
}

// GetInt ...
func GetInt(key string, fallback ...int) int {
	var fb int
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	i, err := parseInt(s)
	if err != nil {
		return fb
	}

	return int(i)
}

// GetInt64 ...
func GetInt64(key string, fallback ...int64) int64 {
	var fb int64
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	i, err := parseInt64(s)
	if err != nil {
		return fb
	}

	return int64(i)
}

// GetUint ...
func GetUint(key string, fallback ...uint) uint {
	var fb uint
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	i, err := parseUint(s)
	if err != nil {
		return fb
	}
	return uint(i)
}

// GetFloat ...
func GetFloat(key string, fallback ...float64) float64 {
	var fb float64
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fb
	}
	return n
}

// GetDuration ...
func GetDuration(key string, fallback ...time.Duration) time.Duration {
	var fb time.Duration
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fb
	}
	return d
}

// GetBool ...
func GetBool(key string, fallback ...bool) bool {
	var fb bool
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	s, ok := Lookup(key)
	if !ok {
		return fb
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return fb
	}
	return b
}

// parse an int, falling back to parsing it as a float
func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(i), nil
	}

	// Try to parse as float, then convert
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %v", s)
	}
	return int(n), nil
}

// parse an int64, falling back to parsing it as a float
func parseInt64(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return int(i), nil
	}

	// Try to parse as float, then convert
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %v", s)
	}
	return int(n), nil
}

// parse an int, falling back to parsing it as a float
func parseUint(s string) (uint, error) {
	if i, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint(i), nil
	}

	// Try to parse as float, then convert
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		if f < 0 {
			return 0, fmt.Errorf("less than zero: %s", s)
		}
		return uint(f), nil
	}
	return 0, fmt.Errorf("invalid int: %v", s)
}
