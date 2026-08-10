package tiktok

import "encoding/json"

/*
{
   "data":{
      "creator_avatar_url": "https://lf16-tt4d.tiktokcdn.com/obj/tiktok-open-platform/8d5740ac3844be417beeacd0df75aef1",
      "creator_username": "tiktok",
      "creator_nickname": "TikTok Official",
      "privacy_level_options": ["PUBLIC_TO_EVERYONE", "MUTUAL_FOLLOW_FRIENDS", "SELF_ONLY"]
      "comment_disabled": false,
      "duet_disabled": false,
      "stitch_disabled": true,
      "max_video_post_duration_sec": 300
   },
    "error": {
         "code": "ok",
         "message": "",
         "log_id": "202210112248442CB9319E1FB30C1073F3"
     }
}
*/

type QueryCreatorInfoResponse struct {
	Data  DataQueryCreatorInfo `json:"data"`
	Error ErrorObject          `json:"error"`
}

type DataQueryCreatorInfo struct {
	CreatorAvatarUrl    string   `json:"creator_avatar_url"`
	CreatorUsername     string   `json:"creator_username"`
	CreatorNickname     string   `json:"creator_nickname"`
	PrivacyLevelOptions []string `json:"privacy_level_options"`
	CommentDisabled     bool     `json:"comment_disabled"`
	DuetDisabled        bool     `json:"duet_disabled"`
	StitchDisabled      bool     `json:"stitch_disabled"`
	// MaxVideoPostDurationSec is the longest video this creator may post. It is
	// the reason TikTok requires a creator_info query before publishing: posting
	// something longer is rejected, and the caller can only know from here.
	MaxVideoPostDurationSec int64 `json:"max_video_post_duration_sec"`
}

type PublishVideoResponse struct {
	Data  DataPublishVideo `json:"data"`
	Error ErrorObject      `json:"error"`
}

type DataPublishVideo struct {
	PublishId string `json:"publish_id"`
	// Deprecated: use PublishId. Kept, and kept populated, because it shipped
	// as public API with the typo; two fields cannot share a JSON tag, so it is
	// filled by UnmarshalJSON instead.
	PubblishId string `json:"-"`
}

// UnmarshalJSON decodes the publish id and mirrors it into the deprecated
// misspelled field.
func (d *DataPublishVideo) UnmarshalJSON(b []byte) error {
	type alias DataPublishVideo
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	a.PubblishId = a.PublishId
	*d = DataPublishVideo(a)
	return nil
}

type PublishStatusFetchResponse struct {
	Data  PublishStatusFetch `json:"data"`
	Error ErrorObject        `json:"error"`
}

type PublishStatusFetch struct {
	Status        string `json:"status"`
	FailReason    string `json:"fail_reason"`
	UploadedBytes int64  `json:"uploaded_bytes"`

	// PubliclyAvailablePostId holds the ids of the posts that became public.
	//
	// Both spellings are accepted off the wire ("publicly_" and TikTok's own
	// "publicaly_"): which one arrives is not something this SDK should bet on,
	// and betting wrong means the field is silently always empty.
	PubliclyAvailablePostId []int64 `json:"-"`
	// Deprecated: use PubliclyAvailablePostId.
	PublicalyAvailablePostId []int64 `json:"-"`
}

// UnmarshalJSON accepts either spelling of the post id field and fills both the
// current and the deprecated one.
func (p *PublishStatusFetch) UnmarshalJSON(b []byte) error {
	var raw struct {
		Status        string  `json:"status"`
		FailReason    string  `json:"fail_reason"`
		UploadedBytes int64   `json:"uploaded_bytes"`
		Correct       []int64 `json:"publicly_available_post_id"`
		Misspelled    []int64 `json:"publicaly_available_post_id"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	ids := raw.Correct
	if ids == nil {
		ids = raw.Misspelled
	}
	*p = PublishStatusFetch{
		Status:                   raw.Status,
		FailReason:               raw.FailReason,
		UploadedBytes:            raw.UploadedBytes,
		PubliclyAvailablePostId:  ids,
		PublicalyAvailablePostId: ids,
	}
	return nil
}

type VideoListResponse struct {
	Data  DataVideoList `json:"data"`
	Error ErrorObject   `json:"error"`
}

type DataVideoList struct {
	Videos []Video `json:"videos"`
}

type Video struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	CoverImageUrl string `json:"cover_image_url"`
}
