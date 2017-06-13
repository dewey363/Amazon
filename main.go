package main

import (
	"fmt"
	//	"os"
	"strings"
	//	"math/rand"
	"strconv"
	//	"time"

	"github.com/Sirupsen/logrus"
)

//var memFile *os.File

var dbLog *logrus.Logger = InitLog("db.log")

func main() {
	//	ttParse()
	mainControll()
}
func mainControll() {
	initSys()

	if RUN_FLAG == 1 {
		fmt.Println("开始采集类别链接...")
		crawer := NewCategoryUrlCrawer()
		crawer.Run()
	} else if RUN_FLAG == 2 {
		fmt.Println("开始采集产品链接...")
		crawer := NewUrlListsCrawer()
		crawer.Run()
	} else if RUN_FLAG == 3 {
		fmt.Println("开始单独采集产品链接...")
		crawer := NewProductListCrawer()
		crawer.Run()
	} else {
		fmt.Println("开始采集产品数据...")
		crawer := NewProductUrlCrawer()
		crawer.Run()
	}

	dbLog.Infoln("采集完成")
}

/*
	初始化系统参数
*/
func initSys() {
	var AMAZON_LINKS = map[string]string{
		"en": "https://www.amazon.com",
		"jp": "https://www.amazon.co.jp",
		"uk": "https://www.amazon.co.uk",
		"fr": "https://www.amazon.fr",
	}
	AMAZON_LINK = AMAZON_LINKS[LANG]

	//	memFile, err := os.OpenFile("memory.log", os.O_RDWR|os.O_CREATE, 0644)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.WriteHeapProfile(memFile)

	//	go func() {
	//		log.Println(http.ListenAndServe("localhost:6060", nil))
	//	}()

	//	hook := graylog.NewGraylogHook(log_server, map[string]interface{}{"from": LOCAL})

	//	logrus.AddHook(hook)
}
func ttParse() {
	files, _ := ListDir("./pp/html", "html")
	xxx := 0
	var err error
	for _, v := range files {
		var ci UrlInfo
		idstr := GetFileNameOnly(v)
		ci.Id, err = strconv.Atoi(idstr)
		if err != nil {
			fmt.Println(err)
		}
		//	cid := 708
		body := ReadAllText(v)
		pi := parase(body, ci)
		if strings.Count(pi.Title, "") < 1 || pi.Price == "0" ||
			strings.Count(pi.Desc, "") < 1 || strings.Count(pi.Imgs, "") < 1 ||
			strings.Count(pi.Attr, "") < 1 {
			//			fmt.Println(pi.Id)/
			ss := fmt.Sprintf("%d\n", pi.Id)
			//			os.Rename(v, v+".txt")
			appendToFile("./urlid.txt", ss)
			xxx = xxx + 1
		}

		//doutFile(pi)

	}
	fmt.Printf("一共有%d条记录，失败的有%d\n", len(files), xxx)
}

//func downHTML() {
//	start := time.Now()
//	var context DBContext
//	context.tableName = PRODUCT_URL_TALBE
//	cuDB := NewDB(context)
//	ipPool := OpenRedisPool()
//	for i := 0; i <= 60000; i++ {
//		r := rand.New(rand.NewSource(time.Now().UnixNano()))
//		x := r.Intn(4650000)
//		ci, _ := cuDB.GetById(x)
//		n := time.Now()
//		proxy := ipPool.PopProxy(MAX_GET_TIME)
//		body, err := GetHTMLFromProxy(ci.Url, proxy.Url)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			//			fmt.Println(body)
//			WriteHTMLToFile(ci.Id, body)
//			ipPool.PushProxy(*proxy)
//		}
//		WaitSecond(n, RandSecond(MAX_GET_TIME))
//	}
//	fmt.Println(time.Now().Sub(start).String())
//	fmt.Println("daili测试完成")
//}
