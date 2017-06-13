package main

import (
	"errors"
	//	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

type UrlListsCrawer struct {
	c *UrlCrawer
}

func NewUrlListsCrawer() *UrlListsCrawer {
	params := &CrawerContext{}
	params.loger = InitLog("urlLists.log")
	var db DBContext

	db.tableName = PRODUCT_LISTS_TABLE
	params.db = NewDB(db)
	params.wg = new(sync.WaitGroup)
	params.prefecthPool = OpenRedisPool()
	params.ipPool = OpenRedisPool()
	params.queueName = PRODUCT_LISTS_URL_KEY

	crawer := &UrlListsCrawer{}
	rule := NewUrlListsRule()
	crawer.c = NewUrlCrawer(params, rule.parseHook)

	return crawer
}
func (self *UrlListsCrawer) Run() {

	self.c.Run()

}

type UrlListsRule struct {
	plDB  *DB
	puDB  *DB
	sbfpu *UrlFilter
	sbfpl *UrlFilter
}

func NewUrlListsRule() *UrlListsRule {
	rule := &UrlListsRule{}
	var context DBContext

	context.tableName = PRODUCT_LIST_TABLE
	rule.plDB = NewDB(context)

	rule.sbfpl = newUrlFilter(rule.plDB)

	var context2 DBContext

	context2.tableName = PRODUCT_URL_TALBE
	rule.puDB = NewDB(context2)
	rule.sbfpu = newUrlFilter(rule.puDB)

	return rule
}
func (self *UrlListsRule) craw(context *CrawerContext, daili string, ci UrlInfo) (string, error) {
	html, err := GetHTMLFromProxy(ci.Url, daili)

	if err != nil {
		return "", err
	}
	if IsRobot(html) {
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
			"ip": daili,
		}).Warnln("验证码页面")
		return "", ERROR_PROXY
	}
	if Is404(html) {
		//如果是404则设置完成
		context.db.SetCompleted(ci.Id)
		return "", PAGE_NOT_FOUND
	}

	return html, nil
}

//因为一个链接需要访问好多次
//&low-price=0&high-price=150 过滤价格
//&page=2 过滤页面
//1.获取网页内容，如果err直接退出
//2.如果len(pis)<24 && opurl!=nil（判断显示格式&lo=digital-text  &lo=videogames）
//修改Url tempci.url=opurl  然后开始从第二页循环
//3.否则把获取的pis插入数据库，从第二页开始保存到product_list，等下次一起处理
func (self *UrlListsRule) parseHook(context *CrawerContext, daili string, ci UrlInfo) error {

	var tempci UrlInfo
	tempci = ci
	tempci.Url = ci.Url + "&low-price=0&high-price=150"
	start := time.Now()
	//1.获取第一页
	html, err := self.craw(context, daili, tempci)
	//第一次如果出错直接返回
	if err != nil {
		return err

	}
	//分析html
	pis, total, viewUrl := self.parase(html, ci)
	if len(pis) <= 0 {
		//		fmt.Println("没有获取到任何产品")
		//没有获取到任何标签，说明采集的标签有问题
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		WriteHTMLToFile(ci.Id*1000+r.Intn(999), html)
		return nil
	}
	//2.如果len(pis)<24 && opurl!=nil 修改Url tempci.url=opurl
	if len(pis) <= 24 && viewUrl != "" {
		tempci.Url = AMAZON_LINK + viewUrl + "&low-price=0&high-price=150"
		//等待
		//		fmt.Println("分析链接")
		WaitSecond(start, RandSecond(MAX_GET_TIME))
		start = time.Now()
		html, err = self.craw(context, daili, tempci)
		WaitSecond(start, RandSecond(MAX_GET_TIME))
		//第一次如果出错直接返回
		if err != nil {
			//			fmt.Println(err)
			return err
		}
		//		WriteHTMLToFile(tempci.Id, html)
		pis, total, _ = self.parase(html, ci)

	}

	if len(pis) <= 0 {
		//		fmt.Println("没有获取到任何产品")
		//没有获取到任何标签，说明采集的标签有问题，输出待人工处理
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		WriteHTMLToFile(ci.Id*1000+r.Intn(999), html)
		return errors.New("没有抓取到信息")
	}
	//保存第一个页面的产品链接，这样下次就少下载一个页面
	self.handlerProductUrl(context, ci, pis)
	//是不是应该保存到product_list中等下次再一次性采集
	//3.把第一页得到的产品链接保存起来，准备第二页
	maxPage := total/len(pis) + 1
	//	fmt.Printf("一共有%d页,每页%d个\n", maxPage, len(pis))
	var plurls []UrlInfo
	for curPage := 2; curPage <= maxPage; curPage++ {
		tempurl := tempci.Url + "&page=" + strconv.Itoa(curPage)
		//保存到product_list等下次再采集
		if !self.sbfpl.TestAndAdd(tempurl) {
			pl := tempci
			pl.Parent = tempci.Parent
			pl.Parentname = tempci.Parentname
			pl.Url = tempurl
			pl.Nums = len(pis)
			plurls = append(plurls, pl)
		}

	}
	context.loger.WithFields(logrus.Fields{
		"id": ci.Id,
	}).Infof("BatchAdd(product list):%d\n", len(plurls))
	self.plDB.BatchAdd(plurls)
	//都添加完后再设置完成标记
	context.db.SetCompleted(ci.Id)
	return nil
}

