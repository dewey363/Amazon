package main

import (
	//	"fmt"

	//	"strconv"
	"testing"
)

//测试单个，中读取一个文件进行测试
func TestProductUrlRule(t *testing.T) {
	//	filePath, err := getFolder(444444443)
	//	if err == nil {
	//		t.Log(filePath)
	//	}
	//	appendToFile(filePath, "ddd")
}

//测试分析单个页面文件
//func TestProductUrlRule(t *testing.T) {
//	filePath, err := getFolder(444444443)
//	if err == nil {
//		t.Log(filePath)
//	}
//	appendToFile(filePath, "ddd")
//}
//func TestAsin(t *testing.T) {
//	var ci UrlInfo
//	ci.Id = 1735337
//	//	cid := 708
//	body := ReadAllText(fmt.Sprintf("pp/html/%d.html", ci.Id))
//	parase(body, ci)
//	//	fmt.Println(parase(body, ci))
//	//	fmt.Println(pis)
//	//	fmt.Printf("一共有%d\n", len(pis))

//}

//func TestParse(t *testing.T) {
//	files, _ := ListDir("./pp/html", "html")
//	for _, v := range files {
//		var ci UrlInfo
//		idstr := GetFileNameOnly(v)
//		ci.Id, _ = strconv.Atoi(idstr)
//		//	cid := 708
//		body := ReadAllText(v)
//		parase(body, ci)

//	}

//}

//测试getprice函数
func TestGetPrice(t *testing.T) {
	t.Logf("%s", getPrice("Prix :EUR dddTous les prix incluent la TVA."))
}

//func TestMultiGetPrice(t *testing.T) {
//	files, _ := ListDir("./test/price", "txt")
//	for _, v := range files {
//		cc := ReadAllText(v)
//		sprice := getPrice(cc)
//		if sprice == "0" {
//			t.Log(v)
//		}
//	}
//}

//func TestGetFolder(t *testing.T) {
//	getFolder(111)
//}

//测试获取页面数据的parase
