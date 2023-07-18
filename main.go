package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

//https://www.xbiquge.la/48/48900/
//爬取小说

//获取小说名
func getStoryName(b []byte) string {
	//找到<h1>里的名称
	name := regexp.MustCompile(`<h1>.*</h1>`)
	str := name.FindAll(b, -1)
	s2 := string(str[0])
	// 找到>第一次出现的位置<最后出现的位置截取
	n1 := strings.Index(s2, ">")
	n2 := strings.LastIndex(s2, "<")
	s3 := s2[n1+1 : n2]
	// fmt.Println(s3)
	return s3
}

//获取小说全部章节和地址
func getPage(b []byte) []map[string]string {
	page := regexp.MustCompile(`<dd>.*</dd>`)
	str := page.FindAll(b, -1)
	//定义切片map保存地址和章节
	//切片长度
	n := len(str)
	var storyPage = make([]map[string]string, n)

	for i := 0; i < n; i++ {
		//截取地址
		// 获取''前后出现的位置
		s := string(str[i])
		i1 := strings.Index(s, "'")
		i2 := strings.LastIndex(s, "'")
		si := s[i1+1 : i2]
		//截取名称
		// l1 := strings.Index(si, "第")
		l1 := s[i2+3:]
		l2 := strings.Index(l1, "<")
		sl := l1[:l2]
		// fmt.Println("si ======", si)
		//保存到map
		storyPage[i] = make(map[string]string)
		storyPage[i][si] = sl

	}

	return storyPage

}

//将map地址和小说写进文本
func inputStore(sPage []map[string]string, newUrl string, fname string) {
	//将map小说名取出和小说章节取出写入文件

	for _, s := range sPage {
		for k, v := range s {
			//新章节地址
			newurl := newUrl + k
			// println("newurl==", newurl)
			//新章节地址数据
			date := getData(newurl)
			//写入文件
			prinListfile(date, v, fname)

		}

	}
}

//写入文件
func prinListfile(contes []byte, pageName string, fname string) {
	//获取文本
	re := regexp.MustCompile(`&nbsp;&nbsp;.*<br /><br />`)
	str := re.FindAll(contes, -1)
	reg := regexp.MustCompile(`[[:^ascii:]]+`)

	str2 := reg.FindAll(str[0], -1)
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("open file error=%v\n", err)
		return
	}
	ls, err := io.WriteString(file, pageName+"\n")
	if err != nil {
		panic(err)
	}
	fmt.Printf("写入 %d 个字节n", ls)
	for i := 0; i < len(str2); i++ {
		n, err := io.WriteString(file, "  "+string(str2[i])+"\n")
		if err != nil {
			panic(err)
		}
		fmt.Printf("写入 %d 个字节n", n)
	}
	defer file.Close()
}

//地址内容获取
func getData(url string) []byte {
	//连接小说地址
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//关闭连接
	defer resp.Body.Close()
	//无法获取数据关闭
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error!!", resp.StatusCode)

	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return all
}

func main() {

	beginTime := time.Now()
	//1.获取小说名称并创建对应文件
	// url := "https://www.xbiquge.la/48/48900/"
	// url := "https://www.xbiquge.la/84/84419/"
	// url := "https://www.xbiquge.la/84/84522/"

	//手动输入
	var url string
	fmt.Scan(&url)
	fmt.Println(url)
	//连接小说地址
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//关闭连接
	defer resp.Body.Close()
	//无法获取数据关闭
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error!!", resp.StatusCode)
		return
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//获取小说
	sroryName := getStoryName(all)

	//创建小说名称的文件
	sname := sroryName + ".txt"

	//返回地址+章节
	storyPage := getPage(all)

	//将url的网址获取
	urli := strings.LastIndex(url, ".")
	//新地址
	newurl := url[:urli+3]
	//传入章节目录和地址,替换地址,雄安说名称
	inputStore(storyPage, newurl, sname)
	fmt.Println("******爬取结束********")
	endTiem := time.Since(beginTime)
	fmt.Println("爬取耗时:", endTiem/1000)

}
