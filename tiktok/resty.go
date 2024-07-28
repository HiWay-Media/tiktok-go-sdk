package tiktok

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func (o *tiktok) HealthCheck() error {
	resp, err := o.restyPost("/", nil)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("healthcheck error")
	}
	return nil
}

func (o *tiktok) IsDebug() bool {
	return o.debug
}

// Resty Methods

func (o *tiktok) restyPost(url string, body interface{}) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("Auhtorization", "Bearer "+o.accessToken).
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *tiktok) restyPostWithQueryParams(url string, body interface{},  queryParams map[string]string) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("Auhtorization", "Bearer "+o.accessToken).
		SetQueryParams(queryParams).
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *tiktok) restyGet(url string, queryParams map[string]string) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Auhtorization", "Bearer "+o.accessToken).
		SetQueryParams(queryParams).
		Get(url)
	//
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *tiktok) debugPrint(data ...interface{}) {
	if o.debug {
		log.Println(data...)
	}
}
