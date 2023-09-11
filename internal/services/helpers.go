package services

import (
	"net/url"
	"strings"
)

// CheckURL verifies if the provided string is a valid URL.
// If valid, it returns true and the possibly modified URL string.
func CheckURL(s string) (urlWithoutProtocol string, isValid bool) {
	if s == "localhost" || s == "http://localhost" || s == "https://localhost" {
		return "localhost", true
	}

	s = EnsureURLScheme(s)

	u, err := url.ParseRequestURI(s)
	if err != nil || !strings.Contains(u.Host, ".") {
		return "", false
	}

	return u.Host, true
}

func EnsureURLScheme(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		// Дефолт на http если схемы нет, чтобы не было ошибок если сайт не поддерживает https
		return "http://" + url
	}
	return url
}
