package tiktok

type UserInfoResponse struct {
	Data  DataUserInfo `json:"data"`
	Error ErrorObject  `json:"error"`
}

type DataUserInfo struct {
	User UserInfo `json:"user"`
}

// UserInfo is the user profile returned by /v2/user/info/.
//
// Every field is optional on the wire: TikTok only returns the ones listed in
// the request's `fields` parameter, so a field left empty here usually means it
// was not asked for.
type UserInfo struct {
	OpenID    string `json:"open_id"`
	UnionID   string `json:"union_id"`
	AvatarUrl string `json:"avatar_url"`

	AvatarUrl100    string `json:"avatar_url_100"`
	AvatarLargeUrl  string `json:"avatar_large_url"`
	DisplayName     string `json:"display_name"`
	BioDescription  string `json:"bio_description"`
	ProfileDeepLink string `json:"profile_deep_link"`
	IsVerified      bool   `json:"is_verified"`
	Username        string `json:"username"`
	FollowerCount   int64  `json:"follower_count"`
	FollowingCount  int64  `json:"following_count"`
	LikesCount      int64  `json:"likes_count"`
	VideoCount      int64  `json:"video_count"`
}
