package main

import (
	"fmt"
	//	"os"
	"testing"

	//	"github.com/astaxie/beego/config"
)

//测试单个插入
//func TestGetSetting(t *testing.T) {
//	fmt.Println(setting.DefaultInt("http_time_out", 180))
//	fmt.Println(setting.DefaultInt("max_get_time", 4))
//	fmt.Println(setting.DefaultInt("max_procs", 1))
//	fmt.Println(setting.DefaultInt("run_flag", 1))
//	fmt.Println(setting.DefaultInt("prefetch_start", 0))
//	fmt.Println(setting.DefaultInt("use_proxy", 0))

//	fmt.Println(setting.DefaultString("amazon_link", "https://www.amazon.com"))
//	fmt.Println(setting.DefaultString("mysql_db", "root:root@/amazon"))
//	fmt.Println(setting.DefaultString("redis_pool_addr", "127.0.0.1:6379"))
//	fmt.Println(setting.DefaultString("redis_pool_passwd", ""))

//}

//测试单个插入
func TestDDHTML(t *testing.T) {
	//	daili := "528a4f8f113ca:7399ab0dd87e1c40824a46cec4a4c8b3@108.187.149.23:16666"
	//	daili := "528a4f8f113ca:7399ab0dd87e1c40824a46cec4a4c8b3@45.43.212.194:16666"
	//	daili := "989945957cdb2b214abbd21214259565:squid@108.187.149.57:3128"
	daili := "989945957cdb2b214abbd21214259565:squid@108.187.149.91:3128"
	//	url := "https://www.amazon.com/Televisions-Video/b/ref=sd_allcat_tv?ie=UTF8&node=1266092011"
	url := "http://www.ptuit.net/ip.php"
	b, err := DDHTML(daili, url)
	if err != nil {
		fmt.Println(err)
	}
	WriteHTMLToFile(1, string(b))
}
