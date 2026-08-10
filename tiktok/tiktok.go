package tiktok

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

type ITiktok interface {
	//
	HealthCheck() error
	IsDebug() bool
	CodeAuthUrl() string
	SetAccessToken(token string)
	GetAccessToken() string
	//
	GetClientAccessTokenManagement() (*AccessTokenManagement, error)
	CreatorInfo() (*QueryCreatorInfoResponse, error)
	PostVideoInit(title, description, videoUrl string, privacyLevel string, disableDuet, disableComment, disableStitch bool) (*PublishVideoResponse, error)
	PublishVideo(publishId string) (*PublishStatusFetchResponse, error)
	GetVideoList(count int64) (*VideoListResponse, error)
	PostPhotoInit(title, description, privacyLevel string, photoUrls []string, photoMode string) (*PublishStatusFetchResponse, error)
	UserInfo() (*UserInfoResponse, error)
	//
}

type tiktok struct {
	restClient   *resty.Client
	debug        bool
	clientKey    string
	clientSecret string
	OAuth2Config *oauth2.Config

	// baseURL, timeout and httpClient are set by the options and consumed once,
	// while the resty client is being built.
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client

	// accessToken is read on every request and can be replaced at any time by
	// SetAccessToken, so it is guarded: a client is meant to be shared across
	// goroutines.
	mu          sync.RWMutex
	accessToken string
}

// NewTikTok builds a client for the TikTok API.
//
// clientKey and clientSecret are required: they identify the app and every
// endpoint refuses the request without them. Options are optional and keep the
// three-argument form source compatible.
func NewTikTok(clientKey, clientSecret string, isDebug bool, opts ...Option) (ITiktok, error) {
	if clientKey == "" {
		return nil, ErrClientKeyRequired
	}
	if clientSecret == "" {
		return nil, ErrClientSecretRequired
	}
	o := &tiktok{
		clientKey:    clientKey,
		clientSecret: clientSecret,
		debug:        isDebug,
		baseURL:      BASE_URL,
		timeout:      DefaultTimeout,
		OAuth2Config: &oauth2.Config{
			ClientID:     clientKey,
			ClientSecret: clientSecret,
			RedirectURL:  "",
			//Scopes:       []string{"user.info.basic", "video.list", },
			Endpoint: Endpoint,
			// /[]string{"user.info.basic", "video.list", "video.publish", "video.delete", }
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.httpClient != nil {
		o.restClient = resty.NewWithClient(o.httpClient)
	} else {
		o.restClient = resty.New()
	}
	o.restClient.SetDebug(o.debug)
	o.restClient.SetBaseURL(o.baseURL)
	if o.timeout > 0 {
		o.restClient.SetTimeout(o.timeout)
	}
	return o, nil
}

func (o *tiktok) SetAccessToken(token string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.accessToken = token
}

func (o *tiktok) GetAccessToken() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.accessToken
}

func (o *tiktok) CodeAuthUrl() string {
	return o.OAuth2Config.AuthCodeURL("state-token", oauth2.ApprovalForce, oauth2.SetAuthURLParam(
		"client_key",
		o.clientKey,
	))
}
