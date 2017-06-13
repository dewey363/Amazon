package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	//	"github.com/jeffail/tunny"
)

var pid *logrus.Logger = InitLog("pid.log")

const (
	SEP_STRING   string = "\n**********\n"
	SPACE_NUMBER int    = 1000
)

type ProductUrlCrawer struct {
	c *UrlCrawer
}

func NewProductUrlCrawer() *ProductUrlCrawer {
	params := &CrawerContext{}
	params.loger = InitLog("ProductUrl.log")
	var db DBContext
	db.tableName = "product_url"
	params.db = NewDB(db)
	params.wg = new(sync.WaitGroup)
	params.prefecthPool = OpenRedisPool()
	params.ipPool = OpenRedisPool()
	params.queueName = PRODUCT_INFO_URL_KEY

	crawer := &ProductUrlCrawer{}
	rule := NewProductUrlRule()
	crawer.c = NewUrlCrawer(params, rule.parseHook)

	return crawer
}
func (self *ProductUrlCrawer) Run() {

	self.c.Run()

}

type ProductUrlRule struct {
	//	sis []SiteInfo
}

func NewProductUrlRule() *ProductUrlRule {
	rule := &ProductUrlRule{}
	//	rule.sis = InitSiteInfoFromDB()
	return rule
}
func (self *ProductUrlRule) parseHook(context *CrawerContext, daili string, ci UrlInfo) error {

	html, err := GetHTMLFromProxy(ci.Url, daili)

	//	fmt.Println(ci.Url)
	//	fmt.Println(daili)
	//	WriteHTMLToFile(ci.Id, html)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if IsRobot(html) {
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
		}).Warnln("验证码页面")
		return ERROR_PROXY
	}
	//如果是404则设置完成
	context.db.SetCompleted(ci.Id)
	if Is404(html) {

		return PAGE_NOT_FOUND
	}

	//分析html
	pi := parase(html, ci)
	//可以实现都没有解析东西出来则写入文件，以便人工分析处理
	if pi.Asin == "" {
		WriteHTMLToFile(ci.Id, html)
	}
	return self.hanlder(context, ci, pi)
	return nil
}

//处理解析后得到的链接，有可能直接保存数据库，或者使用缓存
func (self *ProductUrlRule) hanlder(context *CrawerContext, ci UrlInfo, pi ProductInfo) error {

	//没有解析到东西是否要进队列重新获取？
	if pi.Asin == "" {
		context.loger.WithFields(logrus.Fields{
			"id": ci.Id,
		}).Warnln("没有抓取到信息")
		return errors.New("没有抓取到信息")
	}
	context.loger.WithFields(logrus.Fields{
		"id": ci.Id,
	}).Infoln("save product info")
	//发送到数据库
	if strings.Count(pi.Title, "") > 0 &&
		pi.Price != "0" &&
		strings.Count(pi.Desc, "") > 0 ||
		strings.Count(pi.Imgs, "") > 0 ||
		strings.Count(pi.Attr, "") > 0 {
		//		HandlerPI(self.sis, pi)
		//保存采集的数据到本地数据库
		err := HandlerPILocal(pi)
		if err != nil {
			pid.WithFields(logrus.Fields{
				"id": pi.Id,
			}).Warnln(err)
		}
	} else {
		fmt.Println("数据不对")
	}

	//保存到文件
	//每1000个放在一个文件夹中，每天里面有很多个该文件夹
	//20170320/1-1000/1.txt  2.txt
	//20170320/1001-2000/1001.txt  1002.txt
	//20170321/2001-3000/1.txt  2.txt
	//	fileName, _ := getFolder(ci.Id)
	//	content := pi.Asin + SEP_STRING +
	//		pi.Url + SEP_STRING +
	//		fmt.Sprintf("%d", pi.Id) + SEP_STRING +
	//		fmt.Sprintf("%d", pi.Parent) + SEP_STRING +
	//		pi.Title + SEP_STRING +
	//		pi.Price + SEP_STRING +
	//		pi.Attr + SEP_STRING +
	//		pi.Desc + SEP_STRING +
	//		pi.Imgs + SEP_STRING +
	//		pi.others + SEP_STRING

	//	appendToFile(fileName, content)

	return nil
}

/*
	根据id生成文件或者目录,并返回路径
*/
func getFolder(pid int) (string, error) {
	//txt/20170320/1-1000/1.txt  2.txt
	//获取当天日期字符串

	daystr := time.Now().UTC().Add(8 * time.Hour).Format("20060102")

	idstr := fmt.Sprintf("%d-%d", (pid/SPACE_NUMBER*SPACE_NUMBER + 1), (pid/SPACE_NUMBER+1)*SPACE_NUMBER)
	fileName := fmt.Sprintf("%d.txt", pid)
	dir, _ := os.Getwd()
	saveDir := path.Join(dir, "txt", daystr, idstr)
	err := os.MkdirAll(saveDir, os.ModePerm) //生成多级目录
	if err != nil {
		//		fmt.Println(err)
		return "", err
		//		return "", errors.New("输出到文件错误")
	}
	txtFile := path.Join(saveDir, fileName)
	return txtFile, nil
}

