package test

import (
	"log"
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

func TestNewTikTok(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println(c.IsDebug())
}

func TestAuthCodeUrl(t *testing.T){
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	resp := c.CodeAuthUrl()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println("resp ", resp)
}



func GetTikTok() (tiktok.ITiktok, error) {
	clientKey := os.Getenv("CLIENT_KEY")
	clientSecret := os.Getenv("CLIENT_SECRET")
	//
	c, err := tiktok.NewTikTok(clientKey, clientSecret, false)
	if err != nil {
		return nil, err
	}
	return c, nil
}