//处理解析后得到的链接，有可能直接保存数据库，或者使用缓存
func (self *UrlListsRule) handlerProductUrl(context *CrawerContext, pi UrlInfo, PIS []UrlInfo) error {

	context.loger.WithFields(logrus.Fields{
		"id": pi.Id,
	}).Infof("BatchAdd(product url):%d\n", len(PIS))
	//保存到数据库
	var tempcis []UrlInfo
	for _, v := range PIS {
		//如果不存在
		if !self.sbfpu.TestAndAdd(v.Url) {
			tempcis = append(tempcis, v)
		}
	}
	self.puDB.BatchAdd(tempcis)

	return nil
}

/*
根据字符串和父类别编号获取产品URL信息
可能返回url数组为空，获取过滤后的产品数，显示视图地址

*/
func (crawer *UrlListsRule) parase(body string, ci UrlInfo) ([]UrlInfo, int, string) {
	var PIS []UrlInfo
	total := 0
	var utfBody io.Reader = strings.NewReader(body)
	doc, _ := goquery.NewDocumentFromReader(utfBody)
	//获取总数
	totalTag := doc.Find("#s-result-count").First()
	//	fmt.Printf("正在处理总数标签:%d\n", len(totalTag.Nodes))
	if len(totalTag.Nodes) > 0 {

		sTotal := totalTag.Text()

		total = getTotal(sTotal)

	}
	//获取视图链接
	viewUrl := ""
	doc.Find("#s-result-info-bar").
		Find(".s-layout-toggle-picker a").
		Each(func(i1 int, content *goquery.Selection) {
			url, _ := content.Attr("href")
			url = strings.TrimSpace(url)

			if strings.Count(url, "") > 10 && !strings.Contains(url, "lo=none") {
				viewUrl = url

			}
		})
	//获取产品链接 resultsCol
	//使用二级标签
	regexs := []string{"div#mainResults", "div#resultsCol"} //二维数组的赋值初始化
	for _, v := range regexs {
		resultTag := doc.Find(v).First()
		if len(resultTag.Nodes) > 0 {
			resultTag.Find("a.s-access-detail-page").Each(func(i1 int, content *goquery.Selection) {
				var tempci UrlInfo
				url, _ := content.Attr("href")
				url = strings.TrimSpace(url)
				//如果连接少于10个字符说明有问题
				if strings.Count(url, "") > 10 {
					tempci.Url = url
					tempci.Parent = ci.Parent + SEP_CHARS + strconv.Itoa(ci.Id)
					tempci.Parentname = ci.Parentname + SEP_CHARS + ci.Title
					tempci.Flag = NORMAL
					tempci.Nums = 0
					tempci.Level = ci.Level + 1

					PIS = append(PIS, tempci)
				}

			})
			break
		}
	}

	//能否把价格也解析出来
	return PIS, total, viewUrl
}
