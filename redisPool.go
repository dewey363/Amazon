package main

import (
	"encoding/json"
	//	"fmt"

	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	PROXY_KEY             string = "proxy_queue_key"
	PROXY_FAILED_KEY      string = "proxy_failed_queue_key" //被blocked的IP池，第二天再启用
	CATEGORY_URL_KEY      string = "category_url_queue_key"
	PRODUCT_LISTS_URL_KEY string = "product_lists_url_queue_key"
	PRODUCT_LIST_URL_KEY  string = "product_list_url_queue_key"
	PRODUCT_INFO_URL_KEY  string = "product_info_url_queue_key"
)

//前缀小写，限制外部文件访问
type redisPool struct {
	pool *redis.Pool
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   500,
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				dbLog.Println("Create one connection, a error occurs.", server)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					dbLog.Debugln(err)
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}
func OpenRedisPool() *redisPool {
	pool := newPool(REDIS_POOL_ADDR, REDIS_POOL_PASSWD)
	rPool := &redisPool{
		pool: pool,
	}
	return rPool
}
func CloseRedisPool(rpool *redisPool) error {
	return rpool.pool.Close()
}

//CategroyUrlInfo进队列
func (rPool *redisPool) PushURL(data UrlInfo, queueName string) {
	conn := rPool.pool.Get()
	defer conn.Close()
	value, _ := json.Marshal(data)
	_, err := conn.Do("lpush", queueName, string(value))
	if err != nil {
		//		fmt.Printf("lpush %v", err)
		return
	}
}

//CategroyUrlInfo出队列
func (rPool *redisPool) PopURL(queueName string, timeout int) *UrlInfo {
	var data UrlInfo
	conn := rPool.pool.Get()
	defer conn.Close()
	value, err := redis.Strings(conn.Do("brpop", queueName, timeout))
	if err != nil {
		//		fmt.Println(err)
		return nil
	}
	if len(value) == 2 {
		err = json.Unmarshal([]byte(value[1]), &data)
		if err != nil {
			//			fmt.Println(err)
			return nil
		}
		return &data
	}
	return nil

}

//队列长度
func (rPool *redisPool) LenQueue(queueName string) int {

	conn := rPool.pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("llen", queueName))
	if err != nil {
		//		fmt.Println(err)
		return 0
	}

	return value

}

//清空队列
func (rPool *redisPool) EmptyQueue(queueName string) int {

	conn := rPool.pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("del", queueName))
	if err != nil {
		//		fmt.Println(err)
		return 0
	}

	return value

}

//Proxy进队列
func (rPool *redisPool) PushProxy(data IPInfo) {
	conn := rPool.pool.Get()
	defer conn.Close()
	value, _ := json.Marshal(data)
	_, err := conn.Do("lpush", PROXY_KEY, string(value))
	if err != nil {
		//		fmt.Printf("lpush %v", err)
		return
	}
}

//IPInfo出队列 当timeout=0是阻塞
func (rPool *redisPool) PopProxy(timeout int) *IPInfo {
	var data IPInfo
	conn := rPool.pool.Get()
	defer conn.Close()
	value, err := redis.Strings(conn.Do("brpop", PROXY_KEY, timeout))
	if err != nil {
		//		fmt.Println(err)
		return nil
	}
	if len(value) == 2 {
		err = json.Unmarshal([]byte(value[1]), &data)
		if err != nil {
			//			fmt.Println(err)
			return nil
		}
		if data.Flag == 1 {
			WaitSecond(time.Now(), FAILED_PROXY_SECONDS)
		}
		//如果有失败时间，则等待一段时间
		return &data
	}
	return nil

}