/*
根据字符串和父类别编号获取产品URL信息
可能返回url数组为空

*/
func parase(body string, ci UrlInfo) ProductInfo {
	var PI ProductInfo
	PI.Id = ci.Id
	PI.Parent = ci.Parent
	PI.ParentName = ci.Parentname
	PI.Url = ci.Url
	body = ReplaceHTML(body)
	var utfBody io.Reader = strings.NewReader(body)
	doc, _ := goquery.NewDocumentFromReader(utfBody)

	//查找左边的图片
	//查找主图
	//使用二级标签
	regexs := []string{"div#leftcol", "div#aud_left_col"} //二维数组的赋值初始化
	for _, v := range regexs {
		leftCol := doc.Find(v).First()
		mainImage, ok := leftCol.Find("div#imgTagWrapperId img").First().Attr("src")
		if ok {
			PI.MainImage = mainImage
		}
		if len(leftCol.Nodes) > 0 {
			//使用正则表达是查找所有图片链接：不含有逗号的以http开头，中间包含images,结尾是.jpg的字符串  http((?!,).)*images((?!,).)*\.jpg
			html, err := leftCol.Html()

			if err != nil {
				fmt.Println(nil)
			}
			imagesArr := getImages(html)
			PI.Imgs = strings.Join(imagesArr, "|")
			break
		} //end if

	}
	//	cregexs := []string{"div#centercol", "div#aud_center_col"} //二维数组的赋值初始化
	//	for _, v := range cregexs {
	//		centerCol := doc.Find(v).First()
	//		if len(centerCol.Nodes) > 0 {
	//			//title_feature_div  price_feature_div twister_feature_div featurebullets_feature_div 其它的div取值就可以
	//			centerCol.Find("div").Each(func(i int, content *goquery.Selection) {
	//				id, ok := content.Attr("id")
	//				id = strings.TrimSpace(id)
	//				src, _ := content.Html()
	//				//还得去掉javascript代码和所有的html标签
	//				//				cc := ReplaceJs(src)
	//				if ok {
	//					switch id {
	//					//如果是以下标签，先不处理，等下面再处理
	//					case "featurebullets_feature_div":
	//					case "proddetails":
	//					case "pdiframecontent":
	//					case "title_feature_div":
	//					case "price_feature_div":
	//					case "descriptionAndDetails":
	//					case "aud_center_col":
	//					case "buybox_feature_div":
	//					case "buybox":
	//					case "twister_feature_div":
	//						PI.Flag = 1

	//					default: //剩下的合并在一起
	//						PI.others = PI.others + "|" + ReplaceJs(src)
	//					}

	//				}
	//			})
	//			break
	//		} //endif

	//	}

	//description
	desc_regexs := []string{"div#featurebullets_feature_div", "div#proddetails", "div#descriptionAndDetails", "div#pdiframecontent"} //二维数组的赋值初始化
	for _, v := range desc_regexs {
		div := doc.Find(v).First()
		h, _ := div.Html()

		PI.Desc = PI.Desc + " \n " + ReplaceDesc(h)
	}
	//title
	title_regexs := []string{"div#title_feature_div"} //二维数组的赋值初始化
	for _, v := range title_regexs {
		div := doc.Find(v).First()
		h, _ := div.Html()
		PI.Title = ReplaceJs(h)
		if strings.Count(PI.Title, "") > 10 {
			break
		}
		//		PI.Price = PI.Price + " \n " + ReplaceJs(h)
	}
	//price
	price_regexs := []string{"div#price_feature_div", "div#buybox", "div#buybox_feature_div", "div#aud_center_col"} //二维数组的赋值初始化
	tmpPrice := ""
	for _, v := range price_regexs {
		div := doc.Find(v).First()
		h, _ := div.Html()

		tmpPrice = tmpPrice + " \n " + ReplaceJs(h)
	}
	//字符串价格转数值
	PI.Price = getPrice(tmpPrice)
	if PI.Price == "0" {
		appendToFile(fmt.Sprintf("./test/price/%d.txt", PI.Id), tmpPrice)
	}
	//属性的话，获取里面的option值
	feature_regexs := []string{"div#twister_feature_div"}
	for _, v := range feature_regexs {
		bflag := false
		div := doc.Find(v).First()
		if len(div.Nodes) > 0 {
			div.Find("option").Each(func(i int, content *goquery.Selection) {
				bflag = true //退出标志
				//idstr, _ := content.Attr("id")
				//				if idstr != "-1" {

				//				}
				PI.Attr = PI.Attr + SEP_CHARS + strings.TrimSpace(content.Text())

			})
		}
		if bflag {
			break
		}

	}
	appendToFile(fmt.Sprintf("./test/feature/%d.txt", PI.Id), PI.Attr)
	//asin 页面中的imput 中name=asin的值 <input type="hidden" id="ASIN" name="ASIN" value="B001CS8A6A"> <form method="post" id="addToCart"
	asinTag := doc.Find("input#asin").First()
	if len(asinTag.Nodes) > 0 {
		asin, ret := asinTag.Attr("value")
		//		fmt.Println(asin)
		if ret {
			PI.Asin = asin
		}

	}
	return PI
}

