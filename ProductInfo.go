package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	//	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//产品结构
type ProductInfo struct {
	Id         int
	Asin       string
	Title      string
	Url        string
	Flag       int
	Parent     string
	ParentName string
	Price      string
	Attr       string
	Desc       string
	MainImage  string
	Imgs       string
	others     string
}

//网站数据结构
type SiteInfo struct {
	Id     int
	Domain string
	DBstr  string
	Total  int
	Flag   int
	CNums  int
	CList  map[string]int
	mux    sync.Mutex //支持多线程
}

//初始化
func NewSiteInfo() *SiteInfo {
	psi := &SiteInfo{
		CList: make(map[string]int),
	}
	return psi
}

//判断是否增加
//所有flag=0的网站循环
//if 总数<200000 {
//		if leibie 存在 {
//			插入数据//可以使用异步，或者返回结果，因为速度比较慢，没采集1条有5分钟的时间，完全够时间操作
//			总数+1
//		}else{
//			if leibie<100 {
//				插入数据
//				总数+1
//			}else{
//				退出，下一个网站判断
//			}
//		}
//	}else{
//		设置网站flag=1
//	}
//返回n，则更新数据库中该网站的类别数目n
//返回1，则插入该网站的类别，数目为1
//返回-1，则不进行任何操作，继续下一个网站
//返回-2，设置网站flag，不对数据库进行任何操作，继续下一个网站
func (self *SiteInfo) IsAdd(pi ProductInfo) int {
	self.mux.Lock()
	defer self.mux.Unlock() //这样操作行不行？？？？？？？
	if self.Total < MAX_PRODUCTS_IN_SITE {
		if self.exist(pi.Parent) {
			self.CList[pi.Parent] = self.CList[pi.Parent] + 1
			self.Total = self.Total + 1
			return self.CList[pi.Parent]
		} else {
			if self.CNums < MAX_CATEGORYS_IN_SITE {
				self.CList[pi.Parent] = 1
				self.Total = self.Total + 1
				return 1
			} else {
				return -1
			}
		}
	} else {
		self.Flag = 1
		return -2
	}
}

//判断是否存在
func (self *SiteInfo) exist(key string) bool {
	_, ok := self.CList[key]
	return ok
}

//从数据库中初始化多个siteInfo
//数据库表site  id,domain, host,username,passwd,dbname,flag
//数据库表sitedetail id siteid cname,nums
func InitSiteInfoFromDB() []SiteInfo {
	fmt.Println("初始化siteinfo")
	var sis []SiteInfo
	//读取site表
	sql := "select * from site where flag=0"
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(sql).Values(&maps)
	if err == nil && num > 0 {
		for _, v := range maps {

			id, _ := strconv.Atoi(v["id"].(string))
			domain, _ := v["domain"].(string)
			host, _ := v["host"].(string) //ip地址：格式为：96.126.123.38
			username, _ := v["username"].(string)
			passwd, _ := v["passwd"].(string)
			dbname, _ := v["dbname"].(string)
			flag, _ := strconv.Atoi(v["flag"].(string))

			si := NewSiteInfo()
			si.Id = id
			si.Domain = domain
			si.Flag = flag
			//root:root@/amazonfr
			//orm.RegisterDataBase("db1", "mysql", "root:yZpmTXkuSKdKbYGn@tcp(96.126.123.38:3306)/zadmin_menhirshop?charset=utf8")

			si.DBstr = username + ":" + passwd + "@tcp(" + host + ":3306)/" + dbname + "?charset=utf8"
			si.Total = 0
			si.CNums = 0
			sis = append(sis, *si)

		}

	}
	//设置sis
	//读取数据库sitedetail
	var maps2 []orm.Params
	sql = "select * from sitedetail where siteid=?"
	for i, _ := range sis {
		num, err := o.Raw(sql, sis[i].Id).Values(&maps2)
		if err == nil && num > 0 {
			for _, v := range maps2 {
				cname := v["cname"].(string)
				cnums, _ := strconv.Atoi(v["cnums"].(string))
				sis[i].Total = sis[i].Total + cnums
				sis[i].CList[cname] = cnums
				sis[i].CNums = sis[i].CNums + 1
			}
			if sis[i].Total >= MAX_PRODUCTS_IN_SITE {
				sis[i].Flag = 1
				SetDomainCompleted(sis[i].Id)
			}

		}
	}

	//注册数据库
	for _, si := range sis {
		db := fmt.Sprintf("db%d", si.Id)
		fmt.Println(si.DBstr)
		orm.RegisterDataBase(db, "mysql", si.DBstr)

	}
	fmt.Println(sis)
	return sis
}

//设置域名完成
func SetDomainCompleted(siteid int) {
	sql := "update site set flag=1 where id=?"
	o := orm.NewOrm()
	_, err := o.Raw(sql, siteid).Exec()
	if err != nil {
		fmt.Println(err)
	}
}

