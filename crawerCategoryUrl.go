// Copyright

/*
	Package 包描述
*/
package main

import (
	"fmt"
	"io"

	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

type CategoryUrlCrawer struct {
	c *UrlCrawer
}

func NewCategoryUrlCrawer() *CategoryUrlCrawer {
	params := &CrawerContext{}
	params.loger = InitLog("category.log")
	var db DBContext

	db.tableName = CATEGORY_TABLE
	params.db = NewDB(db)
	params.wg = new(sync.WaitGroup)
	params.prefecthPool = OpenRedisPool()
	params.ipPool = OpenRedisPool()
	params.queueName = CATEGORY_URL_KEY

	crawer := &CategoryUrlCrawer{}
	rule := NewCategoryUrlRule()
	crawer.c = NewUrlCrawer(params, rule.parseHook)

	return crawer
}
func (self *CategoryUrlCrawer) Run() {
	self.c.Run()
}

type CategoryUrlRule struct {
	sbf *UrlFilter
}

func NewCategoryUrlRule() *CategoryUrlRule {
	rule := &CategoryUrlRule{}
	var context DBContext

	context.tableName = CATEGORY_TABLE
	db := NewDB(context)
	rule.sbf = newUrlFilter(db)
	return rule
}
func (self *CategoryUrlRule) parseHook(context *CrawerContext, daili string, ci UrlInfo) error {

	html, err := GetHTMLFromProxy(ci.Url, daili)
	//	fmt.Println(ci.Url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if IsRobot(html) {
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
		}).Warnln("验证码页面")
		//		WriteHTMLToFile(ci.Id, html)
		return ERROR_PROXY
	}
	//如果已经下载下来了就更新下状态，下次不需要再下载
	context.db.SetCompleted(ci.Id)
	if Is404(html) {
		//如果是404则设置完成
		return PAGE_NOT_FOUND
	}

	//分析html
	cis, total := self.parase(html, ci)
	//可以实现都没有解析东西出来则写入文件，以便人工分析处理
	if len(cis) == 0 && total < 10 {
		WriteHTMLToFile(ci.Id, html)
	} else if ci.Level > 2 && len(cis) > 0 && total == 0 {
		//如果是二级目录没有total的话也输出
		WriteHTMLToFile(ci.Id, html)
	}
	return self.handler(context, ci, cis, total)
}

//根据html分析出来的类别总数及产品总数做进一步处理
func (self *CategoryUrlRule) handler(context *CrawerContext, ci UrlInfo, cis []UrlInfo, total int) error {
	ls := len(cis)
	if ls == 0 && total == 0 {
		//如果都没有得到东西
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
		}).Warnln(ERROR_HTML)
		return nil
	}
	//产品总数小于20000，并且大于0
	if total < MAX_PRODUCTS_IN_CATEGORY && total > 0 {
		//保存该链接到PRODUCT_LISTS_DB
		ci.Nums = total
		//保存到数据库

		context.db.SaveProductLists(ci)
		info := fmt.Sprintf("20000>total>0 SaveProductListsToDB: %d\n", len(cis))
		context.loger.WithFields(logrus.Fields{
			"id":    ci.Id,
			"total": total,
		}).Infoln(info)

		//		fmt.Println(info)
		return nil
	}
	//产品总数大于20000
	if ls == 0 {
		//如果没有下一级类别，但是产品数又大于20000，是不是应该保存下来？？？
		//保存该链接到PRODUCT_LISTS_DB或者说是不是有错？？？？？？
		ci.Nums = total
		//保存到数据库
		//		cuDB := NewCategroyUrlDB()
		context.db.SaveProductLists(ci)
		info := fmt.Sprintf("ls=0 and total>20000 SaveProductListsToDB xxxxx: %d", len(cis))
		context.loger.WithFields(logrus.Fields{
			"id":    ci.Id,
			"total": total,
		}).Infoln(info)

		//		fmt.Println(info)

	} else {
		//插入获取到的url到CATEGORY_URL_DB,
		//批量增加之前先判断  如果是total==0 ls>0 是否要判断层次？第三层开始？  ci.parent<100

		if ci.Level >= MAX_LEVEL {
			////这种情况是不是有可能是不需要的链接，直接返回
			//直接添加到product_list，没有去除重复
			//平均分配条数
			var tempcis []UrlInfo
			tempNum := total / (len(cis) + 1)
			for _, v := range cis {
				//如果不存在
				if !self.sbf.TestAndAdd(v.Url) {
					v.Nums = tempNum
					tempcis = append(tempcis, v)
				}
			}
			var context DBContext

			context.tableName = PRODUCT_LISTS_TABLE
			plDB := NewDB(context)
			plDB.BatchAdd(tempcis)
			return nil
		}
		var tempcis []UrlInfo
		for _, v := range cis {
			//如果不存在
			if !self.sbf.TestAndAdd(v.Url) {
				tempcis = append(tempcis, v)
			}
		}
		context.db.BatchAdd(tempcis)

		//记录日志
		info := fmt.Sprintf("ls>0 and total>20000 SaveCategroyUrlInfosToDB: %d", len(cis))
		context.loger.WithFields(logrus.Fields{
			"id":    ci.Id,
			"total": total,
		}).Infoln(info)
		//		fmt.Println(info)

	}
	return nil
}

