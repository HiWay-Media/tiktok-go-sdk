package test

import (
	"os"
	"testing"

	"github.com/HiWay-Media/tiktok-go-sdk/tiktok"
)

func Test(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	resp , err := c.CreatorInfo()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println("resp ", resp)
}