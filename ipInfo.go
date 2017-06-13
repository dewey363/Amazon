package main

import (
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type IPInfo struct {
	Id         int
	Url        string
	Flag       int
	FailedTime time.Time
}

/*
根据IP信息，初始化代理池
*/
func InitIPPool() {
	ipis := GetIPisFromFile()
	rPool := OpenRedisPool()
	//先清空之前的代理池
	rPool.EmptyQueue(PROXY_KEY)
	rPool.EmptyQueue(PROXY_FAILED_KEY)
	defer CloseRedisPool(rPool)
	for _, ip := range ipis {
		rPool.PushProxy(ip)
	}

}

/*
按行读取文件的代理IP

*/
func GetIPisFromFile() []IPInfo {
	ips := ReadFileLines("proxy.lib")
	var ipis []IPInfo
	for i, v := range ips {
		var ip IPInfo
		ip.Id = i
		ip.Url = string(v)
		ip.Url = strings.TrimSpace(strings.TrimRight(strings.TrimRight(ip.Url, "\n"), "\r"))
		ip.Flag = 0
		ip.FailedTime = time.Now()
		ipis = append(ipis, ip)
	}
	return ipis
}
