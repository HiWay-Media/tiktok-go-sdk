//go:build integration

// Tests that talk to the real TikTok API. They are excluded from `go test ./...`
// on purpose (TT-22): a default test run must not depend on live credentials,
// on the network, or on the app's API quota.
//
//	go test -tags integration ./test/
package test

import (
	"log"
	"testing"
)

func TestGetClientAccessTokenManagement(t *testing.T) {
	c := GetTikTok(t)
	resp, err := c.GetClientAccessTokenManagement()
	if err != nil {
		t.Fatal(err)
	}
	if resp.AccessToken == "" {
		t.Fatal("empty access token")
	}
	log.Println("token type ", resp.TokenType, " expires in ", resp.ExpiresIn)
}
