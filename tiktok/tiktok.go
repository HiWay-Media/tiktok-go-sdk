package tiktok

import "github.com/go-resty/resty/v2"

type ITiktok interface {
	//
	HealthCheck() error
	IsDebug() bool
	//
}

type tiktok struct {
	restClient   *resty.Client
	debug        bool
	clientKey    string
	clientSecret string
}

func NewTikTok(clientKey, clientSecret string, isDebug bool) (ITiktok, error) {
	o := &tiktok{
		clientKey:    clientKey,
		clientSecret: clientSecret,
		restClient:   resty.New(),
		debug:        isDebug,
	}
	o.restClient.SetDebug(isDebug)
	return o, nil
}