//addSiteDetail
func addSiteDetail(siteid int, cname string, cnums int) {
	sql := "INSERT INTO sitedetail(siteid,cname,cnums)VALUES(?,?,?)"
	o := orm.NewOrm()
	xx, err := o.Raw(sql, siteid, cname, cnums).Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xx)
}

//updateSiteDetail
func updateSiteDetail(siteid int, cname string, cnums int) {
	sql := "update sitedetail set cnums=? where siteid=? and cname=?"
	o := orm.NewOrm()
	_, err := o.Raw(sql, cnums, siteid, cname).Exec()
	if err != nil {
		fmt.Println(err)
	}
}

//从文件读取域名
func InitSitesFromFile() {
	items := ReadFileLines("domains.txt")
	fmt.Println(len(items))
	sql := "INSERT INTO site(domain,host,username,passwd,dbname,flag)VALUES(?,?,?,?,?,?)"
	o := orm.NewOrm()

	for _, item := range items {
		vals := strings.Split(item, ",")

		if len(vals) == 5 {
			o.Raw(sql, vals[0], vals[1], vals[2], vals[3], strings.TrimSpace(strings.TrimRight(strings.TrimRight(vals[4], "\n"), "\r")), 0).Exec()
		}

	}

}

//orm.RegisterDataBase("db1", "mysql", "root:root@/orm_db2?charset=utf8")
//orm.RegisterDataBase("db2", "sqlite3", "data.db")
//插入产品数据，默认所有网站的产品表都是已经清空的

