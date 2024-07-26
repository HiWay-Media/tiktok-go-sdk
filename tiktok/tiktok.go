package tiktok

import "github.com/go-resty/resty/v2"

type ITiktok interface {
	//
	HealthCheck() error
	IsDebug() bool
	CreatorInfo() (*QueryCreatorInfoResponse, error)
	PostVideoInit(title, videoUrl string, privacyLevel string) (*PublishVideoResponse, error)
	//
}

type tiktok struct {
	restClient   *resty.Client
	debug        bool
	clientKey    string
	clientSecret string
	accessToken string
}

func NewTikTok(clientKey, clientSecret string, isDebug bool) (ITiktok, error) {
	o := &tiktok{
		clientKey:    clientKey,
		clientSecret: clientSecret,
		restClient:   resty.New(),
		debug:        isDebug,
	}
	o.restClient.SetDebug(isDebug)
	o.restClient.SetBaseURL(BASE_URL)
	return o, nil
}


func (o *tiktok) SetAccessToken(token string){
	o.accessToken = token
}

func (o *tiktok) GetAccessToken() string {
	return o.accessToken
}