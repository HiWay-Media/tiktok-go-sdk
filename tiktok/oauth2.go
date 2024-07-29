package tiktok

import "golang.org/x/oauth2"

// Endpoint is TikTok's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://www.tiktok.com/v2/auth/authorize/",
	TokenURL: "https://open.tiktokapis.com/v2/oauth/token",
}

/*
# Client Access Token Management

Client access token is a type of access token that does not need user authorization. This is typically used by clients to access resources about themselves or a TikTok application, rather than to access a user's resources. The use cases are to access Research API and Commercial Content API.

curl --location --request POST 'https://open.tiktokapis.com/v2/oauth/token/' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Cache-Control: no-cache' \
--data-urlencode 'client_key=CLIENT_KEY' \
--data-urlencode 'client_secret=CLIENT_SECRET' \
--data-urlencode 'grant_type=client_credentials'
*/
func  (o *tiktok) GetClientAccessTokenManagement() (*AccessTokenManagement, error) {
	data := map[string]string{}
	data["client_key"] = o.clientKey
	data["client_secret"] = o.clientSecret
	data["grant_type"] = "client_credentials"
	//
	resp, err := o.restyPostFormUrlEncoded(Endpoint.TokenURL, data)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("client access token management error %s", resp.String())
	}
	//
	var obj AccessTokenManagement
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	o.debugPrint(obj)
	return &obj, nil
}

}


