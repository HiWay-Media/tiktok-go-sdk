package tiktok

/*
{
  "post_info": {
    "title": "this will be a funny #cat video on your @tiktok #fyp",
    "privacy_level": "MUTUAL_FOLLOW_FRIENDS",
    "disable_duet": false,
    "disable_comment": true,
    "disable_stitch": false,
    "video_cover_timestamp_ms": 1000
  },
  "source_info": {
      "source": "PULL_FROM_URL",
      "video_url": "https://example.verified.domain.com/example_video.mp4",
  }
}
*/
/*type BasePostInfo struct {
}*/

type PublishVideoRequest struct {
	PostInfo   PostInfo   `json:"post_info"`
	SourceInfo SourceInfo `json:"source_info"`
}

type PostInfo struct {
	Title string `json:"title"`
	// Not a documented post_info field for a video (the documented text field
	// is Title), so it is only sent when the caller sets it.
	Description           string `json:"description,omitempty"`
	PrivacyLevel          string `json:"privacy_level"`
	DisableDuet           bool   `json:"disable_duet"`
	DisableComment        bool   `json:"disable_comment"`
	DisableStitch         bool   `json:"disable_stitch"`
	VideoCoverTimestampMS int64  `json:"video_cover_timestamp_ms"`
}

type SourceInfo struct {
	Source   string `json:"source"`
	VideoUrl string `json:"video_url"`
}

type PublishStatusFetchRequest struct {
	PublishId string `json:"publish_id"`
}

type PublishPhotoRequest struct {
	PostInfo   PostPhotoInfo   `json:"post_info"`
	SourceInfo PhotoSourceInfo `json:"source_info"`
	PostMode   string          `json:"post_mode"`  //  DIRECT_POST | MEDIA_UPLAOD
	MediaType  string          `json:"media_type"` // PHOTO
}

type PostPhotoInfo struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	PrivacyLevel   string `json:"privacy_level"`
	DisableComment bool   `json:"disable_comment"`
	AutoAddMusic   bool   `json:"auto_add_music"`
}

type PhotoSourceInfo struct {
	Source          string   `json:"source"`
	PhotoCoverIndex int      `json:"photo_cover_index"`
	PhotoImages     []string `json:"photo_images"`
}

type VideoListRequest struct {
	MaxCount int64 `json:"max_count"`
}

// VideoPost describes a video to publish with PostVideo. It carries every
// documented post_info option, including VideoCoverTimestampMS, which the
// positional PostVideoInit has no way to express.
type VideoPost struct {
	Title string
	// Description is sent only when non-empty: it is not part of the documented
	// post_info for a video (the documented text field is Title), so an empty
	// one must not end up on the wire.
	Description           string
	VideoUrl              string
	PrivacyLevel          string
	DisableDuet           bool
	DisableComment        bool
	DisableStitch         bool
	VideoCoverTimestampMS int64
}

// PhotoPost describes a photo post for PostPhoto. The zero value posts with
// comments enabled, no auto-added music and cover index 0; PostPhotoInit keeps
// its historical defaults (music on, cover index 1) instead.
type PhotoPost struct {
	Title           string
	Description     string
	PrivacyLevel    string
	PhotoUrls       []string
	PostMode        string
	DisableComment  bool
	AutoAddMusic    bool
	PhotoCoverIndex int
}
