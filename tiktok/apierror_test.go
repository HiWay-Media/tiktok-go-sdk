package tiktok

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// bodyServer answers every request with the given status and body.
func bodyServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// TestErrorInBandIsDetected is the regression test for the defect this package
// shipped with: TikTok answers HTTP 200 and puts the failure in the body, so
// every method used to return (empty response, nil error) — a refused post was
// indistinguishable from a published one.
func TestErrorInBandIsDetected(t *testing.T) {
	const body = `{"data":{},"error":{"code":"spam_risk_too_many_posts","message":"daily post cap reached","log_id":"2026072812"}}`
	srv := bodyServer(t, http.StatusOK, body)
	c := newTestClient(t, srv)

	_, err := c.PostVideoInit("t", "d", "https://example.com/v.mp4", string(SELF_ONLY), false, false, false)
	if err == nil {
		t.Fatal("HTTP 200 with an error envelope must produce an error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T, want *APIError", err)
	}
	if apiErr.Code != "spam_risk_too_many_posts" {
		t.Errorf("Code = %q", apiErr.Code)
	}
	if apiErr.Message != "daily post cap reached" {
		t.Errorf("Message = %q", apiErr.Message)
	}
	if apiErr.LogID != "2026072812" {
		t.Errorf("LogID = %q — it is the only handle TikTok support can trace", apiErr.LogID)
	}
	if apiErr.HTTPStatus != http.StatusOK {
		t.Errorf("HTTPStatus = %d, want 200 (that is the whole point)", apiErr.HTTPStatus)
	}
	// The operation must name the call that failed: every method used to report
	// "post video init error", including the photo and the video-list ones.
	if apiErr.Op != "post video init" {
		t.Errorf("Op = %q", apiErr.Op)
	}
}

// TestOAuthErrorEnvelope covers the second shape: the token endpoint does not
// use ErrorObject, it answers 200 with a flat envelope. Verified against the
// real API: empty credentials return exactly this.
func TestOAuthErrorEnvelope(t *testing.T) {
	const body = `{"error":"invalid_request","error_description":"The request parameters are malformed.","log_id":"20260728055901"}`
	srv := bodyServer(t, http.StatusOK, body)
	c := newTestClient(t, srv)

	tok, err := c.GetClientAccessTokenManagement()
	if err == nil {
		t.Fatalf("want an error, got token %+v", tok)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T, want *APIError", err)
	}
	if apiErr.Code != "invalid_request" {
		t.Errorf("Code = %q", apiErr.Code)
	}
	if !strings.Contains(apiErr.Message, "malformed") {
		t.Errorf("Message = %q, want the error_description", apiErr.Message)
	}
	if apiErr.LogID != "20260728055901" {
		t.Errorf("LogID = %q", apiErr.LogID)
	}
}

// TestSuccessEnvelopeIsNotAnError is the guard against the opposite mistake:
// it must pass both with and without the in-band error check, so seeing it
// green while the tests above are red is the expected result of the
// counter-proof, not a hole in it.
func TestSuccessEnvelopeIsNotAnError(t *testing.T) {
	srv := bodyServer(t, http.StatusOK, `{"data":{"publish_id":"pid"},"error":{"code":"ok","message":"","log_id":"x"}}`)
	c := newTestClient(t, srv)
	if _, err := c.PublishVideo("pid"); err != nil {
		t.Fatalf(`code "ok" must not be an error: %v`, err)
	}

	// A body with no error envelope at all (how the token endpoint answers a
	// successful request) must not be an error either.
	srv2 := bodyServer(t, http.StatusOK, `{"access_token":"clt.token","token_type":"Bearer","expires_in":7200}`)
	c2 := newTestClient(t, srv2)
	tok, err := c2.GetClientAccessTokenManagement()
	if err != nil {
		t.Fatalf("valid token response: %v", err)
	}
	if tok.AccessToken != "clt.token" {
		t.Errorf("AccessToken = %q", tok.AccessToken)
	}
}

func TestHTTPErrorWithUnparsableBody(t *testing.T) {
	srv := bodyServer(t, http.StatusBadGateway, "<html>502 Bad Gateway</html>")
	c := newTestClient(t, srv)

	_, err := c.UserInfo()
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T (%v), want *APIError", err, err)
	}
	if apiErr.HTTPStatus != http.StatusBadGateway {
		t.Errorf("HTTPStatus = %d, want 502", apiErr.HTTPStatus)
	}
	if apiErr.Op != "user info" {
		t.Errorf("Op = %q, want the endpoint that failed", apiErr.Op)
	}
	if !strings.Contains(err.Error(), "http 502") {
		t.Errorf("Error() = %q, want the status when there is no code", err.Error())
	}
}

func TestExcerptBoundsTheBody(t *testing.T) {
	// An HTML error page must not end up whole inside an error string (and from
	// there into a log line).
	long := strings.Repeat("x", 1000)
	if got := excerpt(long); len(got) > 300 {
		t.Errorf("excerpt kept %d chars", len(got))
	}
}

func TestErrorObjectIsError(t *testing.T) {
	cases := []struct {
		code string
		want bool
	}{
		{"ok", false},
		{"", false}, // some endpoints omit the code on success
		{"access_token_invalid", true},
	}
	for _, tc := range cases {
		if got := (ErrorObject{Code: tc.code}).IsError(); got != tc.want {
			t.Errorf("ErrorObject{Code: %q}.IsError() = %v, want %v", tc.code, got, tc.want)
		}
	}
}

func TestAPIErrorMessage(t *testing.T) {
	e := &APIError{Op: "post video init", Code: "spam_risk", Message: "too many posts", LogID: "abc", HTTPStatus: 200}
	want := "tiktok: post video init: spam_risk: too many posts (log_id: abc)"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestInputValidationErrors(t *testing.T) {
	var seen []string
	srv := okServer(t, &seen, nil)
	c := newTestClient(t, srv)

	if _, err := c.PostVideoInit("t", "d", "u", "NOPE", false, false, false); !errors.Is(err, ErrPrivacyLevel) {
		t.Errorf("PostVideoInit: err = %v, want ErrPrivacyLevel", err)
	}
	if _, err := c.PostPhotoInit("t", "d", "NOPE", nil, string(DIRECT_POST)); !errors.Is(err, ErrPrivacyLevel) {
		t.Errorf("PostPhotoInit privacy: err = %v, want ErrPrivacyLevel", err)
	}
	if _, err := c.PostPhotoInit("t", "d", string(SELF_ONLY), nil, "NOPE"); !errors.Is(err, ErrPostMode) {
		t.Errorf("PostPhotoInit mode: err = %v, want ErrPostMode", err)
	}
	// The deprecated names must keep matching: they are the same values.
	if _, err := c.PostVideoInit("t", "d", "u", "NOPE", false, false, false); !errors.Is(err, PrivacyLevelWrong) {
		t.Error("the deprecated alias PrivacyLevelWrong must still match")
	}
	// Validation happens before any HTTP call.
	if len(seen) != 0 {
		t.Errorf("invalid input reached the network: %v", seen)
	}
}
