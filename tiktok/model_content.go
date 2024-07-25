package tiktok

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
	Data  DataQueryCreatorInfo  `json:"data"`
	Error ErrorQueryCreatorInfo `json:"error"`
}

type DataQueryCreatorInfo struct {
	CreatorAvatarUrl    string   `json:"creator_avatar_url"`
	CreatorUsername     string   `json:"creator_username"`
	CreatorNickname     string   `json:"creator_nickname"`
	PrivacyLevelOptions []string `json:"privacy_level_options"`
	CommentDisabled     bool     `json:"comment_disabled"`
	DuetDisabled        bool     `json:"duet_disabled"`
	StitchDisabled      bool     `json:"stitch_disabled"`
}

type ErrorQueryCreatorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LogId   string `json:"log_id"`
}
