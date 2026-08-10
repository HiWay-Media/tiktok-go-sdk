package tiktok

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// APIError is a failure reported by TikTok.
//
// It covers both ways the API can fail, because it can fail in both while
// answering HTTP 200:
//
//   - the /v2/... endpoints wrap the outcome in an "error" object
//     (ErrorObject: code, message, log_id) and use the code "ok" for success;
//   - /v2/oauth/token/ answers with a flat envelope instead
//     ("error", "error_description", "log_id").
//
// Look at Code to tell the cases apart, and always quote LogID when opening a
// ticket with TikTok: it is the only handle they can trace.
type APIError struct {
	Code       string // "access_token_invalid", "invalid_request", ...
	Message    string // human readable description, as sent by TikTok
	LogID      string // TikTok's request id, for support
	HTTPStatus int    // HTTP status of the response that carried the error
	Op         string // SDK operation that failed, e.g. "post video init"
}

func (e *APIError) Error() string {
	var b strings.Builder
	b.WriteString("tiktok: ")
	if e.Op != "" {
		b.WriteString(e.Op)
		b.WriteString(": ")
	}
	if e.Code != "" {
		b.WriteString(e.Code)
	} else {
		fmt.Fprintf(&b, "http %d", e.HTTPStatus)
	}
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	if e.LogID != "" {
		fmt.Fprintf(&b, " (log_id: %s)", e.LogID)
	}
	return b.String()
}

// errorEnvelope models the two shapes at once. "error" is an object for the
// /v2/... endpoints and a string for the OAuth token endpoint, so it is kept
// raw and decoded twice.
type errorEnvelope struct {
	Error            json.RawMessage `json:"error"`
	ErrorDescription string          `json:"error_description"`
	LogId            string          `json:"log_id"`
}

// checkResponse turns a TikTok response into an error, or nil when the call
// actually succeeded.
//
// Checking resp.IsError() alone is not enough and never was: TikTok answers
// HTTP 200 to a rejected request and puts the outcome in the body, so a caller
// that only looks at the status code cannot tell a published post from a
// refused one.
func checkResponse(resp *resty.Response, op string) error {
	apiErr := parseAPIError(resp.Body())
	if apiErr != nil {
		apiErr.Op = op
		apiErr.HTTPStatus = resp.StatusCode()
		return apiErr
	}
	if resp.IsError() {
		// A transport-level or gateway failure with a body we cannot parse
		// (HTML error page, empty body): keep the status and a bounded excerpt.
		return &APIError{
			Op:         op,
			HTTPStatus: resp.StatusCode(),
			Message:    excerpt(resp.String()),
		}
	}
	return nil
}

// parseAPIError extracts the error envelope from a response body, or returns
// nil when the body reports success (or carries no envelope at all, which is
// how the token endpoint answers a successful request).
func parseAPIError(body []byte) *APIError {
	var env errorEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil // not JSON we understand; the caller falls back to the status
	}
	if len(env.Error) == 0 || string(env.Error) == "null" {
		return nil
	}

	// Shape 1: {"error": {"code": "...", "message": "...", "log_id": "..."}}
	var obj ErrorObject
	if err := json.Unmarshal(env.Error, &obj); err == nil {
		if !obj.IsError() {
			return nil
		}
		return &APIError{Code: obj.Code, Message: obj.Message, LogID: obj.LogId}
	}

	// Shape 2: {"error": "invalid_request", "error_description": "...", ...}
	var code string
	if err := json.Unmarshal(env.Error, &code); err == nil {
		oauthErr := OAuthErrorObject{Error: code, ErrorDescription: env.ErrorDescription, LogId: env.LogId}
		if !oauthErr.IsError() {
			return nil
		}
		return &APIError{Code: oauthErr.Error, Message: oauthErr.ErrorDescription, LogID: oauthErr.LogId}
	}
	return nil
}

// excerpt bounds an unparsed body so an HTML error page cannot end up whole
// inside an error message (or a log line).
func excerpt(s string) string {
	s = strings.TrimSpace(s)
	const max = 256
	if len(s) > max {
		return s[:max] + "…"
	}
	return s
}