/*
根据字符串和父类别编号获取下一级URL信息和total总数
可能返回url数组为空，total为0

*/

func (self *CategoryUrlRule) parase(body string, ci UrlInfo) ([]UrlInfo, int) {
	var subCIS []UrlInfo
	var total int = 0

	var utfBody io.Reader = strings.NewReader(body)
	doc, _ := goquery.NewDocumentFromReader(utfBody)

	totalTag := doc.Find("#s-result-count").First()

	//	fmt.Printf("正在处理总数标签:%d\n", len(totalTag.Nodes))

	if len(totalTag.Nodes) > 0 {

		sTotal := totalTag.Text()
		fmt.Println("总数标签内容为" + sTotal)
		total = getTotal(sTotal)

	}
	//使用二级标签
	regexs := [...][2]string{{"div.categoryRefinementsSection", "a"},
		{"div.left_nav", "a"},
		{"div#hybridBrowse", "a.list-item__category-link"},
		{"div#leftNavContainer", "div.a-expander-extend-container a"},
		{"div#hybridBrowse", "div.a-expander-extend-container a"},
		{"div[id^=leftNav]", "div.acs-ln-links a"}, //id以leftNav开头的
		{"div#nav-subnav", "a"},                    //如果很多重复的就不会插入数据库
	} //二维数组的赋值初始化

	for _, v := range regexs {
		categoryTag := doc.Find(v[0]).First()
		if len(categoryTag.Nodes) > 0 {
			categoryTag.Not(".shoppingEngineExpand").Find(v[1]).Each(func(i1 int, content *goquery.Selection) {
				var tempci UrlInfo
				tempci.Title = strings.TrimSpace(content.Text())
				tempci.Title = tempci.Title[0:strings.Index(tempci.Title, "(")] //去除括号及后面的数字

				if strings.Count(tempci.Title, "") > 1 && strings.Index(tempci.Title, "shop") < 0 {
					url, _ := content.Attr("href")
					url = strings.TrimSpace(url)
					//如果连接少于10个字符说明有问题
					if strings.Count(url, "") < 10 || !strings.Contains(url, "http:") {
						tempci.Url = url
						tempci.Parent = ci.Parent + SEP_CHARS + strconv.Itoa(ci.Id)
						tempci.Parentname = ci.Parentname + SEP_CHARS + ci.Title
						tempci.Flag = NORMAL
						tempci.Nums = 0
						tempci.Level = ci.Level + 1
						subCIS = append(subCIS, tempci)
					}
				}

			})
			break
		}
	}
	//如果还是没有,则查找refinements 下面的ul li a中的href 以及a下的 text span class=refinementLink
	if len(subCIS) == 0 {
		leftNav := doc.Find("div#leftNav").First()
		if len(leftNav.Nodes) > 0 {
			leftNav.Find("ul").EachWithBreak(func(i int, content *goquery.Selection) bool {
				content.Find("a").Each(func(i1 int, content1 *goquery.Selection) {
					var tempci UrlInfo
					tempci.Title = strings.TrimSpace(content.Text())
					tempci.Title = tempci.Title[0:strings.Index(tempci.Title, "(")] //去除括号及后面的数字
					if strings.Count(tempci.Title, "") > 1 {
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
							subCIS = append(subCIS, tempci)
						}
					}
				})
				//如果找到则退出
				if len(subCIS) > 0 {
					return false
				}
				return true

			})
		}
	}

	return subCIS, total
}

/*
	根据字符串获取总数 検索結果 380,233件中 1-24件 本
	正则表达式 \d{1,3}(,\d\d\d)*
*/
func getTotal(sTotal string) int {

	var pattern string
	//	pattern = "\\d{1,3}(,\\d\\d\\d)*" //日本
	pattern = "\\d+(\\s+\\d{3})*"
	reg, _ := regexp.Compile(pattern)

	//去除特殊格式
	bs := []byte(sTotal)
	//	fmt.Println(bs)
	//	fmt.Println(string(bs))
	var ns []byte
	for _, v := range bs {
		if v > 0 && v < 128 {
			ns = append(ns, v)
		}
	}
	//	fmt.Println(ns)
	//	fmt.Println(string(ns))
	ss := reg.FindAllString(string(ns), -1)

	if len(ss) == 0 {
		return 0
	}
	//找出最大的就是total
	total := 0
	for _, v := range ss {
		//		fmt.Println(v)
		//去空格，去逗号
		sv := strings.Trim(v, " ")

		sv = strings.Replace(sv, ",", "", 10)
		sv = strings.Replace(sv, " ", "", 10)
		t, err := strconv.Atoi(sv)
		if err == nil {
			if t > total {
				total = t
			}
		}

	}

	return total

}