//从html中根据正则表达式获取图片路径，需要去除重复项
func getImages(html string) []string {
	var pattern string
	pattern = "http([^,].)*images([^,].)*\\.jpg"
	reg, err := regexp.Compile(pattern)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	ss := reg.FindAllString(html, -1)

	return removeDuplicatesAndEmpty(ss)
}

//从字符中获取价格
//123,33
func getPrice(strPrice string) string {
	var price float64 = 0.0

	var pattern string
	//	pattern = "\\d{1,3},\\d\\d" //法国最多999.99eur
	pattern = "\\d{1,3},\\d\\d"
	reg, _ := regexp.Compile(pattern)

	//去除特殊格式
	bs := []byte(strPrice)

	var ns []byte
	for _, v := range bs {
		if v > 0 && v < 128 {
			ns = append(ns, v)
		}
	}

	ss := reg.FindAllString(string(ns), -1)

	if len(ss) == 0 {
		return "0"
	}
	//找出最大的就是 price

	for _, v := range ss {
		//		fmt.Println(v)
		//去空格，去逗号
		sv := strings.Trim(v, " ")

		sv = strings.Replace(sv, ",", ".", 10)
		sv = strings.Replace(sv, " ", "", 10)
		f64, err := strconv.ParseFloat(sv, 64)
		if err == nil {
			if f64 > price {
				price = f64
			}
		}

	}

	return strconv.FormatFloat(price, 'f', 2, 64)
}

/**
 * 数组去重 去空
 */
func removeDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

//去除html中的css,连续换行符，所有标签转换为小写
func ReplaceHTML(src string) string {

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("<style[\\S\\s]+?</style>")
	src = re.ReplaceAllString(src, "")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "")
	src = strings.TrimSpace(src)
	return src
}
func ReplaceJs(src string) string {
	//去除SCRIPT
	re, _ := regexp.Compile("<script[\\S\\s]+?</script>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllString(src, "")

	src = strings.TrimSpace(src)

	var utfBody io.Reader = strings.NewReader(src)
	doc, _ := goquery.NewDocumentFromReader(utfBody)
	return doc.Text()
}

//不需要去除html
func ReplaceDesc(src string) string {
	//去除SCRIPT
	re, _ := regexp.Compile("<script[\\S\\s]+?</script>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	//	re, _ = regexp.Compile("<[\\S\\s]+?>")
	//	src = re.ReplaceAllString(src, "")

	src = strings.TrimSpace(src)

	var utfBody io.Reader = strings.NewReader(src)
	doc, _ := goquery.NewDocumentFromReader(utfBody)
	return doc.Text()
}

//测试输出解析到的数据到文件
func doutFile(pi ProductInfo) {
	appendToFile(fmt.Sprintf("./test/title/%d.txt", pi.Id), pi.Title)

	appendToFile(fmt.Sprintf("./test/price/%d.txt", pi.Id), pi.Price)

	appendToFile(fmt.Sprintf("./test/images/%d.txt", pi.Id), pi.Imgs)

	appendToFile(fmt.Sprintf("./test/feature/%d.txt", pi.Id), pi.Attr)

	appendToFile(fmt.Sprintf("./test/desc/%d.txt", pi.Id), pi.Desc)

	appendToFile(fmt.Sprintf("./test/others/%d.txt", pi.Id), pi.others)
	appendToFile(fmt.Sprintf("./test/asin/%d.txt", pi.Id), pi.Asin)

}

//插入本地数据库处理产品数据
func HandlerPILocal(pi ProductInfo) error {
	return PostProductLocal(pi)
}

//处理产品数据
func HandlerPI(sis []SiteInfo, pi ProductInfo) {
	//对所有的网站

	for _, si := range sis {
		fmt.Printf("当前目标网站为：%d\n", si.Id)
		res := si.IsAdd(pi)
		if res > 0 {
			//该网站就插入数据了，退出循环
			if res == 1 {
				addSiteDetail(si.Id, pi.Parent, res)
			} else {
				updateSiteDetail(si.Id, pi.Parent, res)
			}
			PostProduct(si, pi)
			//			if ok {
			//				break
			//			}
			break
		}
	}

}
