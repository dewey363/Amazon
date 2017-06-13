package main

import (
	"bufio"
	//	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

/*
等待函数，从start开始到现在，如果没有到 longSeconds ,则等待
*/
func WaitSecond(start time.Time, longSeconds int) {
	//	fmt.Printf("等待中...%d", longSeconds)
	diffTime := time.Now().Sub(start)
	//	fmt.Println(diffTime)
	d := time.Duration(longSeconds)*time.Second - diffTime
	//	fmt.Println(d)
	if d > 0 {
		time.Sleep(d)
	}
	//	fmt.Println("等待完成")
}

//等待随时数
func RandSecond(second int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ind := r.Intn(120) - 60 + MAX_GET_TIME //180-300之间
	return ind
}

/*
错误处理
*/
func checkErr(err error) {
	if err != nil {
		dbLog.Errorln(err)
		panic(err)
	}
}

// fileName:文件名字(带全路径),文件必须是已经存在
// content: 写入的内容
func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		//		fmt.Println(fileName + " file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}

	return err
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//保存html到文件中
func WriteHTMLToFile(cid int, body string) error {
	//把返回的html写入到文件中
	saveDir := "html/"
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+saveDir, os.ModePerm) //生成多级目录
	if err != nil {
		//		fmt.Println(err)
		return err
		//		return "", errors.New("输出到文件错误")
	}
	htmlFile := saveDir + strconv.Itoa(cid) + ".html"
	return appendToFile(htmlFile, body)
}

//写cookie到文件中
func WriteCookieToFile(cid int, cookiestr string) {
	//把返回的html写入到文件中
	saveDir := "cookie/"
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+saveDir, os.ModePerm) //生成多级目录
	if err == nil {
		htmlFile := saveDir + strconv.Itoa(cid) + ".txt"
		appendToFile(htmlFile, cookiestr)
	}

}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) ([]string, error) {
	var files []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	//	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			//			files = append(files, dirPth+PthSep+fi.Name())
			files = append(files, path.Join(dirPth, fi.Name()))
		}
	}
	return files, nil
}

//获取路径中的文件名
func GetFileNameOnly(dir string) string {
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(dir) //获取文件名带后缀
	//	fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀
	//	fmt.Println("fileSuffix =", fileSuffix)

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名
	return filenameOnly
}

/*
读取文件内容
*/
func ReadAllText(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	// fmt.Println(string(fd))
	return string(fd)
}

/*
按行读取文件

*/
func ReadFileLines(filename string) []string {
	var ret []string
	f, err := os.Open(filename) //打开文件
	defer f.Close()             //打开文件出错处理

	if nil == err {
		buff := bufio.NewReader(f) //读入缓存
		for {
			line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				break
			}
			if strings.TrimSpace(line) != "" {

				ret = append(ret, line)

			}

		}
	}

	return ret
}

//生成count个[start,end)结束的不重复的随机数
func generateRandomNumber(start int, end int, count int) []string {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]string, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		snum := strconv.Itoa(num)
		//查重
		exist := false
		for _, v := range nums {
			if v == snum {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, snum)
		}
	}

	return nums
}
