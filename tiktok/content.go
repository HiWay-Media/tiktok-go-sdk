package tiktok

import (
	"encoding/json"
	"fmt"
)

/*
Query Creator Info
To initiate a direct post to a creator's account, you must first use the Query Creator Info endpoint to get the target creator's latest information. For more information about why creator information is necessary, refer to these UX guidelines.

Request:

curl --location --request POST 'https://open.tiktokapis.com/v2/post/publish/creator_info/query/' \
--header 'Authorization: Bearer act.example12345Example12345Example' \
--header 'Content-Type: application/json; charset=UTF-8'
*/
func (o *tiktok) CreatorInfo() (*QueryCreatorInfoResponse, error) {
	resp, err := o.restyPost(API_QUERY_CREATOR_INFO, nil)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("creator info error %s", resp.String())
	}
	var obj QueryCreatorInfoResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}

/*
Post a video
To initiate video upload on TikTok's server, you must invoke the Direct Post Video endpoint. You have the following two options:

If you have the video file locally, set the source parameter to FILE_UPLOAD in your request.
If the video is hosted on a URL, set the source parameter to PULL_FROM_URL.
Example
Example using source=FILE_UPLOAD:

Request:

curl --location 'https://open.tiktokapis.com/v2/post/publish/video/init/' \
--header 'Authorization: Bearer act.example12345Example12345Example' \
--header 'Content-Type: application/json; charset=UTF-8' \

	--data-raw '{
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
	}'
*/
func (o *tiktok) PostVideoInit(title, videoUrl string, privacyLevel string) (*PublishVideoResponse, error) {
	if !CheckPrivacyLevel(privacyLevel) {
		return nil, PrivacyLevelWrong
	}
	//
	request := &PublishVideoRequest{
		PostInfo: PostInfo{
			Title:          title,
			PrivacyLevel:   privacyLevel,
			DisableDuet:    false,
			DisableComment: false,
			DisableStitch:  false,
		},
		SourceInfo: SourceInfo{
			Source:   "PULL_FROM_URL",
			VideoUrl: videoUrl,
		},
	}
	resp, err := o.restyPost(API_POST_PUBLISH_VIDEO_INIT, request)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("post video init error %s", resp.String())
	}
	var obj PublishVideoResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}

/*
Get Post Status
For content uploaded with the Content Posting API, two mechanisms are provided for developers to check the status of the post by the TikTok user:

Fetch Status endpoint: An API endpoint for polling the status of the post.
Content Posting webhooks: Events that notify your registered endpoint of the final outcome of the pos

POST /v2/post/publish/status/fetch/ HTTP /1.1
Host: open.tiktokapis.com
Authorization: Bearer {{AccessToken}}
Content-Type: application/json; charset=UTF-8

	{
	    "publish_id": {PUBLISH_ID}
	}
*/
func (o *tiktok) PublishVideo(publishId string) (*PublishStatusFetchResponse, error) {
	request := &PublishStatusFetchRequest{
		PublishId: publishId,
	}
	resp, err := o.restyPost(API_PUBLISH_STATUS_FETCH, request)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("publis video error %s", resp.String())
	}
	var obj PublishStatusFetchResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}

/* Photo
The /v2/post/publish/content/init/ endpoint allows you to post directly or upload photos to TikTok.

curl --location 'https://open.tiktokapis.com/v2/post/publish/content/init/' \
--header 'Authorization: Bearer act.example12345Example12345Example' \
--header 'Content-Type: application/json' \
--data-raw '{
    "post_info": {
        "title": "funny cat",
        "description": "this will be a #funny photomode on your @tiktok #fyp"
    },
    "source_info": {
        "source": "PULL_FROM_URL",
        "photo_cover_index": 1,
        "photo_images": [
            "https://tiktokcdn.com/obj/example-image-01.webp",
            "https://tiktokcdn.com/obj/example-image-02.webp"
        ]
    },
    "post_mode": "MEDIA_UPLOAD",
    "media_type": "PHOTO"
}'
*/

func (o *tiktok) PostPhotoInit(title, description, privacyLevel string, photoUrls []string, photoMode string) (*PublishStatusFetchResponse, error) {
	if !CheckPrivacyLevel(privacyLevel) {
		return nil, PrivacyLevelWrong
	}
	if !CheckPostMode(photoMode) {
		return nil, PhotoModeWrong
	}
	request := &PublishPhotoRequest{
		PostInfo: PostPhotoInfo{
			Title:          title,
			PrivacyLevel:   privacyLevel,
			DisableComment: false,
			AutoAddMusic:   true,
		},
		SourceInfo: PhotoSourceInfo{
			Source:          "PULL_FROM_URL",
			PhotoCoverIndex: 1,
			PhotoImages:     photoUrls,
		},
		PostMode:  photoMode,
		MediaType: "PHOTO",
	}
	resp, err := o.restyPost(POST_PUBLISH_CONTENT_INIT, request)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("post video init error %s", resp.String())
	}
	var obj PublishStatusFetchResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}

/*
	List Videos

The /v2/video/list/ endpoint can return a paginated list for the given user's public TikTok video posts, sorted by create_time in descending order.

curl -L -X POST 'https://open.tiktokapis.com/v2/video/list/?fields=cover_image_url,id,title' \
-H 'Authorization: Bearer act.example12345Example12345Example' \
-H 'Content-Type: application/json' \

	--data-raw '{
	    "max_count": 20
	}'
*/
func (o *tiktok) GetVideoList(count int64) (*VideoListResponse, error) {
	request := &VideoListRequest{
		MaxCount: count,
	}
	resp, err := o.restyPostWithQueryParams(API_VIDEO_LIST, request, map[string]string{
		"fields": "cover_image_url,id,title",
	})
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("post video init error %s", resp.String())
	}
	var obj VideoListResponse
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}
