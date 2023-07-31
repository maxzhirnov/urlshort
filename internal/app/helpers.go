package app

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func generateID(length int) string {

	// minimum id length is 4 symbols
	// TODO: think if it's better to return error here
	if length < 4 {
		length = 4
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CheckURL verifies if the provided string is a valid URL.
// If valid, it returns true and the possibly modified URL string.
func CheckURL(s string) (urlWithoutProtocol string, isValid bool) {
	if s == "localhost" || s == "http://localhost" || s == "https://localhost" {
		return "localhost", true
	}

	if !strings.HasPrefix(s, "http") {
		s = "http://" + s
	}

	u, err := url.ParseRequestURI(s)
	if err != nil || !strings.Contains(u.Host, ".") {
		return "", false
	}

	return u.Host, true
}