func PostProduct(si SiteInfo, pi ProductInfo) bool {
	fmt.Println("插入的产品数据为：")
	fmt.Println(pi)
	//1.根据si.id获取数据库
	db := fmt.Sprintf("db%d", si.Id)
	o1 := orm.NewOrm()
	o1.Using(db)
	//2.if pi.parent not in si.clist {
	//   插入类别
	//	 }
	cids := strings.Split(pi.Parent, SEP_CHARS)
	cnames := strings.Split(strings.Replace(pi.ParentName, SEP_CHARS+SEP_CHARS, SEP_CHARS, -1), SEP_CHARS)

	cdbs := make(map[string]int)
	//把si.clist分割成数组
	for k, _ := range si.CList {
		ks := strings.Split(k, SEP_CHARS)
		for _, v := range ks {
			if strings.TrimSpace(v) != "" {
				cdbs[v] = 1
			}

		}

	}
	//去掉第一个,可以做父类为0
	for i := 1; i < len(cids); i++ {
		//如果不存在数据库中，则插入
		_, ok := cdbs[cids[i]]
		if !ok {
			fmt.Println(cids[i])
			fmt.Println(cnames[i])
			//类别编号直接使用现有的
			sql1 := "INSERT INTO categories(categories_id,parent_id,sort_order,date_added,last_modified,categories_status) values(?,?,?,now(),now(),?)"
			res, err := o1.Raw(sql1, cids[i], cids[i-1], 0, 1).Exec()
			if err != nil {
				fmt.Println(err)
				return false
			} else {
				id, _ := res.LastInsertId()
				fmt.Println(id)
			}
			sql2 := "INSERT INTO categories_description (categories_id,language_id,categories_name,categories_description) values(?,?,?,?)"
			res, err = o1.Raw(sql2, cids[i], 1, cnames[i], "").Exec()
			if err != nil {
				fmt.Println(err)
				//				return false
			} else {
				id, _ := res.LastInsertId()
				fmt.Println(id)
			}
		}
	}

	//3.插入产品，产品编号直接使用现有的,防止重复

	sqlProduct := "INSERT INTO products(products_id,products_type,products_quantity," +
		"products_model,products_image,products_price,products_virtual," +
		"products_date_added,products_last_modified,products_date_available,products_weight," +
		"products_status,get_url)values(?,1,1000,?,?,?,0,now(),now(),now(),0.25,1,?)"
	_, err := o1.Raw(sqlProduct, pi.Id, pi.Asin, pi.MainImage, pi.Price, pi.Url).Exec()
	if err != nil {
		fmt.Println(err)
		return false
	}

	sqlProduct = "INSERT INTO products_description(products_id,language_id,products_name,products_description," +
		"products_url,products_viewed,pd2)VALUES(?,1,?,?,'',20,'')"
	_, err = o1.Raw(sqlProduct, pi.Id, pi.Title, pi.Desc).Exec()
	if err != nil {
		fmt.Println(err)
	}

	sqlProduct = "INSERT INTO products_to_categories(products_id,categories_id)VALUES(?,?)"
	_, err = o1.Raw(sqlProduct, pi.Id, cids[len(cids)-1]).Exec()
	if err != nil {
		fmt.Println(err)
	}
	//4.if pi.attr 不为空 {
	//		插入属性
	//	}
	if pi.Attr != "" {
		//切割属性字符串
		attrs := strings.Split(pi.Attr, SEP_CHARS)
		for i, v := range attrs {
			if strings.TrimSpace(v) != "" {
				var sortorder int
				var products_options_values_id int64
				var err2 error
				if i == 0 {
					sortorder = -100
				} else {
					sortorder = i
				}
				sql := "INSERT INTO products_options_values(language_id,products_options_values_name,products_options_values_sort_order)VALUES(1,?,?)"
				res, err := o1.Raw(sql, v, sortorder).Exec()
				if err == nil {
					products_options_values_id, err2 = res.LastInsertId()
					if err2 != nil {
						fmt.Println(err2)
						return false
					}
				}

				sql = "INSERT INTO products_options_values_to_products_options" +
					"(products_options_values_to_products_options_id,products_options_id,products_options_values_id)" +
					"VALUES(?,1,?)"
				_, err = o1.Raw(sql, products_options_values_id, products_options_values_id).Exec()
				if err != nil {
					fmt.Println(err)
				}
				sql = "INSERT INTO products_attributes(products_id,	options_id,options_values_id,price_prefix,products_attributes_weight_prefix)" +
					"VALUES(?,1,?,'+','+')"
				_, err = o1.Raw(sql, pi.Id, products_options_values_id).Exec()
				if err != nil {
					fmt.Println(err)

				}
			}

		}

	}

	return true
}
func PostProductLocal(pi ProductInfo) error {
	//获取类别数组
	cids := strings.Split(pi.Parent, SEP_CHARS)
	cid := cids[len(cids)-1]

	o1 := orm.NewOrm()
	o1.Using("zencart")
	//3.插入产品，产品编号直接使用现有的,防止重复

	sqlProduct := "INSERT INTO products(products_id,products_type,products_quantity," +
		"products_model,products_image,products_price,products_virtual," +
		"products_date_added,products_last_modified,products_date_available,products_weight," +
		"products_status,get_url)values(?,1,1000,?,?,?,0,now(),now(),now(),0.25,1,?)"
	_, err := o1.Raw(sqlProduct, pi.Id, pi.Asin, pi.MainImage, pi.Price, pi.Url).Exec()
	if err != nil {
		errinfo := fmt.Sprintf("%s,%v", sqlProduct, err)
		return errors.New(errinfo)
	}

	sqlProduct = "INSERT INTO products_description(products_id,language_id,products_name,products_description," +
		"products_url,products_viewed,pd2)VALUES(?,1,?,?,'',20,?)"
	_, err = o1.Raw(sqlProduct, pi.Id, pi.Title, pi.Desc, pi.Imgs).Exec()
	if err != nil {
		errinfo := fmt.Sprintf("%s,%v", sqlProduct, err)
		return errors.New(errinfo)
	}

	sqlProduct = "INSERT INTO products_to_categories(products_id,categories_id)VALUES(?,?)"
	_, err = o1.Raw(sqlProduct, pi.Id, cid).Exec()
	if err != nil {
		errinfo := fmt.Sprintf("%s,%v", sqlProduct, err)
		return errors.New(errinfo)
	}
	//4.if pi.attr 不为空 {
	//		插入属性
	//	}
	if pi.Attr != "" {
		//切割属性字符串
		attrs := strings.Split(pi.Attr, SEP_CHARS)
		for i, v := range attrs {
			if strings.TrimSpace(v) != "" {
				var sortorder int
				var products_options_values_id int64
				var err2 error
				if i == 0 {
					sortorder = -100
				} else {
					sortorder = i
				}
				sql := "INSERT INTO products_options_values(language_id,products_options_values_name,products_options_values_sort_order)VALUES(1,?,?)"
				res, err := o1.Raw(sql, v, sortorder).Exec()
				if err == nil {
					products_options_values_id, err2 = res.LastInsertId()
					if err2 != nil {
						errinfo := fmt.Sprintf("%s,%v", sql, err)
						return errors.New(errinfo)
					}
				}

				sql = "INSERT INTO products_options_values_to_products_options" +
					"(products_options_values_to_products_options_id,products_options_id,products_options_values_id)" +
					"VALUES(?,1,?)"
				_, err = o1.Raw(sql, products_options_values_id, products_options_values_id).Exec()
				if err != nil {
					errinfo := fmt.Sprintf("%s,%v", sql, err)
					return errors.New(errinfo)
				}
				sql = "INSERT INTO products_attributes(products_id,	options_id,options_values_id,price_prefix,products_attributes_weight_prefix)" +
					"VALUES(?,1,?,'+','+')"
				_, err = o1.Raw(sql, pi.Id, products_options_values_id).Exec()
				if err != nil {
					errinfo := fmt.Sprintf("%s,%v", sql, err)
					return errors.New(errinfo)
				}
			}

		}

	}
	return nil
}
