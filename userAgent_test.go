package main

import (
	"time"
	//	"net/http"
	"testing"
)

func TestGetRandUserAgent(t *testing.T) {
	for i := 1; i < 5; i++ {
		agent := GetRandUserAgent()
		time.Sleep(time.Second)
		t.Log(agent)
	}
}
