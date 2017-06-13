package main

import (
	"fmt"
	//	"strings"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//测试单个插入
//func TestInitSitesFromFile(t *testing.T) {
//	InitSitesFromFile()
//}

//func TestAddSiteDetail(t *testing.T) {
//	addSiteDetail(1, "hello", 20)
//}
//func TestUpdateSiteDetail(t *testing.T) {
//	updateSiteDetail(1, "hello", 22)
//}
//func TestSetDomainCompleted(t *testing.T) {
//	SetDomainCompleted(1)
//}
//func TestInitSiteInfoFromDB(t *testing.T) {
//	xx := InitSiteInfoFromDB()
//	t.Log(xx)
//}

//func TestAdd(t *testing.T) {
//	var pi ProductInfo
//	pi.Parent = "hello"
//	xx := InitSiteInfoFromDB()
//	for _, v := range xx {
//		psi := &v
//		ret := psi.Add(pi)

//		t.Log(ret)
//	}

//	t.Log(xx)
//}
//多线程测试
//func TestMultiAdd(t *testing.T) {
//	var pi ProductInfo
//	pi.Parent = "hello"
//	xx := InitSiteInfoFromDB()
//	for _, v := range xx {
//		psi := &v
//		ret := psi.Add(pi)

//		t.Log(ret)
//	}

//	t.Log(xx)
//}
//func TestCinsert(t *testing.T) {
//	//把pi.parent分割成数组
//	parent := "0###35###429"
//	parentname := "###Beauté Prestige######Soins"
//	clist := make(map[string]int)
//	clist["0###39###3429"] = 1
//	clist["0###36###1429"] = 1
//	clist["0###37###2429"] = 1
//	clist["0###35###329"] = 1
//	cids := strings.Split(parent, SEP_CHARS)
//	cnames := strings.Split(strings.Replace(parentname, SEP_CHARS+SEP_CHARS, SEP_CHARS, -1), SEP_CHARS)

//	cdbs := make(map[string]int)
//	for k, _ := range clist {
//		ks := strings.Split(k, SEP_CHARS)
//		for _, v := range ks {
//			if strings.TrimSpace(v) != "" {
//				cdbs[v] = 1
//			}

//		}

//	}
//	fmt.Println(cids)
//	fmt.Println(len(cnames))
//	fmt.Println(clist)
//	fmt.Println(cdbs)
//	//去掉第一个
//	for i := 1; i < len(cids); i++ {
//		//如果不存在数据库中，则插入
//		_, ok := cdbs[cids[i]]
//		if !ok {
//			fmt.Println(cids[i])
//			fmt.Println(cnames[i])
//		}
//	}

//}
func TestInsert(t *testing.T) {
	//	piid := "3333"
	orm.RegisterDataBase("db1", "mysql", "root:yZpmTXkuSKdKbYGn@tcp(96.126.123.38:3306)/zadmin_menhirshop?charset=utf8")
	o1 := orm.NewOrm()
	o1.Using("db1")
	sql1 := "INSERT INTO categories(categories_id,parent_id,sort_order,date_added,last_modified,categories_status)" + "values(?,?,?,now(),now(),?);"
	res, err := o1.Raw(sql1, 1133333, 12, 0, 1).Exec()
	if err != nil {
		fmt.Println(err)
	} else {
		id, _ := res.LastInsertId()
		fmt.Println(id)
	}
	//	sql2 := "INSERT INTO categories_description (categories_id,language_id,categories_name,categories_description)" + "values(?,?,?,?)"
	//	res, err = o1.Raw(sql2, 333333, 1, "Cuisine d'Excellence", "").Exec()
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		id, _ := res.LastInsertId()
	//		fmt.Println(id)
	//	}

	//	sqlProduct := "INSERT INTO products(products_id,products_type,products_quantity," +
	//		"products_model,products_image,products_price,products_virtual," +
	//		"products_date_added,products_last_modified,products_date_available,products_weight," +
	//		"products_status,get_url)values(?,1,1000,?,?,?,0,now(),now(),now(),0.25,1,?)"
	//	res, err = o1.Raw(sqlProduct, 3333, "ddddddddddd", "https://images-na.ssl-images-amazon.com/images/I/814z9tkmi2L._SX425_.jpg", "113.33", "amazon.fr").Exec()
	//	fmt.Println(err)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		id, _ := res.LastInsertId()
	//		fmt.Println(id)
	//	}

	//	sqlProduct = "INSERT INTO products_description(products_id,language_id,products_name,products_description," +
	//		"products_url,products_viewed,pd2)VALUES(?,1,?,?,'',20,'')"
	//	res, err = o1.Raw(sqlProduct, piid, "ddddddddd", "ccccccccccc").Exec()
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		id, _ := res.LastInsertId()
	//		fmt.Println(id)
	//	}

	//	sqlProduct = "INSERT INTO products_to_categories(products_id,categories_id)VALUES(?,?)"
	//	res, err = o1.Raw(sqlProduct, piid, 111).Exec()
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		id, _ := res.LastInsertId()
	//		fmt.Println(id)
	//	}
	//	piAttr := "###Sélectionner###36 EU###37 EU###38 EU###39 EU###40 EU"

	//	if piAttr != "" {
	//		//切割属性字符串
	//		attrs := strings.Split(piAttr, SEP_CHARS)
	//		for i, v := range attrs {
	//			if strings.TrimSpace(v) != "" {
	//				var sortorder int
	//				var products_options_values_id int64
	//				var err2 error
	//				if i == 0 {
	//					sortorder = -100
	//				} else {
	//					sortorder = i
	//				}
	//				sql := "INSERT INTO products_options_values(language_id,products_options_values_name,products_options_values_sort_order)VALUES(1,?,?)"
	//				res, err := o1.Raw(sql, v, sortorder).Exec()
	//				if err == nil {
	//					products_options_values_id, err2 = res.LastInsertId()
	//					if err2 != nil {
	//						fmt.Println(err2)
	//					}
	//				}

	//				sql = "INSERT INTO products_options_values_to_products_options" +
	//					"(products_options_values_to_products_options_id,products_options_id,products_options_values_id)" +
	//					"VALUES(?,1,?)"
	//				res, err = o1.Raw(sql, products_options_values_id, products_options_values_id).Exec()
	//				if err != nil {
	//					fmt.Println(err)
	//				} else {
	//					id, _ := res.LastInsertId()
	//					fmt.Println(id)
	//				}
	//				sql = "INSERT INTO products_attributes(products_id,	options_id,options_values_id,price_prefix,products_attributes_weight_prefix)" +
	//					"VALUES(?,1,?,'+','+')"
	//				res, err = o1.Raw(sql, piid, products_options_values_id).Exec()
	//				if err != nil {
	//					fmt.Println(err)
	//				} else {
	//					id, _ := res.LastInsertId()
	//					fmt.Println(id)
	//				}
	//			}
	//		}

	//	}
}

//func TestPostProduct(t *testing.T) {

//	orm.RegisterDataBase("db1", "mysql", "root:yZpmTXkuSKdKbYGn@tcp(96.126.123.38:3306)/zadmin_menhirshop?charset=utf8")

//	var si SiteInfo
//	var pi ProductInfo
//	ret := PostProduct(si, pi)
//}
