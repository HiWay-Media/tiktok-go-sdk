package tiktok

import (
	"net/http"
	"time"
)

// DefaultTimeout is the per-request timeout applied when the caller does not
// choose one. resty (like net/http) defaults to no timeout at all, which for a
// library called from a server means a single unlucky request can pin a
// goroutine forever.
const DefaultTimeout = 30 * time.Second

// Option customizes a client. Options are applied by NewTikTok after the
// defaults and before the underlying HTTP client is built, so they can replace
// the transport as well as tune it.
type Option func(*tiktok)

// WithBaseURL points the client at another host, e.g. an httptest server in a
// test. In production the default (BASE_URL) is what you want.
func WithBaseURL(baseURL string) Option {
	return func(o *tiktok) { o.baseURL = baseURL }
}

// WithTimeout sets the per-request timeout. A non-positive duration disables
// the timeout entirely — an explicit choice, not the default one.
func WithTimeout(d time.Duration) Option {
	return func(o *tiktok) { o.timeout = d }
}

// WithHTTPClient makes the SDK use the given *http.Client, for callers that
// already have one configured (proxy, custom transport, connection pool,
// instrumentation). WithTimeout still applies on top of it.
func WithHTTPClient(c *http.Client) Option {
	return func(o *tiktok) { o.httpClient = c }
}

// WithDebug turns request/response logging on or off.
//
// Debug logging prints whole objects, access token included: use it while
// developing, never in a shared or production log.
func WithDebug(debug bool) Option {
	return func(o *tiktok) { o.debug = debug }
}

// WithAccessToken sets the user access token at construction time, so a caller
// that already has one does not have to build the client and then mutate it.
func WithAccessToken(token string) Option {
	return func(o *tiktok) { o.accessToken = token }
}
