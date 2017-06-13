package main

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDataBase("default", "mysql", MYSQL_DB, 30)

	orm.RegisterDataBase("zencart", "mysql", ZENCART_DB, 30)
	orm.RegisterModel(new(UrlInfo))

}

//链接
type UrlInfo struct {
	Id         int
	Title      string
	Url        string
	Flag       int
	Parent     string
	Parentname string
	Nums       int
	Level      int
}

type DBContext struct {
	tableName string
}
type DB struct {
	context DBContext
	//	sbf     *UrlFilter
}

//初始化产品链接数据库
func NewDB(params DBContext) *DB {
	db := &DB{
		context: params,
	}
	return db
}

/*
从数据库中读取一级目录flag为0的，存放在map中 map[int]CategroyUrlInfo
*/
func (self *DB) Prefetch(flag int, prefetch int) ([]UrlInfo, error) {

	var cis []UrlInfo

	sql := "SELECT * FROM " + self.context.tableName + " where flag=? order by id  limit " + strconv.Itoa(PREFETCH_START) + "," + strconv.Itoa(prefetch)
	//	x := generateRandomNumber(PREFETCH_START, 4650000, prefetch)
	//	where := "(" + strings.Join(x, ",") + ")"
	//	sql := "SELECT * FROM " + self.context.tableName + " where flag=? and id in " + where

	o := orm.NewOrm()
	_, err := o.Raw(sql, flag).QueryRows(&cis)
	if err != nil {
		panic(err)
	}
	return cis, nil

}

//每次最好只读取1百万条数据
func (self *DB) GetAll(start int, nums int) ([]UrlInfo, error) {

	var cis []UrlInfo

	sql := "SELECT * FROM " + self.context.tableName + " limit " + strconv.Itoa(start) + "," + strconv.Itoa(nums)
	o := orm.NewOrm()
	_, err := o.Raw(sql).QueryRows(&cis)
	if err != nil {
		panic(err)
	}
	return cis, nil

}

/*
设置完成标志
*/
func (self *DB) SetCompleted(id int) {

	sql := "UPDATE " + self.context.tableName + " SET flag=1 WHERE id=?"
	o := orm.NewOrm()
	_, err := o.Raw(sql, id).Exec()
	checkErr(err)
}

/*
批量更新类别数据库中的category表的flag标记
*/
func (self *DB) BatchCompleted(flag int, cis []UrlInfo) {
	if len(cis) <= 0 {
		return
	}

	where := " id in("
	for _, v := range cis {
		where = where + strconv.Itoa(v.Id) + ","
	}
	//删除最后一个逗号
	temp := string(where[0 : strings.Count(where, "")-2])
	where = temp + ")"

	updateSql := "UPDATE " + self.context.tableName + " SET flag=? WHERE " + where
	o := orm.NewOrm()
	_, err := o.Raw(updateSql, flag).Exec()
	checkErr(err)
}

//通过id编号获取数据中category url
func (self *DB) GetById(id int) (UrlInfo, error) {
	var ci UrlInfo

	sql := "SELECT * FROM " + self.context.tableName + " where id=? order by id limit 0,1"

	o := orm.NewOrm()
	err := o.Raw(sql, id).QueryRow(&ci)
	checkErr(err)
	return ci, nil
}

/*
保存单个类别 url 信息到产品列表数据库中的product_lists表
*/
func (self *DB) SaveProductLists(ci UrlInfo) {

	sql := "INSERT INTO product_lists(flag,title,url,parent,parentname,nums,level)VALUES(?,?,?,?,?,?,?)"
	o := orm.NewOrm()
	_, err := o.Raw(sql, ci.Flag, ci.Title, ci.Url, ci.Parent, ci.Parentname, ci.Nums, ci.Level).Exec()
	checkErr(err)

}

/*
插入单个个类别 url 信息到数据库中的 ProductURL 表
*/
func (self *DB) SaveProductURL(pi UrlInfo) {

	sql := "INSERT INTO product_url(url,flag,parent)VALUES(?,?,?,?)"
	o := orm.NewOrm()
	_, err := o.Raw(sql, pi.Url, pi.Flag, pi.Parent, pi.Parentname).Exec()
	checkErr(err)

}

/*
插入多个类别url信息到数据库中的category表
*/
func (self *DB) BatchAdd(cis []UrlInfo) {
	//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

	sql := "INSERT INTO " + self.context.tableName + "(title,url,flag,parent,parentname,nums,level)VALUES(?,?,?,?,?,?,?)"

	o := orm.NewOrm()
	p, err := o.Raw(sql).Prepare()
	defer p.Close()
	checkErr(err)
	for _, ci := range cis {
		p.Exec(ci.Title, ci.Url, ci.Flag, ci.Parent, ci.Parentname, ci.Nums, ci.Level)
	}

	//	p.Close() // 别忘记关闭 statement

}

//判断有几个
func (self *DB) count(flag int) int {

	sql := "SELECT count(*) as c FROM " + self.context.tableName + " WHERE flag = ?"

	urlNum := 0
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(sql, flag).Values(&maps)
	if err == nil && num > 0 {

		urlNum, _ = strconv.Atoi(maps[0]["c"].(string))
	}

	return urlNum
}

//判断从文件中读取总类别，插入到数据库中
func (self *DB) initDBfromFile(filename string) {
	items := ReadFileLines(filename)

	var cis []UrlInfo
	for _, item := range items {
		var ci UrlInfo
		vals := strings.Split(item, "|")

		if len(vals) == 2 {
			ci.Flag = 0
			ci.Nums = 0
			ci.Parent = "0"
			ci.Parentname = ""
			ci.Level = 0
			ci.Title = strings.TrimSpace(vals[0])
			ci.Url = strings.TrimSpace(vals[1])
			cis = append(cis, ci)
		}

	}
	if len(cis) > 0 {
		self.BatchAdd(cis)
	}
}
func (self *DB) SaveLog(workers int, ips int, uncompleted int, completed int, logtime string) {
	sql := "INSERT INTO log(workers,ips,uncompleted,completed,logtime)VALUES(?,?,?,?,?)"

	o := orm.NewOrm()
	_, err := o.Raw(sql, workers, ips, uncompleted, completed, logtime).Exec()
	checkErr(err)

}
func (self *DB) AddTakeTime(urlid int, seconds int, proxyTake int, flag int, url string) {
	sql := "INSERT INTO taketime(urlid,seconds,proxytake,proxyflag,url)VALUES(?,?,?,?,?)"

	o := orm.NewOrm()
	_, err := o.Raw(sql, urlid, seconds, proxyTake, flag, url).Exec()
	checkErr(err)

}
