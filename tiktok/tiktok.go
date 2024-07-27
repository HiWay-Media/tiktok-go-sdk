package tiktok

import (
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

type ITiktok interface {
	//
	HealthCheck() error
	IsDebug() bool
	CodeAuthUrl() string
	CreatorInfo() (*QueryCreatorInfoResponse, error)
	PostVideoInit(title, videoUrl string, privacyLevel string) (*PublishVideoResponse, error)
	PublishVideo(publishId string) (*PublishStatusFetchResponse, error)
	//
}

type tiktok struct {
	restClient   *resty.Client
	debug        bool
	clientKey    string
	clientSecret string
	accessToken  string
	OAuth2Config *oauth2.Config
}

func NewTikTok(clientKey, clientSecret string, isDebug bool) (ITiktok, error) {
	o := &tiktok{
		clientKey:    clientKey,
		clientSecret: clientSecret,
		restClient:   resty.New(),
		debug:        isDebug,
		OAuth2Config: &oauth2.Config{
			ClientID:     clientKey,
			ClientSecret: clientSecret,
			RedirectURL:  "",
			Scopes:       []string{"email"},
			Endpoint:     Endpoint,
		},
	}
	o.restClient.SetDebug(isDebug)
	o.restClient.SetBaseURL(BASE_URL)
	return o, nil
}

func (o *tiktok) SetAccessToken(token string) {
	o.accessToken = token
}

func (o *tiktok) GetAccessToken() string {
	return o.accessToken
}

func (o *tiktok) CodeAuthUrl() string {
	return o.OAuth2Config.AuthCodeURL("state-token", oauth2.ApprovalForce)
}