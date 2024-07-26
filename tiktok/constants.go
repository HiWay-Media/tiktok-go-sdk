package tiktok

import "fmt"

const (
	BASE_URL                = "https://open.tiktokapis.com/"
	QUERY_CREATOR_INFO      = "v2/post/publish/creator_info/query"
	POST_PUBLISH_VIDEO_INIT = "v2/post/publish/video/init"
)

var (
	API_QUERY_CREATOR_INFO      = fmt.Sprintf("%s%s", BASE_URL, QUERY_CREATOR_INFO)
	API_POST_PUBLISH_VIDEO_INIT = fmt.Sprintf("%s%s", BASE_URL, POST_PUBLISH_VIDEO_INIT)
)
