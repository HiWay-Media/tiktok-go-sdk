package tiktok

import "golang.org/x/oauth2"

// Endpoint is TikTok's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://www.tiktok.com/v2/auth/authorize/",
	TokenURL: "https://open.tiktokapis.com/v2/oauth/token",
}
