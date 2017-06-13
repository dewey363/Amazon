package main

import (
	"fmt"
	//	"strconv"
	"testing"
	//	"time"
)

func TestUrlListsRule(t *testing.T) {
	var ci, tempci UrlInfo
	tempci.Id = 1
	tempci.Flag = 1
	tempci.Nums = 1
	tempci.Parent = 1
	tempci.Url = "&low-price=0&high-price=150"

	ci = tempci
	fmt.Println(ci)
	ci.Url = "333"
	fmt.Println(tempci)
	fmt.Println(ci)
	//	ci.Id = ci.Id
	//	ci.Flag = ci.Flag
	//	ci.Nums = ci.Nums
	//	ci.Parent = ci.Parent
	//	ci.Url = ci.Url + "&low-price=0&high-price=150"

}

//func add()[]UrlInfo {
//	var pis []UrlInfo

//}
