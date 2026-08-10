package test

import (
	"log"
	"testing"
)

func TestCreatorInfo(t *testing.T) {
	c := newClient(t)
	log.Println("resp ", c.IsDebug())
}

func TestGetVideoList(t *testing.T) {
	c := newClient(t)
	log.Println("resp ", c.IsDebug())
}
