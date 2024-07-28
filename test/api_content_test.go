package test

import (
	"log"
	"testing"
)

func TestCreatorInfo(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.IsDebug()
	/*resp, err := c.CreatorInfo()
	if err != nil {
		t.Fatalf(err.Error())
	}*/
	log.Println("resp ", c.IsDebug())
}


func TestGetVideoList( t *testing.T){
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.IsDebug()
	resp, err := c.GetVideoList(int64(20))
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println("resp ", resp)
}