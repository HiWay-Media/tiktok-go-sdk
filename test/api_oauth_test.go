package test

import (
	"log"
	"testing"
)

func TestGetClientAccessTokenManagement(t *testing.T) {
	c, err := GetTikTok()
	if err != nil {
		t.Fatalf(err.Error())
	}
	c.IsDebug()
	resp, err := c.GetClientAccessTokenManagement()
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Println("resp ", resp)
}
