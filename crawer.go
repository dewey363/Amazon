package main

import (
	"errors"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/go_tool/spider"
)

//使用fasthttp,使用代理从amazon获取网页数据
//代理从函数变量中获取
func GetHTMLFromProxy(curl string, daili string) (string, error) {

	//构造访问网址
	var desturl string
	if strings.HasPrefix(curl, "http") {
		desturl = curl
	} else {
		desturl = AMAZON_LINK + curl
	}

	//	desturl := "htts://www.baidu.com"
	//	url = "https://www.amazon.com/Televisions-Video/b/ref=sd_allcat_tv?ie=UTF8&node=1266092011"
	//	fmt.Println("正在获取的url是：" + desturl)

	if daili != "" {

		content, err := DDHTML(daili, desturl)
		html := string(content)

		return html, err

	}
	return "", errors.New("没有使用代理不能下载")
}
func DDHTML(ip string, url string) ([]byte, error) {
	proxy := "http://" + ip
	browser, _ := spider.NewSpider(proxy)
	browser.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	browser.Header.Set("Accept-Language", "en-US;q=0.8,en;q=0.5")
	browser.Header.Set("Connection", "keep-alive")
	if strings.Contains(url, "www.amazon.co.jp") {
		browser.Header.Set("Host", "www.amazon.co.jp")
	} else if strings.Contains(url, "www.amazon.de") {
		browser.Header.Set("Host", "www.amazon.de")
	} else if strings.Contains(url, "www.amazon.co.uk") {
		browser.Header.Set("Host", "www.amazon.co.uk")
	} else if strings.Contains(url, "www.amazon.fr") {
		browser.Header.Set("Host", "www.amazon.fr")
	} else {
		browser.Header.Set("Host", "www.amazon.com")
	}
	browser.Header.Set("Upgrade-Insecure-Requests", "1")
	browser.Header.Set("User-Agent", GetRandUserAgent())
	browser.Url = url
	content, err := browser.Get()
	return content, err
}

func IsRobot(content string) bool {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
	text := doc.Find("title").Text()
	// uk usa
	if strings.Contains(text, "Robot Check") {
		return true
	}
	//jp
	if strings.Contains(text, "CAPTCHA") {
		return true
	}
	//de
	if strings.Contains(text, "Bot Check") {
		return true
	}
	return false
}

func Is404(content string) bool {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
	text := doc.Find("title").Text()
	if strings.Contains(text, "Page Not Found") {
		return true
	}
	if strings.Contains(text, "404") {
		return true
	}
	//uk
	if strings.Contains(string(content), "The Web address you entered is not a functioning page on our site") {
		return true
	}
	//de
	if strings.Contains(string(content), "Suchen Sie bestimmte Informationen") {
		return true
	}
	if strings.Contains(string(content), "Suchen Sie etwas bestimmtes") {
		return true
	}
	return false
}
