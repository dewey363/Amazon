package main

import (
	//	"errors"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

type ProductListCrawer struct {
	c *UrlCrawer
}

func NewProductListCrawer() *ProductListCrawer {
	params := &CrawerContext{}
	params.loger = InitLog("productlist.log")
	var db DBContext

	db.tableName = "product_list"
	params.db = NewDB(db)
	params.wg = new(sync.WaitGroup)
	params.prefecthPool = OpenRedisPool()
	params.ipPool = OpenRedisPool()
	params.queueName = PRODUCT_LIST_URL_KEY

	crawer := &ProductListCrawer{}
	rule := NewProductListRule()
	crawer.c = NewUrlCrawer(params, rule.parseHook)

	return crawer
}
func (self *ProductListCrawer) Run() {

	self.c.Run()

}

type ProductListRule struct {
	puDB *DB
	sbf  *UrlFilter
}

func NewProductListRule() *ProductListRule {
	rule := &ProductListRule{}
	var context DBContext

	context.tableName = PRODUCT_URL_TALBE
	rule.puDB = NewDB(context)
	rule.sbf = newUrlFilter(rule.puDB)
	return rule
}
func (self *ProductListRule) parseHook(context *CrawerContext, daili string, ci UrlInfo) error {

	html, err := GetHTMLFromProxy(ci.Url, daili)

	if err != nil {
		return err
	}
	if IsRobot(html) {
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
		}).Warnln("验证码页面")
		return ERROR_PROXY
	}
	context.db.SetCompleted(ci.Id)
	if Is404(html) {
		//如果是404则设置完成

		return PAGE_NOT_FOUND
	}

	//分析html，链接的父节点，就是自己的父节点。所有的父亲都指向category表
	cis := self.parase(html, ci)
	//可以实现都没有解析东西出来则写入文件，以便人工分析处理
	if len(cis) == 0 {
		WriteHTMLToFile(ci.Id, html)
	}
	return self.handler(context, ci, cis)
}

//处理解析后得到的链接，有可能直接保存数据库，或者使用缓存
func (self *ProductListRule) handler(context *CrawerContext, pi UrlInfo, PIS []UrlInfo) error {

	//保存到数据库

	var tempcis []UrlInfo
	for _, v := range PIS {
		//如果不存在
		if !self.sbf.TestAndAdd(v.Url) {
			tempcis = append(tempcis, v)
		}
	}

	self.puDB.BatchAdd(tempcis)
	context.loger.WithFields(logrus.Fields{
		"id": pi.Id,
	}).Infof("BatchAdd(product url):%d\n", len(tempcis))
	return nil
}

/*
根据字符串和父类别编号获取产品URL信息
可能返回url数组为空

*/
func (crawer *ProductListRule) parase(body string, ci UrlInfo) []UrlInfo {
	var PIS []UrlInfo
	var utfBody io.Reader = strings.NewReader(body)
	doc, _ := goquery.NewDocumentFromReader(utfBody)

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
	return PIS
}
