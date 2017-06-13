package main

import (
	"fmt"
	//	"net/http"
	"testing"
)

func TestListDir(t *testing.T) {
	ips, err := ListDir("html_back", "html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ips)
}
