package main

import (
	//	"fmt"
	//	"net/http"
	"testing"
)

//func TestPushShell(t *testing.T) {
//	rPool := OpenRedisPool()
//	var shell ShellInfo
//	shell.Id = 1
//	shell.Url = "http://baidu.com"
//	shell.Flag = 0
//	rPool.PushShell(shell)
//}

//func TestPopShell(t *testing.T) {
//	rPool := OpenRedisPool()
//	for {
//		shell := rPool.PopShell(2)
//		if shell == nil {
//			break
//		}
//		t.Log(shell.Id)
//		t.Log(shell.Url)
//		t.Log(shell.Flag)
//	}
//	fmt.Println("测试完毕")

//}
//func TestPushShell(t *testing.T) {
//	rPool := OpenRedisPool()
//	var shell ShellInfo
//	shell.Id = 1
//	shell.Url = "http://baidu.com"
//	shell.Flag = 0
//	rPool.PushShell(shell)
//}
//func TestPushProxy(t *testing.T) {
//	rPool := OpenRedisPool()
//	ips := GetIPisFromFile()
//	for _, ip := range ips {
//		rPool.PushProxy(ip)
//	}

//}
func TestLenPool(t *testing.T) {
	rPool := OpenRedisPool()
	ips := GetIPisFromFile()
	for _, ip := range ips {
		rPool.PushProxy(ip)
	}
	t.Logf("%d", rPool.LenQueue(PROXY_KEY))

}
