// Copyright

/*
	Package 包描述
*/
package main

import (
	"errors"
	"fmt"

	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jeffail/tunny"
)

var (
	ERROR_HTML     error = errors.New("返回的html有问题，都没有得到东西")
	ERROR_PROXY    error = errors.New("代理有问题")
	PAGE_NOT_FOUND error = errors.New("404页面")
	ERROR_GET_HTML error = errors.New("获取网页失败")
)

type CrawerContext struct {
	loger        *logrus.Logger
	db           *DB
	wg           *sync.WaitGroup
	prefecthPool *redisPool
	ipPool       *redisPool
	workPool     *tunny.WorkPool
	queueName    string
}
type UrlCrawer struct {
	context   *CrawerContext
	ParseFunc func(*CrawerContext, string, UrlInfo) error
}

func NewUrlCrawer(params *CrawerContext, parse func(*CrawerContext, string, UrlInfo) error) *UrlCrawer {
	crawer := &UrlCrawer{
		context:   params,
		ParseFunc: parse,
	}

	return crawer
}
func (crawer *UrlCrawer) Run() {
	start := time.Now()
	fmt.Printf("开始时间为：%v", start)
	//1.category进队列
	crawer.initPool()
	//2.启动定时器，定时检查数据
	crawer.context.wg.Add(1)
	go crawer.initPoolTimer()
	//3.启动定Ip进队列，暂时没有
	if USE_PROXY == 1 {
		InitIPPool()
	}
	//4.启动日志发送，暂时关闭
	go crawer.sendLog()
	//5.开始采集
	crawer.context.workPool = crawer.startCrawerWorker()
	defer crawer.context.workPool.Close()
	//6.结束
	crawer.context.wg.Wait()

	end := time.Now()

	crawer.context.loger.Println("执行时间为：" + end.Sub(start).String())
	crawer.context.loger.Println("执行完成")

}

//	初始化category url队列
//
func (crawer *UrlCrawer) initPool() {

	//进队列
	cis, _ := crawer.context.db.Prefetch(NORMAL, PreFetchSize)
	fmt.Printf("initCUIPool进队列%d:\n", len(cis))
	crawer.context.prefecthPool.EmptyQueue(crawer.context.queueName)
	for _, ci := range cis {
		crawer.context.prefecthPool.PushURL(ci, crawer.context.queueName)
	}

}

//	填充category url队列，如果数据库中没有数据，则返回false
//
func (crawer *UrlCrawer) fillPool() bool {
	//获取shell信息
	rPool := OpenRedisPool()
	defer CloseRedisPool(rPool)
	//如果大于最小则不从数据库读取
	if rPool.LenQueue(crawer.context.queueName) >= Min_Queue_Size {
		return true
	}
	//进队列
	//	cuDB := NewCategroyUrlDB()
	cis, _ := crawer.context.db.Prefetch(NORMAL, PreFetchSize)
	if len(cis) <= 0 {
		return false
	}
	for _, ci := range cis {
		//		fmt.Println(ci.Url)
		rPool.PushURL(ci, crawer.context.queueName)
	}
	return true
}

//	定时读取数据到队列
//	每个IP在4(MAX_GET_TIME)分钟内只能处理一条，50*15/小时。
//  启动的时候读取Proc_Fetch_Size(19),定时4分钟，查询队列，如果小于(3)则从数据库读取19
func (crawer *UrlCrawer) initPoolTimer() {
	i := 0
	for {
		start := time.Now()
		//如果数据库没有了，则退出循环
		fmt.Println("initCUIPoolTimer定时读取数据")
		if !crawer.fillPool() {
			i = i + 1
			if i >= Proc_Fetch_Size {
				break
			}

		} else {
			i = 0
		}

		WaitSecond(start, PreFetchTime)

	}
	crawer.context.wg.Done()

}

