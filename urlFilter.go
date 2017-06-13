package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/tylertreat/BoomFilters"
)

type UrlFilter struct {
	mux    sync.Mutex
	filter *boom.ScalableBloomFilter
}

func newUrlFilter(db *DB) *UrlFilter {

	u := &UrlFilter{}
	//	u.filter = boom.NewDefaultScalableBloomFilter(0.01)
	u.filter = boom.NewScalableBloomFilter(1000000, 0.01, 0.8)
	MM := 1000000
	x1 := db.count(0)
	x2 := db.count(1)
	x := x1 + x2
	t := x / MM
	for i := 0; i <= t; i++ {
		start := time.Now()
		urls, _ := db.GetAll(i*MM, MM)
		for _, v := range urls {
			u.filter.TestAndAdd([]byte(v.Url))
		}
		fmt.Println(len(urls))
		fmt.Println(time.Now().Sub(start).Seconds())
	}

	return u
}
func (self *UrlFilter) TestAndAdd(url string) bool {
	self.mux.Lock()
	ret := self.filter.TestAndAdd([]byte(url))
	self.mux.Unlock()
	return ret
}
