package test

import (
	"os"
	"testing"
    "github.com/HiWay-Media/tiktok-go-sdk/tiktok"
)

func TestMain(m *testing.M) {
    c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.IsDebug()
}


func GetTikTok() (tiktok.ITiktok, error) {
	clientKey := os.GetEnv("clientKey")
	clientSecret := os.GetEnv("clientSecret")
	//
	c, err := tiktok.NewTikTok(clientKey, clientSecret, false)
	if err != nil {
		return nil, err
	}
	return c, nil
}