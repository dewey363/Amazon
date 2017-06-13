package main

import (
	//	"fmt"
	"os"
	//	"path/filepath"
	"strconv"

	"github.com/astaxie/beego/config"
)

// 配置文件涉及的默认配置。
const (
	CONFIG               string = "config.ini"
	max_level            int    = 4
	http_time_out        int    = 180 //链接超时，单位秒
	failed_proxy_seconds int    = 30 * 60
	max_get_time         int    = 4 * 60 //每个链接时间，单位秒数
	mysql_db             string = "root:root@/amazon2"
	redis_pool_addr      string = "127.0.0.1:6379"
	redis_pool_passwd    string = "" //redis密码
	max_procs            int    = 51 //多线程数
	run_flag             int    = 1  //1采集类别 2采集链接 3.采集产品
	prefetch_start       int    = 0  //数据开始位置
	//	use_proxy            int    = 1    //是否使用代理
	lang string = "en" //en 美国 jp 日本

	//	log_server string = "45.33.34.189:12202"

	//	local    string = "测试"
	zencart_db string = "root:root@/zencart154"
	log_time   int    = 30
)

var setting = func() config.Configer {

	iniconf, err := config.NewConfig("ini", CONFIG)
	if err != nil {
		file, err := os.Create(CONFIG)
		file.Close()
		iniconf, err = config.NewConfig("ini", CONFIG)
		if err != nil {
			panic(err)
		}
		defaultConfig(iniconf)
		iniconf.SaveConfigFile(CONFIG)
	} else {
		trySet(iniconf)
	}

	return iniconf
}()

func defaultConfig(iniconf config.Configer) {
	iniconf.Set("http_time_out", strconv.Itoa(http_time_out))
	iniconf.Set("max_get_time", strconv.Itoa(max_get_time))
	iniconf.Set("max_procs", strconv.Itoa(max_procs))
	iniconf.Set("run_flag", strconv.Itoa(run_flag))
	iniconf.Set("prefetch_start", strconv.Itoa(prefetch_start))
	//	iniconf.Set("use_proxy", strconv.Itoa(use_proxy))
	iniconf.Set("failed_proxy_seconds", strconv.Itoa(failed_proxy_seconds))
	iniconf.Set("max_level", strconv.Itoa(max_level))

	iniconf.Set("mysql_db", mysql_db)
	iniconf.Set("redis_pool_addr", redis_pool_addr)
	iniconf.Set("redis_pool_passwd", redis_pool_passwd)
	iniconf.Set("lang", lang)
	//	iniconf.Set("log_server", log_server)
	iniconf.Set("zencart_db", zencart_db)
	iniconf.Set("log_time", strconv.Itoa(log_time))

}

func trySet(iniconf config.Configer) {
	if v, e := iniconf.Int("http_time_out"); v <= 0 || e != nil {
		iniconf.Set("http_time_out", strconv.Itoa(http_time_out))
	}
	if v, e := iniconf.Int("max_level"); v <= 0 || e != nil {
		iniconf.Set("max_level", strconv.Itoa(max_level))
	}
	if v, e := iniconf.Int("max_get_time"); v <= 0 || e != nil {
		iniconf.Set("max_get_time", strconv.Itoa(max_get_time))
	}
	if v, e := iniconf.Int("max_procs"); v <= 0 || e != nil {
		iniconf.Set("max_procs", strconv.Itoa(max_procs))
	}
	if v, e := iniconf.Int("run_flag"); v <= 0 || e != nil {
		iniconf.Set("run_flag", strconv.Itoa(run_flag))
	}
	//	if v, e := iniconf.Int("prefetch_start"); v <= 0 || e != nil {
	//		iniconf.Set("prefetch_start", strconv.Itoa(prefetch_start))
	//	}
	//	if v, e := iniconf.Int("use_proxy"); v <= 0 || e != nil {
	//		iniconf.Set("use_proxy", strconv.Itoa(use_proxy))
	//	}
	if v, e := iniconf.Int("log_time"); v <= 0 || e != nil {
		iniconf.Set("log_time", strconv.Itoa(log_time))
	}
	if v, e := iniconf.Int("failed_proxy_seconds"); v <= 0 || e != nil {
		iniconf.Set("failed_proxy_seconds", strconv.Itoa(failed_proxy_seconds))
	}

	if v := iniconf.String("mysql_db"); v == "" {
		iniconf.Set("mysql_db", mysql_db)
	}
	if v := iniconf.String("lang"); v == "" {
		iniconf.Set("lang", lang)
	}

	if v := iniconf.String("redis_pool_addr"); v == "" {
		iniconf.Set("redis_pool_addr", redis_pool_addr)
	}
	if v := iniconf.String("zencart_db"); v == "" {
		iniconf.Set("zencart_db", zencart_db)
	}
	//	if v := iniconf.String("log_server"); v == "" {
	//		iniconf.Set("log_server", log_server)
	//	}

	//	if v := iniconf.String("local"); v == "" {
	//		iniconf.Set("local", local)
	//	}

	iniconf.SaveConfigFile(CONFIG)
}
