package tiktok

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// newTestClient builds a client pointed at srv. It returns the concrete type so
// the tests can look at the transport the options produced.
func newTestClient(t *testing.T, srv *httptest.Server, opts ...Option) *tiktok {
	t.Helper()
	opts = append([]Option{WithBaseURL(srv.URL)}, opts...)
	c, err := NewTikTok("test-key", "test-secret", false, opts...)
	if err != nil {
		t.Fatalf("NewTikTok: %v", err)
	}
	return c.(*tiktok)
}

// okServer answers every request with an empty-but-valid TikTok envelope and
// records the path and headers it saw.
func okServer(t *testing.T, seen *[]string, hdr *http.Header) *httptest.Server {
	t.Helper()
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		*seen = append(*seen, r.URL.Path)
		if hdr != nil {
			*hdr = r.Header.Clone()
		}
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{},"error":{"code":"ok","message":"","log_id":"test"}}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestNewTikTokRequiresCredentials(t *testing.T) {
	// Empty credentials used to produce a working-looking client whose every
	// call was rejected by TikTok much later, with an opaque message.
	if _, err := NewTikTok("", "secret", false); !errors.Is(err, ErrClientKeyRequired) {
		t.Errorf("empty key: err = %v, want ErrClientKeyRequired", err)
	}
	if _, err := NewTikTok("key", "", false); !errors.Is(err, ErrClientSecretRequired) {
		t.Errorf("empty secret: err = %v, want ErrClientSecretRequired", err)
	}
	if _, err := NewTikTok("key", "secret", false); err != nil {
		t.Errorf("valid credentials: err = %v, want nil", err)
	}
}

func TestDefaultTimeout(t *testing.T) {
	c, err := NewTikTok("key", "secret", false)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.(*tiktok).restClient.GetClient().Timeout; got != DefaultTimeout {
		t.Errorf("default timeout = %v, want %v", got, DefaultTimeout)
	}

	c, _ = NewTikTok("key", "secret", false, WithTimeout(2*time.Second))
	if got := c.(*tiktok).restClient.GetClient().Timeout; got != 2*time.Second {
		t.Errorf("WithTimeout: timeout = %v, want 2s", got)
	}

	// Disabling the timeout must stay possible, but only on purpose.
	c, _ = NewTikTok("key", "secret", false, WithTimeout(0))
	if got := c.(*tiktok).restClient.GetClient().Timeout; got != 0 {
		t.Errorf("WithTimeout(0): timeout = %v, want 0 (disabled)", got)
	}
}

func TestTimeoutIsEnforced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	}))
	defer srv.Close()

	c := newTestClient(t, srv, WithTimeout(50*time.Millisecond))
	start := time.Now()
	if _, err := c.CreatorInfo(); err == nil {
		t.Fatal("want a timeout error, got nil")
	}
	if elapsed := time.Since(start); elapsed > 250*time.Millisecond {
		t.Errorf("call took %v: the timeout did not fire", elapsed)
	}
}

func TestWithHTTPClient(t *testing.T) {
	var seen []string
	srv := okServer(t, &seen, nil)

	own := &http.Client{}
	c := newTestClient(t, srv, WithHTTPClient(own))
	if c.restClient.GetClient() != own {
		t.Fatal("WithHTTPClient: the SDK is not using the provided *http.Client")
	}
	// The timeout is applied on top of the caller's client, not instead of it.
	if own.Timeout != DefaultTimeout {
		t.Errorf("timeout on the provided client = %v, want %v", own.Timeout, DefaultTimeout)
	}
}

// TestEndpointsUseTheBaseURL is the test that makes every other endpoint test
// possible: if a method hardcodes an absolute URL, resty ignores the base URL
// and the call escapes to the real API.
func TestEndpointsUseTheBaseURL(t *testing.T) {
	var seen []string
	srv := okServer(t, &seen, nil)
	c := newTestClient(t, srv)
	c.SetAccessToken("act.token")

	calls := []struct {
		name string
		want string
		call func() error
	}{
		{"GetClientAccessTokenManagement", OAUTH_TOKEN, func() error {
			_, err := c.GetClientAccessTokenManagement()
			return err
		}},
		{"CreatorInfo", QUERY_CREATOR_INFO, func() error { _, err := c.CreatorInfo(); return err }},
		{"PostVideoInit", POST_PUBLISH_VIDEO_INIT, func() error {
			_, err := c.PostVideoInit("t", "d", "https://example.com/v.mp4", string(SELF_ONLY), false, false, false)
			return err
		}},
		{"PublishVideo", PUBLISH_STATUS_FETCH, func() error { _, err := c.PublishVideo("pid"); return err }},
		{"PostPhotoInit", POST_PUBLISH_CONTENT_INIT, func() error {
			_, err := c.PostPhotoInit("t", "d", string(SELF_ONLY), []string{"https://example.com/a.webp"}, string(MEDIA_UPLOAD))
			return err
		}},
		{"GetVideoList", VIDEO_LIST, func() error { _, err := c.GetVideoList(20); return err }},
		{"UserInfo", USER_INFO, func() error { _, err := c.UserInfo(); return err }},
	}
	for _, tc := range calls {
		seen = seen[:0]
		if err := tc.call(); err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		if len(seen) != 1 || seen[0] != tc.want {
			t.Errorf("%s hit %v, want [%s]", tc.name, seen, tc.want)
		}
	}
}

func TestAccessTokenIsSentAsBearer(t *testing.T) {
	var seen []string
	var hdr http.Header
	srv := okServer(t, &seen, &hdr)

	c := newTestClient(t, srv)
	c.SetAccessToken("act.example")
	if _, err := c.CreatorInfo(); err != nil {
		t.Fatal(err)
	}
	if got := hdr.Get("Authorization"); got != "Bearer act.example" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer act.example")
	}

	// A token set at construction time must reach the wire the same way.
	c2 := newTestClient(t, srv, WithAccessToken("act.from-option"))
	if _, err := c2.CreatorInfo(); err != nil {
		t.Fatal(err)
	}
	if got := hdr.Get("Authorization"); got != "Bearer act.from-option" {
		t.Errorf("WithAccessToken: Authorization = %q", got)
	}
}

// TestAccessTokenConcurrentAccess is meaningful under -race: the token used to
// be a plain field written by SetAccessToken while requests read it.
func TestAccessTokenConcurrentAccess(t *testing.T) {
	var seen []string
	srv := okServer(t, &seen, nil)
	c := newTestClient(t, srv)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				c.SetAccessToken("act.rotating")
				_ = c.GetAccessToken()
				if _, err := c.CreatorInfo(); err != nil {
					t.Errorf("CreatorInfo: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}
