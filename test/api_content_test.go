package test

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	resp, err := c.CreatorInfo()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println("resp ", resp)
}
