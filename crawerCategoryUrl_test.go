package main

import (
	"fmt"
	"strconv"
	"testing"
	//	"time"
)

//测试单个，从html_back中读取一个文件进行测试
//func TestParseCategoryHTML(t *testing.T) {
//	body := ReadAllText("html_back/4.html")
//	cis, total := ParseCategoryHTML(body, 4)
//	fmt.Println(cis)
//	fmt.Println(total)
//}

//测试多个个，从html_back中读取所有文件进行测试
//func TestParseCategoryHTML(t *testing.T) {
//	files, _ := ListDir("html_back", "html")
//	for _, f := range files {
//		idstr := GetFileNameOnly(f)
//		id, err := strconv.Atoi(idstr)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			body := ReadAllText(f)
//			cis, total := ParseCategoryHTML(body, id)
//			fmt.Printf("类别%d有%d,总数有%d\n", id, len(cis), total)

//		}

//	}

//}

//测试单个个，从html_back中读取所有文件进行测试
//func TestCategoryUrlCrawer(t *testing.T) {
//	cid := 708
//	body := ReadAllText(fmt.Sprintf("html/%d.html", cid))
//	crawer := NewCategoryUrlCrawer()
//	cis, total := crawer.parase(body, cid)
//	fmt.Printf("一共有%d\n", len(cis))
//	fmt.Printf("一共有%d\n", total)

//}

//测试多个，从html_back中读取所有文件进行测试
//func TestCategoryUrlCrawer(t *testing.T) {
//	crawer := NewCategoryUrlCrawer()
//	//	cuDB := NewCategroyUrlDB()

//	files, _ := ListDir("html", "html")
//	for _, f := range files {
//		idstr := GetFileNameOnly(f)
//		id, err := strconv.Atoi(idstr)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			body := ReadAllText(f)
//			cis, total := crawer.parase(body, id)
//			if len(cis) == 0 && total == 0 {
//				fmt.Printf("类别%d有%d,总数有%d\n", id, len(cis), total)
//			}

//			//			ci, err := cuDB.GetById(id)
//			//			if err != nil {
//			//				fmt.Println(err)
//			//				continue
//			//			}
//			//			ci.Nums = total
//			//			crawer.handler(ci, cis, total)

//		}

//	}
//}

func TestCategoryUrlRule(t *testing.T) {
	//	crawer := NewCategoryUrlRule()
	//	sTotal := "1-16 sur 5 221 961 résultats pour Boutique Kindle : Ebooks Kindle"
	sTotal := "1-24 sur 16548 rsultats pour 0  150 EUR : Produits Handmade : Bijoux"

	total := getTotal(sTotal)
	fmt.Println(strconv.Itoa(total))
}
