package tiktok

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// captureServer records the JSON body of the last request it served.
func captureServer(t *testing.T, body *map[string]interface{}) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var decoded map[string]interface{}
		_ = json.Unmarshal(raw, &decoded)
		*body = decoded
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{},"error":{"code":"ok"}}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func postInfo(t *testing.T, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	pi, ok := body["post_info"].(map[string]interface{})
	if !ok {
		t.Fatalf("no post_info in request body: %v", body)
	}
	return pi
}

func sourceInfo(t *testing.T, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	si, ok := body["source_info"].(map[string]interface{})
	if !ok {
		t.Fatalf("no source_info in request body: %v", body)
	}
	return si
}

func TestPostPhotoUsesTheCallersOptions(t *testing.T) {
	var body map[string]interface{}
	c := newTestClient(t, captureServer(t, &body))

	_, err := c.PostPhoto(PhotoPost{
		Title:           "t",
		PrivacyLevel:    string(SELF_ONLY),
		PhotoUrls:       []string{"https://example.com/a.webp", "https://example.com/b.webp"},
		PostMode:        string(DIRECT_POST),
		DisableComment:  true,
		AutoAddMusic:    false,
		PhotoCoverIndex: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	pi := postInfo(t, body)
	if pi["disable_comment"] != true {
		t.Errorf("disable_comment = %v, want true — it used to be hardcoded false", pi["disable_comment"])
	}
	if pi["auto_add_music"] != false {
		t.Errorf("auto_add_music = %v, want false — it used to be hardcoded true", pi["auto_add_music"])
	}
	if got := sourceInfo(t, body)["photo_cover_index"]; got != float64(2) {
		t.Errorf("photo_cover_index = %v, want 2 — it used to be hardcoded 1", got)
	}
}

// TestPostPhotoInitKeepsItsDefaults pins the behaviour of the old entry point:
// callers who never chose these values must keep getting what they had.
func TestPostPhotoInitKeepsItsDefaults(t *testing.T) {
	var body map[string]interface{}
	c := newTestClient(t, captureServer(t, &body))

	if _, err := c.PostPhotoInit("t", "d", string(SELF_ONLY), []string{"https://example.com/a.webp"}, string(MEDIA_UPLOAD)); err != nil {
		t.Fatal(err)
	}
	pi := postInfo(t, body)
	if pi["auto_add_music"] != true || pi["disable_comment"] != false {
		t.Errorf("legacy defaults changed: %v", pi)
	}
	if got := sourceInfo(t, body)["photo_cover_index"]; got != float64(1) {
		t.Errorf("photo_cover_index = %v, want the historical 1", got)
	}
}

func TestPostVideoSendsTheCoverTimestamp(t *testing.T) {
	var body map[string]interface{}
	c := newTestClient(t, captureServer(t, &body))

	_, err := c.PostVideo(VideoPost{
		Title:                 "t",
		VideoUrl:              "https://example.com/v.mp4",
		PrivacyLevel:          string(PUBLIC_TO_EVERYONE),
		VideoCoverTimestampMS: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	pi := postInfo(t, body)
	// The field was in the struct but no entry point could ever set it.
	if got := pi["video_cover_timestamp_ms"]; got != float64(1000) {
		t.Errorf("video_cover_timestamp_ms = %v, want 1000", got)
	}
	// description is not a documented post_info field for a video: an empty one
	// must not be sent at all.
	if _, present := pi["description"]; present {
		t.Errorf("an empty description reached the wire: %v", pi)
	}
}

func TestPostRequiresItsMedia(t *testing.T) {
	var seen []string
	srv := okServer(t, &seen, nil)
	c := newTestClient(t, srv)

	if _, err := c.PostPhoto(PhotoPost{PrivacyLevel: string(SELF_ONLY), PostMode: string(DIRECT_POST)}); !errors.Is(err, ErrPhotoUrlsRequired) {
		t.Errorf("photo post with no image: err = %v, want ErrPhotoUrlsRequired", err)
	}
	if _, err := c.PostVideo(VideoPost{PrivacyLevel: string(SELF_ONLY)}); !errors.Is(err, ErrVideoUrlRequired) {
		t.Errorf("video post with no URL: err = %v, want ErrVideoUrlRequired", err)
	}
	if len(seen) != 0 {
		t.Errorf("an incomplete post reached the network: %v", seen)
	}
}
