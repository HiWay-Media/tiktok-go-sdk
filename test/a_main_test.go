package test

import (
	"os"
	"testing"

	"github.com/HiWay-Media/tiktok-go-sdk/tiktok"
)

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "" {
		err := os.Setenv("APP_ENV", "test")
		if err != nil {
			panic("could not set test env")
		}
	}
	//env.Load()
	m.Run()
}

func TestNewCompress(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.IsDebug()
}

func GetTikTok() (tiktok.ITiktok, error) {
	clientKey := os.Getenv("clientKey")
	clientSecret := os.Getenv("clientSecret")
	//
	c, err := tiktok.NewTikTok(clientKey, clientSecret, false)
	if err != nil {
		return nil, err
	}
	return c, nil
}
