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
type PublishVideoRequest struct {
    PostInfo PostInfo       `json:"post_info"`
    SourceInfo SourceInfo   `json:"source_info"`
}

type PostInfo struct { 
    Title                   string `json:"title"`
    PrivacyLevel            string `json:"privacy_level"`
    DisableDuet             bool `json:"disable_duet"`
    DisableComment          bool `json:"disable_comment"`
    DisableStitch           bool `json:"disable_stitch"`
    VideoCoverTimestampMS   int `json:"video_cover_timestamp_ms"`
}

type SourceInfo struct {
    Source      string `json:"source"`
    VideoUrl    string `json:"video_url"`
}