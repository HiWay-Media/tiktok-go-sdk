package tiktok

import (
	"encoding/json"
)

/*
Get User Info

The /v2/user/info/ endpoint returns some basic information for a given TikTok user.

# HTTP URL

Request:
curl -L -X GET 'https://open.tiktokapis.com/v2/user/info/?fields=open_id,union_id,avatar_url' \
-H 'Authorization: Bearer act.example12345Example12345Example'
*/
func (o *tiktok) UserInfo() (*UserInfoResponse, error) {
	resp, err := o.restyGet(USER_INFO, nil)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(resp, "user info"); err != nil {
		return nil, err
	}
	var obj UserInfoResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}
