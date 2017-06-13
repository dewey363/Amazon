package main

import (
	//	"fmt"
	//	"net/http"
	"testing"
)

//func TestGetCategroyUrlInfo(t *testing.T) {
//	cis, _ := GetCategroyUrlInfo(NORMAL, 100)
//	t.Logf("cis rows:%d\n", len(cis))
//}

//func TestUpdateBatchCategoryInfoToDB(t *testing.T) {
//	cis, _ := GetCategroyUrlInfo(1, 100)
//	t.Logf("cis rows:%d\n", len(cis))
//	UpdateBatchCategoryInfoToDB(NORMAL, cis)
//	cis, _ = GetCategroyUrlInfo(NORMAL, 100)
//	t.Logf("cis rows:%d\n", len(cis))
//}

//func TestUpdateCategoryInfoToDB(t *testing.T) {
//	cis, _ := GetCategroyUrlInfo(0, 100)
//	t.Logf("cis rows:%d\n", len(cis))
//	UpdateBatchCategoryInfoToDB(1, cis)
//	ids := []int{12, 29, 33, 36, 37,
//		39, 41, 43, 45, 46, 47, 48, 50, 51, 52, 54, 55, 56,
//		57, 58, 59, 60, 61, 64, 65, 66, 71, 76, 82,
//		86, 88, 91, 92, 93}
//	for _, v := range ids {
//		fmt.Println(v)
//		var ci CategroyUrlInfo
//		ci.Id = v
//		ci.Flag = NORMAL
//		ci.Nums = 0
//		UpdateCategoryInfoToDB(ci)
//	}

//}
//func TestCheckCategroyUrlExist(t *testing.T) {
//	fmt.Println(CheckCategroyUrlExist("/Arts-Crafts-Sewing/b/ref=topnav_storetab_ac/156-8609269-7173439?ie=UTF8&node=2617941011"))
//	fmt.Println(CheckCategroyUrlExist("/Arts-Crafts-Sewing/b/ref=topnav_storetab_ac/156-8609269-7173439?ie=UTF8&node=26179410111"))
//}
func TestCategroyUrlDB(t *testing.T) {
	var context DBContext
	context.tableName = CATEGORY_TABLE
	cuDB := NewDB(context)
	cuDB.initDBfromFile("fr.txt")
}
