package util

import (
	"net/http"

	"github.com/tomasen/realip"
)

// GetClientIP tries to get the client's real IP even behind proxies
func GetClientIP(r *http.Request) string {
	return realip.FromRequest(r)
}
