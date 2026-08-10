package tiktok

import (
	"encoding/json"
	"testing"
)

func TestDataPublishVideoKeepsBothSpellings(t *testing.T) {
	var d DataPublishVideo
	if err := json.Unmarshal([]byte(`{"publish_id":"v_pub_url~v2.123"}`), &d); err != nil {
		t.Fatal(err)
	}
	if d.PublishId != "v_pub_url~v2.123" {
		t.Errorf("PublishId = %q", d.PublishId)
	}
	// The misspelled field shipped as public API: code that reads it must keep
	// working, or the fix is a silent breakage instead of a rename.
	if d.PubblishId != d.PublishId {
		t.Errorf("deprecated PubblishId = %q, want %q", d.PubblishId, d.PublishId)
	}
}

func TestPublishStatusFetchAcceptsEitherSpelling(t *testing.T) {
	cases := map[string]string{
		"tiktok's spelling": `{"status":"PUBLISH_COMPLETE","publicaly_available_post_id":[123,456]}`,
		"correct spelling":  `{"status":"PUBLISH_COMPLETE","publicly_available_post_id":[123,456]}`,
	}
	for name, body := range cases {
		var p PublishStatusFetch
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if len(p.PubliclyAvailablePostId) != 2 || p.PubliclyAvailablePostId[0] != 123 {
			t.Errorf("%s: PubliclyAvailablePostId = %v", name, p.PubliclyAvailablePostId)
		}
		if len(p.PublicalyAvailablePostId) != 2 {
			t.Errorf("%s: the deprecated field must stay populated, got %v", name, p.PublicalyAvailablePostId)
		}
		if p.Status != "PUBLISH_COMPLETE" {
			t.Errorf("%s: Status = %q", name, p.Status)
		}
	}
}

func TestCreatorInfoExposesMaxDuration(t *testing.T) {
	const body = `{"data":{"creator_username":"tiktok","max_video_post_duration_sec":300,
		"privacy_level_options":["PUBLIC_TO_EVERYONE","SELF_ONLY"]},"error":{"code":"ok"}}`
	var r QueryCreatorInfoResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatal(err)
	}
	// Without this field the caller cannot do what creator_info exists for:
	// check the video fits before uploading it.
	if r.Data.MaxVideoPostDurationSec != 300 {
		t.Errorf("MaxVideoPostDurationSec = %d, want 300", r.Data.MaxVideoPostDurationSec)
	}
}

func TestUserInfoDecodesTheWholeProfile(t *testing.T) {
	const body = `{"data":{"user":{"open_id":"o","union_id":"u","display_name":"Nome",
		"username":"nome","is_verified":true,"follower_count":42,"likes_count":7}},"error":{"code":"ok"}}`
	var r UserInfoResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatal(err)
	}
	u := r.Data.User
	if u.DisplayName != "Nome" || u.Username != "nome" || !u.IsVerified || u.FollowerCount != 42 || u.LikesCount != 7 {
		t.Errorf("user = %+v: fields beyond open_id/union_id/avatar_url are dropped", u)
	}
}
