package main

import (
	"testing"
)

//func TestUrlFilter(t *testing.T) {
//	var params DBContext
//	params.dbName = MYSQL_DB
//	params.tableName = CATEGORY_TABLE

//	sbf := newUrlFilter(params)
//	url := "/informatique-ordinateurs-imprimantes-r%C3%A9seaux-composants/b/ref=sd_allcat_allpc?ie=UTF8&node=34085803"
//	fmt.Println(sbf.TestAndAdd(url))
//	url = "/informatique-ordinateurs-imprimantes-r%C3%A9seaux-composants/b/ref=sd_allcat_allpc?ie=UTF8&node=34085803111"
//	fmt.Println(sbf.TestAndAdd(url))

//}
func TestNewUrlFilter(t *testing.T) {
	var params DBContext
	params.tableName = CATEGORY_TABLE

	db := &DB{context: params}
	newUrlFilter(db)

}
