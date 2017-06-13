package main

//读取HTML配置
var (
	MAX_LEVEL                int = setting.DefaultInt("max_level", 4)       //最大5层
	USE_PROXY                int = 1                                        //setting.DefaultInt("use_proxy", 1)       //        //是否使用代理,0不使用，1使用
	HTTP_TIME_OUT            int = setting.DefaultInt("http_time_out", 180) //每个连接超时秒数
	MAX_GET_TIME             int = setting.DefaultInt("max_get_time", 4*60) // 每个IP每个页面至少4分钟，如果不够则等待
	PREFETCH_START           int = setting.DefaultInt("prefetch_start", 0)  //预读记录的开始位置
	RUN_FLAG                 int = setting.DefaultInt("run_flag", 1)        //1.抓取类别，2.抓取产品链接3.抓取产品信息
	MAX_PROCS                int = setting.DefaultInt("max_procs", 1)       //最大线程数,线程数是IP的80%，能够保证IP的充足。如果没有代理，则设置为1
	MAX_PRODUCTS_IN_CATEGORY int = 48 * 400                                 //一个类别中最多产品数
	MAX_PRODUCTS_IN_SITE     int = 20000                                    //每个站点产品数
	MAX_CATEGORYS_IN_SITE    int = 100                                      //每个站点类别数

	LANG string = setting.DefaultString("lang", "en")
	//		AMAZON_LINK              string = setting.DefaultString("amazon_link", "https://www.amazon.com")
	AMAZON_LINK          string = "https://www.amazon.com"
	FAILED_PROXY_SECONDS int    = setting.DefaultInt("failed_proxy_seconds", 30*60) // 30 * 60 //IP被封了需要等待多久，以秒为单位
)

//数据库参数配置
var (
	//	PRODUCT_LISTS_TABLE string = "root:Aazsedcftgbhujmko123456@/amazon" //产品列表页面，准备生成400页
	PRODUCT_LISTS_TABLE string = "product_lists"                                                       //产品列表页面，准备生成400页
	PRODUCT_LIST_TABLE  string = "product_list"                                                        //产品列表页面，准备获取产品url，每个页面大概48个产品
	PRODUCT_INFO_TABLE  string = "product_info"                                                        //产品信息数据库
	PRODUCT_URL_TALBE   string = "product_url"                                                         //产品url信息数据库
	CATEGORY_TABLE      string = "category"                                                            //产品类别页面链接表，用于获取子类别或者产品列表页面
	MYSQL_DB            string = setting.DefaultString("mysql_db", "root:root@/amazon?charset=utf8")   //产品列表页面，准备生成400页
	ZENCART_DB          string = setting.DefaultString("zencart_db", "root:root@/amazon?charset=utf8") //产品列表页面，准备生成400页

)

//日志服务器参数
var (
	//	LOG_SERVER string = setting.DefaultString("log_server", "45.33.34.189:12202") //日志服务器IP和端口
	//	LOCAL      string = setting.DefaultString("local", "测试主机")                    //本机IP或描述

	LOG_TIME int = setting.DefaultInt("log_time", 30) //发送频率  30分钟
)

//缓存参数配置
var (
	//	REDIS_POOL_ADDR   string = "127.0.0.1:6379"        //redis服务器地址
	//	REDIS_POOL_PASSWD string = "asdfghjkl123456!@#$%^" //redis服务器密码
	REDIS_POOL_ADDR   string = setting.DefaultString("redis_pool_addr", "127.0.0.1:6379") //redis服务器地址
	REDIS_POOL_PASSWD string = setting.DefaultString("redis_pool_passwd", "")             //redis服务器密码

	Proc_Fetch_Size int = 19                          //每个线程每次从数据库预先读取个数，
	Min_Queue_Size  int = 2*MAX_PROCS + 1             //线程池中最少个数
	PreFetchTime    int = MAX_GET_TIME                //每次从数据库读取间隔时间,单位秒数
	PreFetchSize    int = Proc_Fetch_Size * MAX_PROCS //每次从数据库预先读取总个数，

)

/*
类型信息flag标记值定义
*/
const (
	NORMAL    int    = 0
	COMPLETED int    = 1
	SEP_CHARS string = "###"
)
