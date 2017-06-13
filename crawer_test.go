package main

import (
	"fmt"
	"math/rand"
	"time"
	//	"net/http"
	//	"errors"
	//	"strings"
	//	"sync"
	"testing"
)

//var shellqueue *goqueue.Queue

////单个网址测试
//func TestGetHTMLFromProxy(t *testing.T) {
//	b, err := GetHTMLFromProxy(2, "", "")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(b)

//}

////多个网址测试
//func TestMultGetHTMLFromProxy(t *testing.T) {
//	start := time.Now()
//	cuDB := NewCategroyUrlDB()
//	cis, _ := cuDB.Prefetch(0, 100)
//	for _, ci := range cis {
//		n := time.Now()

//		body, err := GetHTMLFromProxy(ci.Url, "137.175.66.155")
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			//			fmt.Println(body)
//			WriteHTMLToFile(ci.Id, body)
//		}
//		WaitSecond(n, RandSecond(MAX_GET_TIME))
//	}
//	fmt.Println(time.Now().Sub(start).String())
//	fmt.Println("daili测试完成")

//}

////多个网址多线程测试
//func TestMultGetHTMLFromProxy(t *testing.T) {
//	ipqueue = InitIPQueue()
//	wg := new(sync.WaitGroup)
//	cis, _ := GetCategroyUrlInfo(0, 100)
//	for _, ci := range cis {
//		wg.Add(1)
//		go func(cid int, curl string) {
//			defer wg.Done()
//			daili := GetDaili()
//			//			daili = ""
//			body, err := GetHTMLFromProxy(cid, curl, daili)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//			t.Logf("body rows:%d\n", len(body))
//		}(ci.Id, ci.Url)
//	}
//	wg.Wait()
//	fmt.Println("daili测试完成")

//}
func TestMultGetHTMLFromProxy(t *testing.T) {
	start := time.Now()
	var context DBContext
	context.tableName = PRODUCT_URL_TALBE
	cuDB := NewDB(context)

	for i := 0; i <= 100; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		x := r.Intn(2000000)
		ci := cuDB.GetById(x)
		n := time.Now()

		body, err := GetHTMLFromProxy(ci.Url, "989945957cdb2b214abbd21214259565:squid@108.187.149.3:3128")
		if err != nil {
			fmt.Println(err)
		} else {
			//			fmt.Println(body)
			WriteHTMLToFile(ci.Id, body)
		}
		WaitSecond(n, RandSecond(MAX_GET_TIME))
	}
	fmt.Println(time.Now().Sub(start).String())
	fmt.Println("daili测试完成")

}
