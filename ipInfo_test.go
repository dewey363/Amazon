package main

import (
	//	"net/http"
	"testing"
)

func TestGetIPsFromFile(t *testing.T) {
	//	ips := GetIPsFromFile("proxy.lib")
	//	t.Logf("ips rows:%d\n", len(ips))
	proxy := &IPInfo{}
	//	ipPool := OpenRedisPool()
	//	proxy = ipPool.PopProxy(10)
	//	if proxy == nil {

	//		t.Log("no")
	//	} else {
	//		t.Log(proxy.Url == "")
	//	}
	t.Log(proxy.Url == "")
}

//func TestSaveIPSToDB(t *testing.T) {
//	SaveIPSToDB("proxy.lib")
//	t.Logf("请查看数据库中是否有数据")
//}

//func TestGetIPIs(t *testing.T) {
//	ipis := GetIPIs(0)
//	t.Logf("rows:%d\n", len(ipis))
//}
//func TestInitIPQueue(t *testing.T) {
//	queue := InitIPQueue()
//	t.Logf("rows:%d\n", queue.Size())
//	item, _ := queue.Get(0)
//	ip := (item).(IPInfo)
//	t.Logf(ip.Url)
//}
//func TestGetIPFromKuaidaili(t *testing.T) {
//	GetIPFromKuaidaili()
//}
