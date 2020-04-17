package cacheheaders

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

// X-Cache-Channel: Used by Varnish
const CacheChannelHeader = "X-Cache-Channel"

// Cache-Tag: Used by Cloudflare (Enterprise only)
const CacheTagHeader = "Cache-Tag"

// CacheChannels - A middleware struct that outputs cache channel / cache tag headers
// These headers can be used by cache proxies to ban / invalidate large groups of cached items in one go
type CacheChannels struct {
	Varnish    bool // if true, SendHeaders() will send the Varnish "X-Cache-Channel" header
	Cloudflare bool // if true, SendHeaders() will send the Cloudflare Enterprise "Cache-Tag" header
	channels   []string
}

// Add - Prunes, then adds the specified channels to the channel slice
// Removes any character that's not in the legal range: a-z, A-Z, 0-9, _, -
func (cc *CacheChannels) Add(channels ...string) {
	rx, err := regexp.Compile(`[^\w-]+`)
	if err != nil { // shouldn't ever happen (but still)
		log.Printf("Regex compile error: %s", err.Error())
		return
	}

	for _, ch := range channels {
		cc.channels = append(cc.channels, rx.ReplaceAllString(ch, ""))
	}
}

// Set - Replaces any existing channels, with the specified channels
func (cc *CacheChannels) Set(channels ...string) {
	cc.channels = []string{}
	cc.Add(channels...)
}

// SendHeaders - A middleware function compatible with most routers
// Sends the configured cache channels as "X-Cache-Channel" headers
func (cc *CacheChannels) SendHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cc.Varnish {
			w.Header().Set(CacheChannelHeader, strings.Join(cc.channels, ", "))
		}

		if cc.Cloudflare {
			w.Header().Set(CacheTagHeader, strings.Join(cc.channels, ","))
		}

		next.ServeHTTP(w, r)
	})
}
