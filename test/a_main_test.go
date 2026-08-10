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

// newClient builds a client with placeholder credentials, for the tests that
// never reach the network. NewTikTok rejects empty credentials, so these must
// not be blank.
func newClient(t *testing.T) tiktok.ITiktok {
	t.Helper()
	c, err := tiktok.NewTikTok("test-client-key", "test-client-secret", false)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// GetTikTok builds a client with the real app credentials, and skips the test
// when they are not in the environment. Only the integration tests need it: a
// missing credential must skip, never fail — otherwise a fork without the
// repository secrets sees a red CI it cannot fix.
func GetTikTok(t *testing.T) tiktok.ITiktok {
	t.Helper()
	clientKey := os.Getenv("CLIENT_KEY")
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientKey == "" || clientSecret == "" {
		t.Skip("CLIENT_KEY/CLIENT_SECRET not set")
	}
	c, err := tiktok.NewTikTok(clientKey, clientSecret, false)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNewTikTok(t *testing.T) {
	c := newClient(t)
	log.Println(c.IsDebug())
}

func TestAuthCodeUrl(t *testing.T) {
	c := newClient(t)
	resp := c.CodeAuthUrl()
	if resp == "" {
		t.Fatal("CodeAuthUrl returned an empty URL")
	}
	log.Println("resp ", resp)
}