/*
	开始采集工作
*/
func (crawer *UrlCrawer) startCrawerWorker() *tunny.WorkPool {

	pool, _ := tunny.CreatePool(MAX_PROCS, func(object interface{}) interface{} {
		//		ci, _ := object.(*CategroyUrlInfo)
		threadId := object.(int)
		fmt.Printf("开始第%d个工作线程\n", threadId)
		var ci *UrlInfo
		var n time.Time
		var err error
		for {

			fmt.Println("开始处理数据...")

			//等4分钟如果没有的话，则退出？？？
			ci = crawer.context.prefecthPool.PopURL(crawer.context.queueName, MAX_GET_TIME)
			if ci == nil {
				crawer.context.loger.Warnln("没有category url 可用")
				return ""
			}
			pn := time.Now()
			proxy := &IPInfo{}
			if USE_PROXY == 1 {
				proxy = crawer.context.ipPool.PopProxy(MAX_GET_TIME)

				if proxy == nil {
					crawer.context.loger.Warnln("没有proxy可用")
					return ""
				}

			}
			proxyTake := int(time.Now().Sub(pn).Seconds())

			fmt.Printf("启动线程处理%d\n", ci.Id)
			n = time.Now()
			err = crawer.craw(*ci, proxy.Url)

			//把URL重新入队，不需要重新入队，数据库中没有设置完成标志，下次自动读出
			if USE_PROXY == 1 {
				//设置代理失败时间，如果等待一段时间后才可以使用
				if err == ERROR_PROXY {
					proxy.Flag = 1

				} else {
					proxy.Flag = 0
				}
				proxy.FailedTime = time.Now()
				crawer.context.ipPool.PushProxy(*proxy)

			}
			//计算每条记录处理多长时间，保存到数据库中，代理失败标记也保存在数据库，可以知道代理失败的概率
			takeTime := int(time.Now().Sub(n).Seconds())
			crawer.context.db.AddTakeTime(ci.Id, takeTime, proxyTake, proxy.Flag, proxy.Url)
			fmt.Printf("处理完%d\n", ci.Id)

			WaitSecond(n, RandSecond(MAX_GET_TIME))
		}

		return ""

	}).Open()

	//	defer pool.Close()

	for i := 0; i < MAX_PROCS; i++ {
		crawer.context.wg.Add(1)
		go func(i int) {
			pool.SendWork(i)
			fmt.Printf("启动%d个进程\n", MAX_PROCS)
			crawer.context.wg.Done()
		}(i)
	}
	fmt.Println("退出函数")

	return pool

}

// 执行采集分析类别url工作

func (crawer *UrlCrawer) craw(ci UrlInfo, daili string) error {
	//把获取网页的暂时移到外部

	//进行后续处理
	return crawer.ParseFunc(crawer.context, daili, ci)

}

/*
	发送日志：
	需要日志服务器IP及端口,本机IP,初始化的时候已经设置
	包含内容：当前进程数，代理IP数，待采集数量，已采集数量
	//category表：flag=0,flag=1,level00,level01,level10,level11,level20,level21
	//product_lists表:flag=0,flag=1,sum(nums)
	//product_list表:flag=0,flag=1
	//product_url表:flag=0,flag=1
*/
func (crawer *UrlCrawer) sendLog() {

	for {
		start := time.Now()
		var workers, ips, x, y int

		if crawer.context.workPool != nil {
			workers = crawer.context.workPool.NumWorkers()
		}
		if crawer.context.ipPool != nil {
			ips = crawer.context.ipPool.LenQueue(PROXY_KEY)
		}
		if crawer.context.db != nil {
			x = crawer.context.db.count(0)
			y = crawer.context.db.count(1)
		}

		//		logrus.WithFields(logrus.Fields{
		//			"workers":     strconv.Itoa(workers),
		//			"ips":         strconv.Itoa(ips),
		//			"uncompleted": strconv.Itoa(x),
		//			"completed":   strconv.Itoa(y),
		//		}).Infoln(time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02 15:04:05"))
		if crawer.context.db != nil {
			crawer.context.db.SaveLog(workers, ips, x, y, time.Now().UTC().Add(8*time.Hour).Format("2006-01-02 15:04:05"))

		}

		WaitSecond(start, LOG_TIME*60) //30分钟发送一次
	}

}
