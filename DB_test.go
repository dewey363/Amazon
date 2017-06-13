package main

import (
	"fmt"
	//	"net/http"
	//	"fmt"
	//	"sync"
	"testing"
	//	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//测试单个查询
//func TestDB(t *testing.T) {
//	var context DBContext
//	context.tableName = "category"
//	db := NewDB(context)
//	fmt.Println(db.count(0))

//}

//测试批量插入
//func TestDB(t *testing.T) {
//	start := time.Now()
//	for i := 0; i < 1000; i++ {
//		var pi ProductUrlInfo
//		pi.Flag = 0
//		pi.Parent = 0
//		pi.Url = "https://www.amazon.com/dp/B00MAOCXZ4"
//		InsertProductUrlInfoToDB(pi)
//	}

//	end := time.Now()

//	fmt.Println("执行时间为：" + end.Sub(start).String())
//}
//测试多线程插入
//func TestMultInsert(t *testing.T) {

//	var context DBContext
//	context.tableName = "product_url"
//	db := NewDB(context)
//	wg := new(sync.WaitGroup)
//	start := time.Now()
//	pool, _ := tunny.CreatePool(MAX_PROCS, func(object interface{}) interface{} {
//		pi, _ := object.(UrlInfo)
//		db.SaveProductLists(pi)
//		return ""

//	}).Open()

//	defer pool.Close()

//	wg.Add(1)

//	go func() {
//		for i := 1; i < 100000; i++ {
//			var pi UrlInfo
//			pi.Flag = 0
//			pi.Parent = 0
//			pi.Url = "https://www.amazon.com/dp/B00MAOCXZ4"
//			pool.SendWork(pi)
//		}
//		wg.Done()
//	}()

//	wg.Wait()

//	end := time.Now()

//	fmt.Println("执行时间为：" + end.Sub(start).String())
//	fmt.Println("执行完成")
//}
//func TestDB(t *testing.T) {

//	var context DBContext
//	context.tableName = "product_url"
//	db := NewDB(context)
//	MM := 1000000
//	x1 := db.count(0)
//	x2 := db.count(1)
//	x := x1 + x2
//	tt := x / MM
//	for i := 0; i < tt; i++ {
//		urls, err := db.GetAll(i*MM, MM)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println(len(urls))
//		}
//	}

//}
//func TestPrefetch(t *testing.T) {
//	var context DBContext
//	context.tableName = "category"
//	db := NewDB(context)
//	cis, _ := db.Prefetch(0, 1)
//	fmt.Println(cis[0].Parentname)
//	//	fmt.Println(err)
//}

//测试远程数据库
func TestPrefetch(t *testing.T) {
	fmt.Println(ZENCART_DB)
	orm.RegisterDataBase("db1", "mysql", "root:Aazsedcftgbhujmko123456!!@tcp(45.33.34.189:3306)/amazon?charset=utf8")
	orm.RegisterDataBase("db1", "mysql", "root:Aazsedcftgbhujmko123456!!@tcp(45.33.34.189:3306)/amazon?charset=utf8")
	//	o1 := orm.NewOrm()
	//	o1.Using("db1")
	//	sql1 := "select count(*) as c from product_url"
	//	var maps []orm.Params
	//	num, err := o1.Raw(sql1).Values(&maps)
	//	if err == nil && num > 0 {

	//		fmt.Println(maps[0]["c"].(string))
	//	}
}
