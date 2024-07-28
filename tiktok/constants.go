package tiktok

import "fmt"

const (
	BASE_URL                    = "https://open.tiktokapis.com"
	QUERY_CREATOR_INFO          = "/v2/post/publish/creator_info/query"
	POST_PUBLISH_VIDEO_INIT     = "/v2/post/publish/video/init"
    PUBLISH_STATUS_FETCH        = "/v2/post/publish/status/fetch/"
    POST_PUBLISH_CONTENT_INIT   = "/v2/post/publish/content/init/"
    USER_INFO  			 		= "/v2/user/info/"
	VIDEO_LIST 					= "/v2/video/list/"
	RESEARCH_VIDEO_QUERY		= "/v2/research/video/query/"
)

var (
	API_QUERY_CREATOR_INFO          = fmt.Sprintf("%s%s", BASE_URL, QUERY_CREATOR_INFO)
	API_POST_PUBLISH_VIDEO_INIT     = fmt.Sprintf("%s%s", BASE_URL, POST_PUBLISH_VIDEO_INIT)
	API_PUBLISH_STATUS_FETCH        = fmt.Sprintf("%s%s", BASE_URL, PUBLISH_STATUS_FETCH)
	API_POST_PUBLISH_CONTENT_INIT   = fmt.Sprintf("%s%s", BASE_URL, POST_PUBLISH_CONTENT_INIT)
	API_USER_INFO   				= fmt.Sprintf("%s%s", BASE_URL, USER_INFO)
	API_VIDEO_LIST   				= fmt.Sprintf("%s%s", BASE_URL, VIDEO_LIST)
	API_RESEARCH_VIDEO_QUERY		= fmt.Sprintf("%s%s", BASE_URL, RESEARCH_VIDEO_QUERY)
)