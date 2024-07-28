package tiktok

type UserInfoResponse struct {
	Data  DataUserInfo `json:"data"`
	Error ErrorObject  `json:"error"`
}

type DataUserInfo struct {
	User UserInfo `json:"user"`
}

type UserInfo struct {
	OpenID    string `json:"open_id"`
	UnionID   string `json:"union_id"`
	AvatarUrl string `json:"avatar_url"`
}
